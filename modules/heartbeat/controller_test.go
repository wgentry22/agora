package heartbeat_test

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/api"
	"github.com/wgentry22/agora/modules/heartbeat"
	"github.com/wgentry22/agora/types/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("HeartbeatController", func() {
	var (
		appInfo = config.Info{
			Name: "heartbeat",
			Version: config.SemanticVersion{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			Env: config.Production,
		}

		apiConfig = config.API{
			Port:       8123,
			PathPrefix: "/app",
			Timeout: config.TimeoutOptions{
				Read:  5 * time.Second,
				Write: 5 * time.Second,
			},
		}.WithInfo(appInfo)

		router          api.Router
		heartbeatConfig config.Heartbeat
	)

	BeforeEach(func() {
		router = api.NewRouter(apiConfig)
	})

	Context("when registering with api.Router", func() {
		var ts *httptest.Server

		BeforeEach(func() {
			heartbeatConfig = config.Heartbeat{
				PathPrefix: "/heartbeat",
				Timeout: config.TimeoutOptions{
					Read:  5 * time.Second,
					Write: 5 * time.Second,
				},
			}.WithInfo(appInfo)

			testPulser := &TestPulser{}

			heartbeat.RegisterPulser(testPulser)

			router.Register(heartbeat.NewHeartbeatController(heartbeatConfig))

			ts = httptest.NewServer(router.Handler())
		})

		AfterEach(func() {
			heartbeat.ClearPulsers()
			ts.Close()
		})

		It("should serve health checks without error", func() {
			res, err := http.Get(fmt.Sprintf("%s%s%s/health", ts.URL, apiConfig.PathPrefix, heartbeatConfig.PathPrefix))
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(http.StatusOK))

			data, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			var response heartbeat.HealthCheckResponse
			err = json.Unmarshal(data, &response)
			Expect(err).To(BeNil())

			Expect(response.Status).To(Equal(heartbeat.StatusOK))
			Expect(response.App).To(Equal(appInfo.Name))
			Expect(response.Version).To(Equal(appInfo.Version.String()))
			Expect(response.Env).To(Equal(appInfo.Env.String()))
			Expect(response.Dependencies).To(HaveLen(1))
			Expect(response.Dependencies[0].Status).To(Equal(heartbeat.StatusOK))
			Expect(response.Dependencies[0].Dependencies).To(HaveLen(0))
		})
	})

	Context("when serving responses", func() {
		It("should return 503 SERVICE UNAVAILABLE if at least one dependency returns StatusCritical", func() {
			pulses := []heartbeat.Pulse{
				{
					Component:    "test",
					Status:       heartbeat.StatusCritical,
					Dependencies: nil,
				},
				{
					Component:    "test-warn",
					Status:       heartbeat.StatusWarn,
					Dependencies: nil,
				},
				{
					Component:    "test-ok",
					Status:       heartbeat.StatusOK,
					Dependencies: nil,
				},
			}

			response := heartbeat.NewHealthCheckResponse(appInfo, pulses)

			Expect(response.HTTPStatus()).To(Equal(http.StatusServiceUnavailable))
			Expect(response.Status).To(Equal(heartbeat.StatusCritical))
		})

		It("should return 424 FAILED DEPENDENCY if at least one dependency returns StatusWarn and no StatusCritical", func() {
			pulses := []heartbeat.Pulse{
				{
					Component:    "test",
					Status:       heartbeat.StatusWarn,
					Dependencies: nil,
				},
				{
					Component:    "test-ok",
					Status:       heartbeat.StatusOK,
					Dependencies: nil,
				},
			}

			response := heartbeat.NewHealthCheckResponse(appInfo, pulses)

			Expect(response.HTTPStatus()).To(Equal(http.StatusFailedDependency))
			Expect(response.Status).To(Equal(heartbeat.StatusWarn))
		})

		It("should return 200 OK if one dependency returns StatusOK", func() {
			pulses := []heartbeat.Pulse{
				{
					Component:    "test",
					Status:       heartbeat.StatusOK,
					Dependencies: nil,
				},
			}

			response := heartbeat.NewHealthCheckResponse(appInfo, pulses)

			Expect(response.HTTPStatus()).To(Equal(http.StatusOK))
			Expect(response.Status).To(Equal(heartbeat.StatusOK))
		})
	})
})

type TestPulser struct{}

func (tp *TestPulser) Component() string {
	return "test-pulser"
}

func (tp *TestPulser) Pulse(ctx context.Context, pulsec chan<- heartbeat.Pulse) {
	pulse := heartbeat.Pulse{
		Component:    tp.Component(),
		Status:       heartbeat.StatusOK,
		Dependencies: nil,
	}

	pulsec <- pulse
}
