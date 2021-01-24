package agora_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/errwrap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/wgentry22/agora"
	"github.com/wgentry22/agora/testutil"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var _ = Describe("Configuration", func() {

	Context("when using invalid values in callbacks", func() {
		It("should panic when path is not absolute", func() {
			defer ConfigurationRecover(agora.ErrConfigPathMustBeAbsolute)

			agora.New(agora.ConfigPath("some/relative/path"))
		})

		It("should panic when path is empty", func() {
			defer ConfigurationRecover(agora.ErrConfigPathMustBeAbsolute)

			agora.New(agora.ConfigPath(""))
		})

		It("should panic when ConfigFile callback uses a file type other than TOML", func() {
			defer ConfigurationRecover(agora.ErrConfigMustBeTOML)

			agora.New(agora.ConfigName("config.json"))
		})
	})

	Context("when using valid values in callbacks", func() {
		var (
			pgContainerArgs = testutil.TestContainerArgs{
				Image: "bitnami/postgresql",
				Env: map[string]string{
					"POSTGRES_USER":     "test",
					"POSTGRES_PASSWORD": "test",
					"POSTGRES_DB":       "test",
				},
				Port:       "5432/tcp",
				WaitForLog: "database system is ready to accept connections",
				AutoRemove: true,
			}
			pgContainer testcontainers.Container
			tcContext   context.Context
			port        nat.Port
		)

		It("should setup successfully", func() {
			path := fmt.Sprintf("%s/app.toml", os.TempDir())

			err := writeTCPortToFile(path, port)
			Expect(err).To(BeNil())

			app := agora.New(agora.ConfigPath(os.TempDir()))
			Expect(app).ToNot(BeNil())
		})

		BeforeEach(func() {
			ctx := context.Background()

			pgContainer = testutil.NewPGTestContainer(ctx, pgContainerArgs)

			var err error

			if port, err = testutil.PortFromContainer(ctx, pgContainerArgs, pgContainer); err != nil {
				errwrap.Walk(err, func(err error) {
					log.Printf("Error while getting port from TC: %s", err)
				})
				panic(err)
			}

			tcContext = ctx
		})

		AfterEach(func() {
			if err := pgContainer.Terminate(tcContext); err != nil {
				panic(err)
			}
		})
	})
})

func ConfigurationRecover(expectedErr error) {
	if r := recover(); r != nil {
		err, isErr := r.(error)
		if !isErr {
			Fail("Expected to panic when passing an invalid value to Configuration Callback")
		}

		Expect(err).To(Equal(expectedErr))
	}
}

func writeTCPortToFile(path string, port nat.Port) error {
	var buf bytes.Buffer

	tmpl := template.Must(template.New("test").Parse(`
[info]
name = "agora"
version = "1.2.3"
env = "qa"

[api]
port = 9123
pathPrefix = "prefix"

[logging]
level = "trace"
outputPaths = ["stdout"]
[logging.fields]
from = "toml"

[db]
vendor = "postgres"
user = "test"
password = "test"
host = "localhost"
port = {{.Port}}
name = "test"
[db.args]
sslmode = "disable"
`))

	data := struct{ Port int }{Port: port.Int()}

	if err := tmpl.ExecuteTemplate(&buf, "test", &data); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
