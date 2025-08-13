package cache

import (
	"context"
	"time"

	"github.com/qoinlyid/qore"
	"github.com/redis/go-redis/v9"
)

// setter is a method-chaining configuration struct for cache store operations.
type setter struct {
	base
	ttl time.Duration

	// setFn is a closure function that called to stores cache in the backend.
	setFn func(set *setter, val any, borrow ...bool) (redis.UniversalClient, error)
}

func (s *setter) cleanup() {
	if s.cancel != nil {
		s.cancel()
	}
	// Zero out all fields to help GC or prepare for reuse
	*s = setter{}
}

// Set initializes a new setter instance for the given key.
// If the provided context is nil, a new background context with
// a default timeout will be created. The cancel function is stored
// in the setter and will be called automatically when the final
// method (e.g., Put) completes.
//
// This is the entry point for method chaining.
//
//	cache.Set(ctx, "myKey").
//		SetPrefix("user").
//		SetTTL(5 * time.Minute).
//		Put("value")
func (i *Instance) Set(ctx context.Context, key string) *setter {
	// Value modifier.
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	}

	// Return setter.
	return &setter{
		base: base{
			ctx:    ctx,
			cancel: cancel,
			key:    key,
		},
		setFn: i.set,
	}
}

// SetPrefix sets the key prefix for the cache entry.
// This is optional and is useful for namespacing keys.
//
//	s.SetPrefix("session")
func (s *setter) SetPrefix(prefix string) *setter {
	if !qore.ValidationIsEmpty(prefix) {
		s.key = prefix + DefaultKeySeparator + s.key
	}
	return s
}

// SetTTL sets the time-to-live (TTL) duration for the cache entry.
// If TTL is <= 0, the default value will be used is 1 minutes.
//
//	s.SetTTL(10 * time.Second)
func (s *setter) SetTTL(ttl time.Duration) *setter { s.ttl = ttl; return s }

// Put stores the given value in the cache using the configured
// context, prefix, key, and TTL. This is the final method in the
// method-chaining sequence. Once executed, the associated cancel
// function (if any) will be called to release context resources.
//
//	err := s.Put("value")
//	if err != nil {
//		log.Println(err)
//	}
func (s *setter) Put(value any) (ttl time.Duration, err error) {
	if s.ttl <= 0 {
		s = s.SetTTL(DefaultTTL)
	}
	ttl = s.ttl
	_, err = s.setFn(s, value)
	return
}

// PutForever stores the given value in the cache without an expiration time.
// This method is similar to Put but overrides the TTL (time-to-live) value
// to zero, indicating that the cache entry should persist indefinitely
// (depending on the cache backend's configuration and eviction policy).
//
// This is a terminal method in the method-chaining sequence, meaning it should
// be called after all desired configuration methods (e.g., SetPrefix, SetTTL).
//
//	err := s.PutForever("value")
//	if err != nil {
//		log.Println(err)
//	}
func (s *setter) PutForever(value any) (ttl time.Duration, err error) {
	s.ttl = 0
	ttl = s.ttl
	_, err = s.setFn(s, value)
	return
}

// RateLimitOnce attempts to set a rate limit for the given period.
// It returns `allowed = true` if the key was successfully set (meaning the action is permitted),
// and `allowed = false` if the key already exists (meaning the action is rate-limited).
//
// The TTL (time-to-live) is set to the specified period, ensuring the limit
// only applies once within that time window.
//
//	allowed, err := s.RateLimitOnce(time.Minute)
//	if err != nil {
//		log.Println(err)
//	}
//	if !allowed {
//		log.Println("blocked, already reach the limit!")
//	}
func (s *setter) RateLimitOnce(period time.Duration) (allowed bool, err error) {
	if !qore.ValidationIsEmpty(s.key) {
		s = s.SetPrefix(KeyRateLimit)
	}
	if period > 0 {
		s = s.SetTTL(period)
	}

	// Borrow client from instance.
	client, e := s.setFn(s, nil, true)
	defer s.cleanup()
	if e != nil {
		err = e
		return
	}

	// Perform SET with NX (only if key doesn't exist) and EX (expire after TTL).
	return client.SetNX(s.ctx, s.key, 1, s.ttl).Result()
}
