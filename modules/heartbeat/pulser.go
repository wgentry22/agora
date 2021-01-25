package heartbeat

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wgentry22/agora/types/config"
	"net/http"
	"sync"
)

type HealthCheckStatus int8

const (
	StatusOK HealthCheckStatus = iota
	StatusWarn
	StatusCritical
	StatusUnknown
)

var (
	statusDisplay = []string{"ok", "warn", "critical", "unknown"}
	statusLookup  = map[string]HealthCheckStatus{
		"ok":       StatusOK,
		"warn":     StatusWarn,
		"critical": StatusCritical,
		"unknown":  StatusUnknown,
	}
	m                 sync.Mutex
	registeredPulsers = make(map[string]Pulser)
)

func (h HealthCheckStatus) String() string {
	return statusDisplay[h]
}

func (h HealthCheckStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *HealthCheckStatus) UnmarshalJSON(data []byte) error {
	var status string

	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	found, ok := statusLookup[status]
	if !ok {
		*h = StatusUnknown
	} else {
		*h = found
	}

	return nil
}

type HealthCheckResponse struct {
	App          string            `json:"app"`
	Version      string            `json:"version"`
	Env          string            `json:"env"`
	Status       HealthCheckStatus `json:"status"`
	Dependencies []Pulse           `json:"dependencies"`
}

func (h *HealthCheckResponse) HTTPStatus() int {
	status := http.StatusOK

	histo := make(map[HealthCheckStatus]int)

	for _, dep := range h.Dependencies {
		determineStatus(dep, histo)
	}

	if crit, ok := histo[StatusCritical]; ok && crit > 0 {
		h.Status = StatusCritical

		status = http.StatusServiceUnavailable
	} else if warn, ok := histo[StatusWarn]; ok && warn > 0 {
		h.Status = StatusWarn

		status = http.StatusFailedDependency
	} else {
		h.Status = StatusOK
	}

	return status
}

func determineStatus(pulse Pulse, histo map[HealthCheckStatus]int) {
	for _, dep := range pulse.Dependencies {
		determineStatus(dep, histo)
	}

	histo[pulse.Status]++
}

func NewHealthCheckResponse(info config.Info, deps []Pulse) HealthCheckResponse {
	return HealthCheckResponse{
		App:          info.Name,
		Version:      info.Version.String(),
		Env:          info.Env.String(),
		Status:       StatusUnknown,
		Dependencies: deps,
	}
}

type Pulser interface {
	Component() string
	Pulse(ctx context.Context, pulsec chan<- Pulse)
}

type Pulse struct {
	Component    string            `json:"component"`
	Status       HealthCheckStatus `json:"status"`
	Dependencies []Pulse           `json:"dependencies,omitempty"`
}

func NewPulse(component string) Pulse {
	return Pulse{
		Component:    component,
		Status:       StatusUnknown,
		Dependencies: nil,
	}
}

func ClearPulsers() {
	m.Lock()
	defer m.Unlock()

	registeredPulsers = make(map[string]Pulser)
}

func RegisterPulser(pulsers ...Pulser) {
	m.Lock()
	defer m.Unlock()

	for _, pulser := range pulsers {
		registeredPulsers[pulser.Component()] = pulser
	}
}

func HealthHandler(conf config.Heartbeat) func(*gin.Context) {
	return func(c *gin.Context) {
		withTimeout, cancel := context.WithTimeout(context.Background(), conf.Timeout.Read)
		defer cancel()

		pulsec := make(chan Pulse, len(registeredPulsers))

		for _, pulser := range registeredPulsers {
			go pulser.Pulse(withTimeout, pulsec)
		}

		dependencies := make([]Pulse, 0)

		finished := len(registeredPulsers) == 0

		for !finished {
			select {
			case pulse, ok := <-pulsec:
				if !ok {
					finished = true
				} else {
					dependencies = append(dependencies, pulse)

					if len(dependencies) == len(registeredPulsers) {
						finished = true
					}
				}
			}
		}

		response := NewHealthCheckResponse(conf.Info(), dependencies)

		c.JSON(response.HTTPStatus(), response)
	}
}
