package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheGenDefaultName(t *testing.T) {
	i := New()
	assert.Equal(t, "Cache", i.Name(), "Name should be 'Cache'")
}

func TestCacheGenDefaultPriority(t *testing.T) {
	i := New()
	assert.Equal(
		t,
		defaultConfig.DependencyPriority,
		i.Priority(),
		fmt.Sprintf("Default DependencyPriority should be %d", defaultConfig.DependencyPriority),
	)
}
