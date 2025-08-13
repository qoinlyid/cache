package cache

import (
	"context"
	"time"

	"github.com/qoinlyid/qore"
	"github.com/redis/go-redis/v9"
)

// Instance defines Cache dependency singleton.
type Instance struct {
	// Define dependency singleton here.
	client redis.UniversalClient

	// Private field.
	cfg        *Config
	startTime  time.Time
	clustering bool
	*instanceGen
}

// New creates singleton dependency instance.
func New() *Instance {
	config := loadConfig()
	instance := &Instance{
		cfg:         config,
		instanceGen: &instanceGen{priority: config.DependencyPriority},
	}
	return instance
}

// HealthCheck returns statistics for dependency health check.
func (i *Instance) HealthCheck(ctx context.Context) *qore.DependencyStats {
	uptime := time.Since(i.startTime)
	stats := &qore.DependencyStats{
		UptimeSeconds: uptime.Seconds(),
		UptimeHuman:   uptime.String(),
	}
	if i.client == nil {
		return stats
	}

	start := time.Now()
	ping, err := i.client.Ping(ctx).Result()
	latency := time.Since(start)
	if err != nil {
		stats.PINGResponse = err.Error()
		return stats
	}
	stats.PINGLatencyMillis = latency.Milliseconds()
	stats.PINGLatencyHuman = latency.String()
	stats.PINGResponse = ping
	return stats
}

// Open an backend connection or construct the dependency.
func (i *Instance) Open() error {
	// Open connection.
	if err := i.open(); err != nil {
		return err
	}

	// Set another instance field.
	i.startTime = time.Now()

	// Return.
	return nil
}

// Close an backend connection or destruct the dependency.
func (i *Instance) Close() error {
	// Close connection.
	if i.client == nil {
		return nil
	}
	return i.client.Close()
}

// Client returns redis client interface.
func (i *Instance) Client() (redis.UniversalClient, error) {
	if err := i.validateClient(); err != nil {
		return nil, err
	}
	return i.client, nil
}
