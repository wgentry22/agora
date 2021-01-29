package orm

import (
  "database/sql"
  "reflect"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/wgentry22/agora/modules/heartbeat"
)

var (
  MetricsNames = []string{
    "max_open_connections",
    "open_connections",
    "connections_in_use",
    "idle_connections",
    "wait_count",
    "wait_duration",
    "max_idle_closed",
    "max_idle_time_closed",
    "lifetime_closed",
  }

  MetricsHelp = []string{
    "Maximum number of open connections to the database.",
    "The number of established connections both in use and idle.",
    "The number of connections currently in use.",
    "The number of idle connections.",
    "The total number of connections waited for.",
    "The total time blocked waiting for a new connection.",
    "The total number of connections closed due to SetMaxIdleConns.",
    "The total number of connections closed due to SetConnMaxIdleTime.",
    "The total number of connections closed due to SetConnMaxLifetime.",
  }
)

func ormMetricDescriptions() []*prometheus.Desc {
  descriptions := make([]*prometheus.Desc, len(MetricsNames))

  for i, name := range MetricsNames {
    descriptions[i] = prometheus.NewDesc(name, MetricsHelp[i], nil, nil)
  }

  return descriptions
}

func ormMetrics(stats sql.DBStats) []prometheus.Metric {
  metrics := make([]prometheus.Metric, len(MetricsNames))

  descriptions := ormMetricDescriptions()

  statsRef := reflect.ValueOf(stats)

  for i := 0; i < statsRef.NumField(); i++ {
    var metric float64

    if val, ok := statsRef.Field(i).Interface().(int); ok {
      metric = float64(val)
    } else if val, ok := statsRef.Field(i).Interface().(int64); ok {
      metric = float64(val)
    }

    metrics[i] = prometheus.MustNewConstMetric(descriptions[i], prometheus.GaugeValue, metric)
  }

  return metrics
}

func RegisterPacer() {
  m.Lock()
  defer m.Unlock()

  heartbeat.RegisterPacers(ormInstance)
}

func (o *orm) RegisterWith(registry *prometheus.Registry) {
  registry.MustRegister(&Collector{rawDB})
}

type Collector struct {
  rawDB *sql.DB
}

func (c *Collector) Describe(desc chan<- *prometheus.Desc) {
  for _, d := range ormMetricDescriptions() {
    desc <- d
  }
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
  for _, metric := range ormMetrics(c.rawDB.Stats()) {
    metrics <- metric
  }
}
