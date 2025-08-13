package cache

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/qoinlyid/qore"
	"github.com/spf13/viper"
)

// Config defines Cache config.
type Config struct {
	// DependencyPriority defines priority of cache dependency.
	DependencyPriority int `json:"CACHE_DEPENDENCY_PRIORITY" mapstructure:"CACHE_DEPENDENCY_PRIORITY"`

	// Namespace defines cache key prefix that always be used.
	Namespace string `json:"CACHE_NAMESPACE" mapstructure:"CACHE_NAMESPACE"`

	// DB defines redis logical DB that used by cache dependency. For clustering mode, this will be ignored.
	DB int `json:"CACHE_DB" mapstructure:"CACHE_DB"`

	// Username defines redis username.
	Username string `json:"CACHE_USERNAME" mapstructure:"CACHE_USERNAME"`

	// Password defines redis password.
	Password string `json:"CACHE_PASSWORD" mapstructure:"CACHE_PASSWORD"`

	// Addresses defines redis addresses. Use comma separated to set redis cluster.
	// Even cluster single-endpoint give comma as a suffix.
	Addresses string `json:"CACHE_ADDRESSES" mapstructure:"CACHE_ADDRESSES"`

	// SentinelAddresses defines redis sentinel addresses. Use comma separated to set redis cluster.
	SentinelAddresses string `json:"CACHE_SENTINEL_ADDRESSES" mapstructure:"CACHE_SENTINEL_ADDRESSES"`

	// SentinelMaster defines redis sentinel master name.
	SentinelMaster string `json:"CACHE_SENTINEL_MASTER" mapstructure:"CACHE_SENTINEL_MASTER"`

	// SentinelUsername defines redis sentinel username, it can be same or different with redis username.
	SentinelUsername string `json:"CACHE_SENTINEL_USERNAME" mapstructure:"CACHE_SENTINEL_USERNAME"`

	// SentinelPassword defines redis sentinel password, it can be same or different with redis password.
	SentinelPassword string `json:"CACHE_SENTINEL_PASSWORD" mapstructure:"CACHE_SENTINEL_PASSWORD"`

	// SentinelCluster defines redis sentinel backend is using cluster mode.
	SentinelCluster bool `json:"CACHE_SENTINEL_CLUSTER" mapstructure:"CACHE_SENTINEL_CLUSTER"`
}

// Default config.
var defaultConfig = &Config{
	DependencyPriority: 10,
}

// Load config.
func loadConfig() *Config {
	var e error
	config := defaultConfig

	// Get used config from OS env.
	configSource := os.Getenv(qore.CONFIG_USED_KEY)
	if qore.ValidationIsEmpty(configSource) {
		configSource = "OS"
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	switch strings.ToUpper(configSource) {
	case "OS":
		if err := viper.Unmarshal(&config); err != nil {
			e = errors.Join(fmt.Errorf("failed to parse OS env value to config: %w", err))
		}
	default:
		ext := strings.ToLower(filepath.Ext(configSource))
		switch ext {
		case ".env":
			viper.SetConfigFile(configSource)
			viper.SetConfigType("env")
			if err := viper.ReadInConfig(); err != nil {
				e = errors.Join(fmt.Errorf("failed to read env file %s: %w", configSource, err))
			} else {
				if err := viper.Unmarshal(&config); err != nil {
					e = errors.Join(fmt.Errorf("failed to parse env file %s value to config: %w", configSource, err))
				}
			}
		case ".json", ".yml", ".yaml", ".toml":
			viper.SetConfigFile(configSource)
			if err := viper.ReadInConfig(); err != nil {
				e = errors.Join(fmt.Errorf("failed to read config file %s: %w", configSource, err))
			} else {
				if err := viper.Unmarshal(&config); err != nil {
					e = errors.Join(fmt.Errorf("failed to parse config file %s value to config: %w", configSource, err))
				}
			}
		}
	}
	if e != nil {
		log.Printf("dependency config - failed to load config: %s\n", e.Error())
	}

	// Config value modifier.
	if config.DependencyPriority == 0 {
		config.DependencyPriority = defaultConfig.DependencyPriority
	}
	if qore.ValidationIsEmpty(config.Namespace) {
		config.Namespace = DefaultNameSpace
	}
	return config
}
