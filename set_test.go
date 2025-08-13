package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/qoinlyid/qore"
	"github.com/stretchr/testify/assert"
)

const (
	testPrefix      = "test"
	testKey         = "TestSet"
	testKeyRemember = "TestRemember"
	testValue       = "test_value"
)

const testValueRemember int = 15

func TestSet(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	keyExpected := testPrefix + DefaultKeySeparator + testKey

	// Set.
	set := i.Set(
		t.Context(),
		testKey,
	).SetPrefix(testPrefix).SetTTL(time.Second)
	assert.NotNil(t, set, "setter must not be nil")
	assert.Equal(t, keyExpected, set.key, fmt.Sprintf("Set key must be equal %s", keyExpected))
}

func TestSetPutDefaultTTL(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()
	ttl, err := i.Set(
		t.Context(),
		testKey,
	).SetPrefix(testPrefix).Put(testValue)
	assert.NoError(t, err, "Set Put must be no error")
	assert.Equal(t, DefaultTTL, ttl, fmt.Sprintf("Time-to-live value for setter must be %s", DefaultTTL))
}

func TestSetPutForever(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()
	ttl, err := i.Set(
		t.Context(),
		testKey,
	).SetPrefix(testPrefix).PutForever(testValue)
	assert.NoError(t, err, "Set PutForever must be no error")
	expectedTTL := time.Duration(0)
	assert.Equal(t, expectedTTL, ttl, fmt.Sprintf("Time-to-live value for setter must be %s", expectedTTL))
}

func TestRateLimitOnce(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	allowed, err := i.Set(t.Context(), "LimitKey").RateLimitOnce(time.Minute)
	assert.NoError(t, err, "Must be no error")
	assert.True(t, allowed, "LimitOnce must be allowed once perform for the first time")
}

func TestRateLimitOnceBlocked(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	allowed, err := i.Set(t.Context(), "LimitKeyBLocked").RateLimitOnce(time.Minute)
	assert.NoError(t, err, "Must be no error")
	assert.True(t, allowed, "LimitOnce must be allowed once perform for the first time")

	allowed, err = i.Set(t.Context(), "LimitKeyBLocked").RateLimitOnce(time.Minute)
	assert.NoError(t, err, "Must be no error")
	assert.False(t, allowed, "LimitOnce must be blocked")
}
