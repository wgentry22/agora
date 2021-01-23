package testutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/errwrap"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
)

var (
	ErrFailedToStartTestcontainer = func(args TestContainerArgs) error {
		return fmt.Errorf("failed to start testcontainer `%s`", args.Image)
	}
	ErrUnexpectedPortFormat = errors.New("unexpected port format")
	expectedPortPartsLen    = 2
)

type TestContainerArgs struct {
	Image      string
	Env        map[string]string
	Port       string
	WaitForLog string
	AutoRemove bool
}

func NewPGTestContainer(ctx context.Context, args TestContainerArgs) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:           args.Image,
		Env:             args.Env,
		ExposedPorts:    []string{args.Port},
		WaitingFor:      wait.ForLog(args.WaitForLog),
		AutoRemove:      args.AutoRemove,
		AlwaysPullImage: true,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(errwrap.Wrap(ErrFailedToStartTestcontainer(args), err))
	}

	return container
}

func PortFromContainer(ctx context.Context, args TestContainerArgs, container testcontainers.Container) (port nat.Port, err error) { //nolint:lll
	portParts := strings.Split(args.Port, "/")
	if len(portParts) != expectedPortPartsLen {
		err = ErrUnexpectedPortFormat

		return
	}

	port, err = nat.NewPort(portParts[1], portParts[0])
	if err != nil {
		err = errwrap.Wrap(ErrUnexpectedPortFormat, err)

		return
	}

	port, err = container.MappedPort(ctx, port)

	return
}
