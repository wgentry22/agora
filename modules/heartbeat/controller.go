package heartbeat

import (
	"github.com/wgentry22/agora/modules/api"
	"github.com/wgentry22/agora/types/config"
)

func NewHeartbeatController(conf config.Heartbeat) api.Controller {
	controller := api.NewController(conf.PathPrefix)

	controller.Register(api.NewGETRoute("/health", HealthHandler(conf)))
	controller.Register(api.NewGETRoute("/metrics", MetricsHandler(conf)))

	return controller
}
