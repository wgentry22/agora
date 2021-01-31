package agora

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/wgentry22/agora/modules/api"
	"github.com/wgentry22/agora/modules/heartbeat"
	"github.com/wgentry22/agora/modules/logg"
	"github.com/wgentry22/agora/modules/orm"
	"github.com/wgentry22/agora/types/config"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	logger logg.Logger
)

type Application struct {
	errors chan error
	conf   config.Application
	router api.Router
	quit   chan os.Signal
}

func (a *Application) RegisterController(controllers ...api.Controller) {
	for _, controller := range controllers {
		a.router.Register(controller)
	}
}

func (a *Application) Errors() <-chan error {
	return a.errors
}

func (a *Application) Logger() logg.Logger {
	return logg.NewLogrusLogger(a.conf.Logging())
}

func (a *Application) Router() http.Handler {
	return a.router.Handler()
}

func (a *Application) Consumer() *kafka.Consumer {
	return config.NewSubscriber(a.conf.Broker())
}

func (a *Application) Producer() *kafka.Producer {
	return config.NewPublisher(a.conf.Broker())
}

func (a *Application) Setup() {
	go func() {
		if err, ok := <-a.errors; ok {
			logger.WithError(err).Panic("Unable to continue")
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err, isErr := r.(error)
			if isErr {
				errwrap.Walk(err, func(e error) {
					logger.WithError(e).Error("Found wrapped error")
				})
				a.errors <- err
			}
		}
	}()

	logger = logg.NewLogrusLogger(a.conf.Logging())

	orm.UseConfig(a.conf.DB())
	orm.UseLoggingConfig(a.conf.Logging())
	orm.RegisterPulser()
	orm.RegisterPacer()

	a.router.Register(heartbeat.NewHeartbeatController(a.conf.Heartbeat()))
}

func (a *Application) Start() {
	server := a.router.Server()

	signal.Notify(a.quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err, ok := <-a.errors; ok {
			logger.WithError(err).Error("Error while using application")
			a.quit <- syscall.SIGTERM
		}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errors <- err
		}
	}()

	<-a.quit
	logger.Info("Shutting down server...")

	withTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(withTimeout); err != nil {
		logger.WithError(err).Panic("Failed to shutdown server")
	}

	logger.Info("Server exiting.")
}
