package cache

import (
	"context"
	"time"

	"github.com/qoinlyid/qore"
)

type RememberFn func() (forever bool, val any, err error)

// getter is a method-chaining configuration struct for cache retrieve operations.
type getter struct {
	base

	// getFn is a closure function that called to retrieve cache from the backend.
	getFn func(get *getter, out any, rem ...RememberFn) error
}

// Get initializes a new getter instance for the given key & prefix (if any).
// If the provided context is nil, a new background context with
// a default timeout will be created. The cancel function is stored
// in the getter and will be called automatically when the final
// method (e.g., Pull) completes.
//
// This is the entry point for method chaining.
//
//	get := cache.Get(ctx, "myKey")
func (i *Instance) Get(ctx context.Context, key string, prefix ...string) *getter {
	// Value modifier.
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	}
	if len(prefix) > 0 {
		if !qore.ValidationIsEmpty(prefix[0]) {
			key = prefix[0] + DefaultKeySeparator + key
		}
	}

	// Return getter.
	return &getter{
		base: base{
			ctx:    ctx,
			cancel: cancel,
			key:    key,
		},
		getFn: i.get,
	}
}

// Pull retrieves item(s) from cache and parse it to the given output.
//
//	var out any
//	err := get.Pull(&out)
//	if err != nil {
//		log.Println(err)
//	}
func (g *getter) Pull(out any) error {
	return g.getFn(g, out)
}

// Remember retrieves item(s) from cache and parse it to the given output,
// but also store a default value if the requested item(s) does not exist.
//
// In the `RememberFn` first return value is bool to determines is the cache entry should be persist or not.
//
//	var out string
//	err := get.Remember(&out, func() (forever bool, val any, err error) {
//		return true, "myValue", nil
//	})
//	if err != nil {
//		log.Println(err)
//	}
func (g *getter) Remember(out any, rem RememberFn) error {
	return g.getFn(g, out, rem)
}
