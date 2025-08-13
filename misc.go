package cache

import (
	"context"
	"strings"
	"time"

	"github.com/qoinlyid/qore"
	"github.com/redis/go-redis/v9"
)

// Has checks whether the specified key exists in cache.
// An optional prefix can be provided, which is prepended to the key along with
// the default key separator.
// If the context is nil, a 1-second timeout context is used.
// Returns true if the key exists, false otherwise.
//
//	exists := cache.Has(ctx, "myKey")
func (i *Instance) Has(ctx context.Context, key string, prefix ...string) bool {
	// Validate.
	if err := i.validateClient(); err != nil {
		return false
	}
	if ctx == nil {
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ctx = c
	}
	if len(prefix) > 0 {
		if !qore.ValidationIsEmpty(prefix[0]) {
			key = prefix[0] + DefaultKeySeparator + key
		}
	}
	key = i.cfg.Namespace + DefaultKeySeparator + key

	// Exec.
	res, err := i.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return res > 0
}

type Keyer struct {
	Namespace string
	Prefix    string
	Key       string
}

// GetAllKeys retrieves all Redis keys that match the specified prefix.
// The prefix is normalized to ensure it ends with the default key separator
// before performing a SCAN operation.
//
// A non-empty prefix is required; otherwise, ErrEmptyPrefix is returned.
// If the provided context is nil, a new context with a 1-second timeout is used.
// Returns a slice of matching keys or an error if the SCAN operation fails.
//
//	keys, err := cache.GetAllKeys(ctx, "myPrefix")
//	if err != nil {
//		log.Println(err)
//	}
func (i *Instance) GetAllKeys(ctx context.Context, prefix string) (keys []Keyer, err error) {
	// Validate.
	if e := i.validateClient(); e != nil {
		err = e
		return
	}
	if qore.ValidationIsEmpty(prefix) {
		err = ErrEmptyPrefix
		return
	}
	if ctx == nil {
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ctx = c
	}
	match := i.cfg.Namespace + DefaultKeySeparator + strings.TrimSuffix(prefix, "*") + DefaultKeySeparator + "*"

	// Perform SCAN in the each cluster node.
	if i.clustering {
		clusterClient, ok := i.client.(*redis.ClusterClient)
		if !ok {
			err = ErrClientNotCluster
			return
		}
		clusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			return scan(client, ctx, match, &keys)
		})
	}

	// Perform SCAN.
	err = scan(i.client, ctx, match, &keys)
	return
}

func scan(rdb redis.UniversalClient, ctx context.Context, match string, founds *[]Keyer) error {
	iter := rdb.Scan(ctx, 0, match, 0).Iterator()
	for iter.Next(ctx) {
		vals := strings.Split(iter.Val(), DefaultKeySeparator)
		reverseStrings(vals)
		keyer := Keyer{}
		for i, v := range vals {
			switch i {
			case 0:
				keyer.Key = v
			case 1:
				keyer.Prefix = v
			default:
				keyer.Namespace = v
			}
		}
		*founds = append(*founds, keyer)
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	return nil
}
