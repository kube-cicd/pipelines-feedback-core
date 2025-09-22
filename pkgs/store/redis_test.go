package store

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRedis_Flow(t *testing.T) {
	ctx := context.Background()

	//
	// 0. Setup Redis container
	//
	req := testcontainers.ContainerRequest{
		Image:        "quay.io/opstree/redis:v7.4.4",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.Nil(t, err)
	ip, ipErr := container.ContainerIP(ctx)
	assert.Nil(t, ipErr)

	adapter := NewRedis()
	os.Setenv("REDIS_HOST", ip+":6379")
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
