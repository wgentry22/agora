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
  "github.com/wgentry22/agora/modules/broker"
  "github.com/wgentry22/agora/modules/heartbeat"
  "github.com/wgentry22/agora/modules/logg"
  "github.com/wgentry22/agora/modules/orm"
  "github.com/wgentry22/agora/types/config"
)

var (
  logger    logg.Logger
  consumer  broker.Consumer
  publisher broker.Publisher
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

func (a *Application) Publisher() broker.Publisher {
  return publisher
}

func (a *Application) Consumer() broker.Consumer {
  return consumer
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

  if a.conf.Broker().Role == config.RoleProducer {
    publisher = broker.NewPublisher(a.conf.Broker())
  } else if a.conf.Broker().Role == config.RoleConsumer {
    consumer = broker.NewConsumer(a.conf.Broker())
  }
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

  if consumer != nil {
    go func(c broker.Consumer) {
      if err := <-c.Errors(); err != nil {
        a.errors <- err
      }
    }(consumer)
  }

  if publisher != nil {
    go func(p broker.Publisher) {
      if err := <-p.Errors(); err != nil {
        a.errors <- err
      }
    }(publisher)
  }

  <-a.quit
  logger.Info("Shutting down server...")

  withTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  if err := server.Shutdown(withTimeout); err != nil {
    logger.WithError(err).Panic("Failed to shutdown server")
  }

  logger.Info("Server exiting.")
}
