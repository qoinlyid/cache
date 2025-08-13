package cache

import "errors"

var (
	ErrClientNil        = errors.New("redis client is null")
	ErrClientNotCluster = errors.New("redis set to cluster mode, but unfortunately the client is not cluster client")
	ErrEmptyKey         = errors.New("cache key cannot be empty")
	ErrEmptyPrefix      = errors.New("prefix cannot be empty")
	ErrOutNonPointer    = errors.New("out type non-pointer")
)
