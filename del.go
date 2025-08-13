package cache

import (
	"context"
	"time"

	"github.com/qoinlyid/qore"
)

// getter is a method-chaining configuration struct for cache delete operations.
type deleter struct {
	base

	// delFn is a closure function that called to delete cache from the backend.
	delFn func(del *deleter) (int64, error)
}

// Delete initializes a new deleter instance for the given key & prefix (if any).
// If the provided context is nil, a new background context with
// a default timeout will be created. The cancel function is stored
// in the getter and will be called automatically when the final
// method (e.g., Perform) completes.
//
// This is the entry point for method chaining.
//
//	del := cache.Delete(ctx, "myKey")
func (i *Instance) Delete(ctx context.Context, key string, prefix ...string) *deleter {
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

	// Return deleter.
	return &deleter{
		base: base{
			ctx:    ctx,
			cancel: cancel,
			key:    key,
		},
		delFn: i.del,
	}
}

// Perform delete entry from the cache.
//
//	count, err := get.Perform(&out)
//	if err != nil {
//		log.Println(err)
//	}
//	log.Println(count)
func (del *deleter) Perform() (int64, error) {
	return del.delFn(del)
}
