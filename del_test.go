package cache

import (
	"fmt"
	"testing"

	"github.com/qoinlyid/qore"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	expectedKey := testPrefix + DefaultKeySeparator + testKey
	del := i.Delete(t.Context(), testKey, testPrefix)
	assert.NotNil(t, del, "Delete must be not nil")
	assert.Equal(t, expectedKey, del.key, fmt.Sprintf("Key value should be %s", expectedKey))
}

func TestDeletePerformWrongKey(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	count, err := i.Delete(t.Context(), testKey).Perform()
	assert.NoError(t, err, "Must be no error even if the key is wrong")
	assert.Equal(t, int64(0), count, "The returned count should be 0")
}

func TestDeletePerform(t *testing.T) {
	// Put forever.
	TestSetPutForever(t)

	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	count, err := i.Delete(t.Context(), testKey, testPrefix).Perform()
	assert.NoError(t, err, "Must be no error")
	assert.Greater(t, count, int64(0), "The returned count should be greather than 0")
}
