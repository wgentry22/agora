package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pelletier/go-toml"
	"github.com/wgentry22/agora/types/config"
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
			}))

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
`)
		)

		It("should be parsed correctly", func() {
			err := toml.Unmarshal(tomlData, &app)
			Expect(err).To(BeNil())

			Expect(app.Info()).To(Equal(config.Info{
				Name: "agora",
				Version: config.SemanticVersion{
					Major: 1,
					Minor: 2,
					Patch: 3,
				},
				Env: config.QualityAssurance,
			}))

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
			}))

			connStr := config.ConnectionString(app.DB())
			expectedConnStr := "user=test password=test host=localhost port=6000 dbname=test sslmode=disable"

			Expect(expectedConnStr).To(Equal(connStr))
		})
	})
})
