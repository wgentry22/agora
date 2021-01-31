package agora_test

import (
	"context"
	"log"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/errwrap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/wgentry22/agora/testutil"
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

func TestAgora(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agora Suite")
}
