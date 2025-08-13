package cache

import (
	"fmt"
	"testing"

	"github.com/qoinlyid/qore"
	"github.com/stretchr/testify/assert"
)

func TestNewWithoutConfig(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	assert.Equal(
		t, defaultConfig.DependencyPriority, i.cfg.DependencyPriority,
		fmt.Sprintf("Value for DependencyPriority of unloaded config must be %d", defaultConfig.DependencyPriority),
	)
	assert.Empty(t, i.cfg.Addresses, "Value of Addresses of unloaded config must be empty")
}

func TestNewWithConfig(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	ok := i.cfg.Addresses != "" || i.cfg.SentinelAddresses != ""
	assert.True(t, ok, "Either c.Addresses or c.SentinelAddresses must not be empty")
}

func TestOpen(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	err := i.Open()
	defer i.Close()
	assert.NoError(t, err, "Open must be no error")
}

func TestClient(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	err := i.Open()
	defer i.Close()
	assert.NoError(t, err, "Open must be no error")

	client, err := i.Client()
	assert.NoError(t, err, "Must be no error while perform opening was successfull")
	assert.NotNil(t, client, "Client must be not nil")
}

func TestHealthCheckHealthy(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	err := i.Open()
	assert.NoError(t, err)
	stats := i.HealthCheck(t.Context())
	assert.NotNil(t, stats, "DependencyStats value must be not nil")
	assert.NotEmpty(t, stats.PINGLatencyMillis, "PINGLatencyMillis value must be not empty")
}

func TestClose(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	err := i.Open()
	assert.NoError(t, err, "Open must be no error")
	err = i.Close()
	assert.NoError(t, err, "Close must be no error")
}
