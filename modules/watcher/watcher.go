package watcher

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/errwrap"
	"github.com/pelletier/go-toml"
	"github.com/wgentry22/agora/types/config"
	"io/ioutil"
	"os"
)

var (
	ErrFailedToGetWatcher = errors.New("failed to create file watcher")
	ErrFailedToWatchFile  = func(fileName string) error {
		return fmt.Errorf("failed to watch file `%s`", fileName)
	}
	ErrWatcherGotDirectory = errors.New("got a directory, wanted a file")
	ErrPathDoesNotExist    = func(filePath string) error {
		return fmt.Errorf("`%s` does not exist on this filesystem", filePath)
	}
)

type ConfigWatcher interface {
	Watch(errc chan error) <-chan error
	Changes() <-chan config.Application
}

type configWatcher struct {
	filePath string
	changes  chan config.Application
}

func (c *configWatcher) Watch(errc chan error) <-chan error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errc <- errwrap.Wrap(ErrFailedToGetWatcher, err)

		return errc
	}

	if err := watcher.Add(c.filePath); err != nil {
		errc <- errwrap.Wrap(ErrFailedToWatchFile(c.filePath), err)

		return errc
	}

	go func(errorChannel chan error) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if ok && event.Op&fsnotify.Write == fsnotify.Write {
					var appConfig config.Application

					if data, err := ioutil.ReadFile(event.Name); err != nil {
						errorChannel <- err
					} else if err = toml.Unmarshal(data, &appConfig); err != nil {
						errorChannel <- err
					} else {
						c.changes <- appConfig
						errorChannel <- nil
					}
				}
			case err, ok := <-watcher.Errors:
				if ok && err != nil {
					errorChannel <- err
				}
			}
		}
	}(errc)

	return errc
}

func (c *configWatcher) Changes() <-chan config.Application {
	return c.changes
}

func NewConfigWatcher(filePath string) ConfigWatcher {
	info, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		panic(ErrPathDoesNotExist(filePath))
	}

	if info.IsDir() {
		panic(ErrWatcherGotDirectory)
	}

	return &configWatcher{filePath, make(chan config.Application)}
}
