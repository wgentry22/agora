package heartbeat

import (
  "github.com/gin-gonic/gin"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/wgentry22/agora/types/config"
)

var (
  registeredPacers = make(map[string]Pacer)
)

func RegisterPacers(pacers ...Pacer) {
  m.Lock()
  defer m.Unlock()

  for _, pacer := range pacers {
    if _, ok := registeredPacers[pacer.Component()]; !ok {
      registeredPacers[pacer.Component()] = pacer
    }
  }
}

func ClearPacers() {
  m.Lock()
  defer m.Unlock()

  registeredPacers = make(map[string]Pacer)
}

type Pacer interface {
  Component() string
  RegisterWith(registry *prometheus.Registry)
}

func MetricsHandler(conf config.Heartbeat) func(*gin.Context) {
  registry := prometheus.NewRegistry()

  registry.MustRegister(prometheus.NewGoCollector())

  for _, pacer := range registeredPacers {
    pacer.RegisterWith(registry)
  }

  handler := promhttp.HandlerFor(
    registry,
    promhttp.HandlerOpts{
      Registry:            registry,
      Timeout:             conf.Timeout.Read,
      EnableOpenMetrics:   true,
    },
  )

  return func(c *gin.Context) {
    handler.ServeHTTP(c.Writer, c.Request)
  }
}
