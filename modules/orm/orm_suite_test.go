package orm_test

import (
  "context"
  "github.com/docker/go-connections/nat"
  "github.com/hashicorp/errwrap"
  "github.com/pelletier/go-toml"
  "github.com/testcontainers/testcontainers-go"
  "github.com/wgentry22/agora/testutil"
  "github.com/wgentry22/agora/types/config"
  "log"
  "testing"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

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

  _ = BeforeSuite(func() {
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

  _ = AfterSuite(func() {
    if err := pgContainer.Terminate(tcContext); err != nil {
      panic(err)
    }
  })
)

func TestOrm(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Orm Suite")
}

func mustParseConfig() config.DB {
  anonConfig := struct {
    Vendor   string                 `toml:"vendor"`
    User     string                 `toml:"user"`
    Password string                 `toml:"password"`
    Host     string                 `toml:"host"`
    Port     int                    `toml:"port"`
    DBName   string                 `toml:"name"`
    Args     map[string]interface{} `toml:"args"`
  }{
    Vendor:   "postgres",
    User:     pgContainerArgs.Env["POSTGRES_USER"],
    Password: pgContainerArgs.Env["POSTGRES_PASSWORD"],
    Host:     "localhost",
    Port:     port.Int(),
    DBName:   pgContainerArgs.Env["POSTGRES_DB"],
    Args: map[string]interface{}{
      "sslmode": "disable",
    },
  }

  // Hacky, but effective workaround to populate config.DB since fields are not exported
  conf := new(config.DB)

  data, err := toml.Marshal(&anonConfig)
  if err != nil {
    panic(err)
  }

  var dataMap map[string]interface{}

  err = toml.Unmarshal(data, &dataMap)
  if err != nil {
    panic(err)
  }

  if err = conf.UnmarshalTOML(dataMap); err != nil {
    panic(err)
  }

  return *conf
}
