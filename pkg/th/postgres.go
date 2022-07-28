package th

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer(ctx context.Context) (string, testcontainers.Container, error) {
	port := "5432/tcp"
	dbname := "internal"
	credential := "postgres"
	dsnTemplate := "postgres://%s:%s@localhost:%s/%s?sslmode=disable"

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:13-alpine",
			ExposedPorts: []string{port},
			Env: map[string]string{
				"POSTGRES_USER":     credential,
				"POSTGRES_PASSWORD": credential,
				"POSTGRES_DB":       dbname,
				"POSTGRES_SSL_MODE": "disable",
			},
			Cmd: []string{
				"postgres", "-c", "fsync=off",
			},
			WaitingFor: wait.ForSQL(
				nat.Port(port),
				"postgres",
				func(p nat.Port) string {
					return fmt.Sprintf(dsnTemplate, credential, credential, p.Port(), dbname)
				},
			).Timeout(time.Second * 5),
		},
	})
	if err != nil {
		return "", container, err
	}

	p, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return "", container, fmt.Errorf("failed to get container external port: %w", err)
	}

	dbDSN := fmt.Sprintf(dsnTemplate, credential, credential, p.Port(), dbname)

	return dbDSN, container, nil
}
