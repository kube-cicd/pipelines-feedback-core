package store

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestRedis_Flow(t *testing.T) {
	ctx := context.Background()

	//
	// 0. Setup Redis container
	//
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/mirrorshub/docker/redis:7.0.7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.Nil(t, err)
	_, ipErr := container.ContainerIP(ctx)
	assert.Nil(t, ipErr)

	adapter := NewRedis()
	adapter.Initialize()

	//
	// Test
	//

	// CASE: When no key submitted - then report "No such key"
	_, err = adapter.Get("non-existing-key")
	assert.Equal(t, "No such key", err.Error())

	// CASE: Set() and Get()
	assert.Nil(t, adapter.Set("Hello/World", "Bread", 100))
	get, getErr := adapter.Get("Hello/World")
	assert.Equal(t, "Bread", get)
	assert.Nil(t, getErr)
}
