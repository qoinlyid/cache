package cache

import (
	"fmt"
	"testing"

	"github.com/qoinlyid/qore"
	"github.com/stretchr/testify/assert"
)

func TestHasExists(t *testing.T) {
	// Put forever.
	TestSetPutForever(t)

	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()
	exists := i.Has(t.Context(), testKey, testPrefix)
	assert.True(t, exists, "Has return must be true")
}

func TestHasNotExists(t *testing.T) {
	// Put forever.
	TestSetPutForever(t)

	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()
	exists := i.Has(t.Context(), testKey)
	assert.False(t, exists, "Has return must be false")
}

func TestGetAllKeys(t *testing.T) {
	// Put forever.
	TestSetPutForever(t)

	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()
	keys, err := i.GetAllKeys(t.Context(), testPrefix)
	assert.NoError(t, err, "GetAllKeys must be no error")
	var found bool
	for _, key := range keys {
		if key.Key == testKey {
			found = true
			break
		}
	}
	assert.True(t, found, fmt.Sprintf("Keys must contains %s", testKey))
}
