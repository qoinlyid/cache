package cache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/qoinlyid/qore"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	expectedKey := testPrefix + DefaultKeySeparator + testKey
	get := i.Get(t.Context(), testKey, testPrefix)
	assert.NotNil(t, get, "Get must be not nil")
	assert.Equal(t, expectedKey, get.key, fmt.Sprintf("Key value should be %s", expectedKey))
}

func TestGetPull(t *testing.T) {
	// Put forever.
	TestSetPutForever(t)

	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	var val string
	err := i.Get(t.Context(), testKey, testPrefix).Pull(&val)
	assert.NoError(t, err, "Get Pull must be no error")
	assert.Equal(t, testValue, val, fmt.Sprintf("Value should be %s", testKey))
}

func TestGetPullTargetNotPointer(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	var val string
	err := i.Get(t.Context(), testKey, testPrefix).Pull(val)
	assert.ErrorIs(t, err, ErrOutNonPointer, fmt.Sprintf("Error should be wrapped with %s", ErrOutNonPointer))
}

func TestGetRemember(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	var val int
	err := i.Get(
		t.Context(), testKeyRemember, testPrefix,
	).Remember(&val, func() (forever bool, val any, err error) {
		return false, testValueRemember, nil
	})
	assert.NoError(t, err, "Get Remember must be no error")
	assert.Equal(t, testValueRemember, val, fmt.Sprintf("Value should be %d", testValueRemember))
}

func TestGetRememberClosureError(t *testing.T) {
	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	testKeyRemember := testKeyRemember + "Failed"
	var val int
	err := i.Get(
		t.Context(), testKeyRemember, testPrefix,
	).Remember(&val, func() (forever bool, val any, err error) {
		return false, nil, errors.New("closure error")
	})
	assert.NotNil(t, err, "Should be return error")
	assert.Equal(t, 0, val, fmt.Sprintf("Value should be %d due closure error", 0))
}

func TestGetRememberClosureSkip(t *testing.T) {
	// Put forever.
	TestSetPutForever(t)

	t.Setenv(qore.CONFIG_USED_KEY, "./.env")
	i := New()
	i.Open()
	defer i.Close()

	var val string
	err := i.Get(
		t.Context(), testKey, testPrefix,
	).Remember(&val, func() (forever bool, val any, err error) {
		return false, nil, errors.New("closure error")
	})
	assert.NoError(t, err, "Must be no error")
	assert.Equal(t, testValue, val, fmt.Sprintf("Value should be %s even closure returning error", testValue))
}
