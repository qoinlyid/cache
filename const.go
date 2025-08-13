package cache

import "time"

// String.
const (
	DefaultNameSpace    = "cache-app"
	DefaultKeySeparator = ":"
)

// Numeric
const (
	DefaultTTL = time.Minute
)

// Rate Limit.
const (
	KeyRateLimit = "rate-limit"
)
