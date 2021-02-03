package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pelletier/go-toml"
	"github.com/wgentry22/agora/types/config"
	"time"
)

var _ = Describe("Parsing config.Application", func() {

	Context("when file is TOML, but empty", func() {
		var (
			app      config.Application
			tomlData = []byte(``)
		)

		It("should resolve defaults", func() {
			err := toml.Unmarshal(tomlData, &app)
			Expect(err).To(BeNil())

			Expect(app.Info()).To(Equal(config.Info{
				Name: "agora-app",
				Version: config.SemanticVersion{
					Major: 0,
					Minor: 0,
					Patch: 1,
				},
				Env: config.Development,
			}))

			Expect(app.Logging()).To(Equal(config.Logging{
				Level:       "debug",
				OutputPaths: []string{"stdout"},
				Fields: map[string]interface{}{
					"name":    "agora-app",
					"version": "0.0.1",
					"env":     "dev",
				},
			}))

			Expect(app.API()).To(Equal(config.API{
				Port:       8123,
				PathPrefix: "/v1",
			}.WithDefaultInfo()))

			connStr := config.ConnectionString(app.DB())
			Expect(connStr).To(Equal("unknown"))
		})
	})

	Context("when file is TOML", func() {
		var (
			app      config.Application
			tomlData = []byte(`
[info]
name = "agora"
version = "1.2.3"
env = "qa"

[api]
port = 9123
pathPrefix = "prefix"
[api.timeout]
read = 5678
write = 1234
[api.cors]
allow-origins = ["http://localhost:1234"]
allow-methods = ["GET","PUT","POST","PATCH","DELETE"]
allow-headers = ["Access-Control-Allow-Origin"]
expose-headers = ["Access-Control-Allow-Origin"]
allow-credentials = true

[heartbeat]
pathPrefix = "/ekg"
[heartbeat.timeout]
read = 2345
write = 3456

[logging]
level = "trace"
outputPaths = ["stdout", "/var/log/app.log"]
[logging.fields]
from = "toml"

[db]
vendor = "postgres"
user = "test"
password = "test"
host = "localhost"
port = 6000
name = "test"
[db.args]
sslmode = "disable"

[broker]
id = "some-id"
vendor = "kafka"
role = "producer"
servers = ["localhost:1234", "localhost:2345"]
timeout = 250
buffer_size = 1000
[broker.args]
"auto.group.offset" = "smallest"
`)
		)

		It("should be parsed correctly", func() {
			err := toml.Unmarshal(tomlData, &app)
			Expect(err).To(BeNil())

			expectedInfo := config.Info{
				Name: "agora",
				Version: config.SemanticVersion{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
				Env: config.QualityAssurance,
			}

			Expect(app.Info()).To(Equal(expectedInfo))

			Expect(app.Logging()).To(Equal(config.Logging{
				Level:       "trace",
				OutputPaths: []string{"stdout", "/var/log/app.log"},
				Fields: map[string]interface{}{
					"name":    "agora",
					"version": "1.2.3",
					"env":     "qa",
					"from":    "toml",
				},
			}))

			Expect(app.API()).To(Equal(config.API{
				Port:       9123,
				PathPrefix: "/prefix",
				Timeout: config.TimeoutOptions{
					Read:  5678 * time.Millisecond,
					Write: 1234 * time.Millisecond,
				},
				Cors: config.CORS{
					AllowOrigins:     []string{"http://localhost:1234"},
					AllowMethods:     []string{"GET", "PUT", "POST", "PATCH", "DELETE"},
					AllowHeaders:     []string{"Access-Control-Allow-Origin"},
					ExposeHeaders:    []string{"Access-Control-Allow-Origin"},
					AllowCredentials: true,
				},
			}.WithInfo(expectedInfo)))

			Expect(app.Heartbeat()).To(Equal(config.Heartbeat{
				PathPrefix: "/ekg",
				Timeout: config.TimeoutOptions{
					Read:  2345 * time.Millisecond,
					Write: 3456 * time.Millisecond,
				},
			}.WithInfo(expectedInfo)))

			connStr := config.ConnectionString(app.DB())
			expectedConnStr := "user=test password=test host=localhost port=6000 dbname=test sslmode=disable"

			Expect(expectedConnStr).To(Equal(connStr))

			Expect(app.Broker()).To(Equal(config.Broker{
				ID:         "some-id",
				Servers:    []string{"localhost:1234", "localhost:2345"},
				BufferSize: 1000,
				Vendor:     config.BrokerVendorKafka,
				Role:       config.BrokerRoleProducer,
				Timeout:    250,
				Args: map[string]interface{}{
					"auto.group.offset": "smallest",
				},
			}))
		})
	})
})
