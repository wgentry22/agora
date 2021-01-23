package watcher_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/watcher"
	"github.com/wgentry22/agora/types/config"
	"os"
)

var _ = Describe("Watcher", func() {
	var (
		fileName = "temp.toml"
		filePath = fmt.Sprintf("%s/%s", os.TempDir(), fileName)
		cw       watcher.ConfigWatcher
	)

	JustBeforeEach(func() {
		if err := writeDataToFile(filePath, []byte{}); err != nil {
			panic(err)
		}
	})

	Describe("initializing from factory method", func() {
		var (
			errc chan error
		)

		JustBeforeEach(func() {
			errc = make(chan error, 1)
		})

		Context("when passed a directory", func() {
			It("should panic", func() {
				defer func() {
					if r := recover(); r != nil {
						err, isErr := r.(error)
						if !isErr {
							Fail("expected to panic when config watcher is passed a directory")
						}

						Expect(err).To(Equal(watcher.ErrWatcherGotDirectory))
					}
				}()

				cw = watcher.NewConfigWatcher(os.TempDir())
			})
		})

		Context("when passed a path that does not exist", func() {
			It("should panic", func() {
				defer func() {
					if r := recover(); r != nil {
						err, isErr := r.(error)
						if !isErr {
							Fail("expected to panic when config watcher is passed a directory")
						}

						Expect(err).To(Equal(watcher.ErrPathDoesNotExist("/does/not/exist")))
					}
				}()

				cw = watcher.NewConfigWatcher("/does/not/exist")
			})
		})

		Context("when passed a file", func() {
			It("should NOT panic", func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovering", r)
						Fail("expected NOT to panic when config watcher is passed a file")
					}
				}()

				cw = watcher.NewConfigWatcher(filePath)

				// Begin watching - config is empty at this point.
				errors := cw.Watch(errc)

				// Simulate change in watched file.
				err := writeDataToFile(filePath, preChangeData)
				Expect(err).To(BeNil())

				// Listen for changes on ConfigWatcher.
				appConfig, ok := <-cw.Changes()
				Expect(ok).To(BeTrue())

				// Assert config changed.
				Expect(appConfig.Info()).To(Equal(config.Info{
					Name: "agora",
					Version: config.SemanticVersion{
						Major: 1,
						Minor: 2,
						Patch: 3,
					},
					Env: config.QualityAssurance,
				}))

				// Assert errors channel sent nil on successful file change.
				Eventually(errors).Should(Receive(nil))

				// Simulate another file change.
				err = writeDataToFile(filePath, postChangeData)
				Expect(err).To(BeNil())

				// Listen for changes on ConfigWatcher.
				appConfig, ok = <-cw.Changes()
				Expect(ok).To(BeTrue())

				// Assert config changed
				Expect(appConfig.Info()).To(Equal(config.Info{
					Name: "agora",
					Version: config.SemanticVersion{
						Major: 1,
						Minor: 2,
						Patch: 4,
					},
					Env: config.Production,
				}))

				// Assert errors channel sent nil on successful file change.
				Eventually(errors).Should(Receive(nil))
			})
		})
	})
})

func writeDataToFile(filePath string, data []byte) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
