package cache

import (
	"os"
	"testing"

	"github.com/qoinlyid/qore"
	"github.com/stretchr/testify/assert"
)

func configLoaderTest(file string) *Config {
	os.Setenv(qore.CONFIG_USED_KEY, file)
	return loadConfig()
}

func TestLoadConfigDotEnv(t *testing.T) {
	c := configLoaderTest("./.env")
	assert.NotNil(t, c)
	ok := c.Addresses != "" || c.SentinelAddresses != ""
	assert.True(t, ok, "Either c.Addresses or c.SentinelAddresses must not be empty")
}

func TestLoadConfigJSON(t *testing.T) {
	c := configLoaderTest("./.env.json")
	assert.NotNil(t, c)
	ok := c.Addresses != "" || c.SentinelAddresses != ""
	assert.True(t, ok, "Either c.Addresses or c.SentinelAddresses must not be empty")
}

func TestLoadConfigTOML(t *testing.T) {
	c := configLoaderTest("./.env.toml")
	assert.NotNil(t, c)
	ok := c.Addresses != "" || c.SentinelAddresses != ""
	assert.True(t, ok, "Either c.Addresses or c.SentinelAddresses must not be empty")
}

func TestLoadConfigYAML(t *testing.T) {
	c := configLoaderTest("./.env.yaml")
	assert.NotNil(t, c)
	ok := c.Addresses != "" || c.SentinelAddresses != ""
	assert.True(t, ok, "Either c.Addresses or c.SentinelAddresses must not be empty")
}
