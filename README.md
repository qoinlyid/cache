# Cache
Simple Redis-backed cache library for Go, fully compatible with the Qore toolkit.

[![Go Report Card](https://goreportcard.com/badge/github.com/qoinlyid/qore)](https://goreportcard.com/report/github.com/qoinlyid/qore)

## Features

- **Redis Support**: Standalone, Cluster, and Sentinel modes
- **Fluent API**: Method chaining for intuitive cache operations
- **Dependency Management**: Implements `qore.Dependency` interface
- **Health Checks**: Built-in health monitoring with ping latency
- **Multiple Config Sources**: Support for environment variables, JSON, YAML, TOML, and .env files
- **Type Safety**: Generic encoding/decoding with support for primitives and complex types
- **Rate Limiting**: Built-in rate limiting functionality
- **Remember Pattern**: Cache-aside pattern with automatic fallback
- **Context Support**: Full context cancellation and timeout support

## Installation

```bash
go get github.com/qoinlyid/cache@latest
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/qoinlyid/cache"
)

func main() {
    // Create cache instance
    cache := cache.New()
    
    // Open connection
    if err := cache.Open(); err != nil {
        log.Fatal(err)
    }
    defer cache.Close()
    
    ctx := context.Background()
    
    // Store data
    _, err := cache.Set(ctx, "user:123").
        SetPrefix("session").
        SetTTL(5 * time.Minute).
        Put("user_data")
    if err != nil {
        log.Fatal(err)
    }
    
    // Retrieve data
    var result string
    err = cache.Get(ctx, "user:123", "session").Pull(&result)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Retrieved:", result)
}
```

## Configuration

### Environment Variables

Set these environment variables for basic configuration:

```bash
CACHE_ADDRESSES="localhost:6379"
CACHE_PASSWORD="your_password"
CACHE_DB=0
CACHE_NAMESPACE="myapp"
```

### Configuration Files

You can use `OS environment variables`, dotenv file `.env`, `.json` file, `.toml` file or `.yaml` file.
Example file:

**JSON Example (.json):**
```json
{
  "CACHE_ADDRESSES": "localhost:6379",
  "CACHE_PASSWORD": "your_password",
  "CACHE_DB": 0,
  "CACHE_NAMESPACE": "myapp"
}
```

**YAML Example (.yaml):**
```yaml
CACHE_ADDRESSES: "localhost:6379"
CACHE_PASSWORD: "your_password"
CACHE_DB: 0
CACHE_NAMESPACE: "myapp"
```

**Environment file (.env):**
```env
CACHE_ADDRESSES=localhost:6379
CACHE_PASSWORD=your_password
CACHE_DB=0
CACHE_NAMESPACE=myapp
```

To use a specific config file (Standalone mode):
```go
os.Setenv("QORE_CONFIG_USED", "./.env.json")
cache := cache.New()
```

### Redis Modes

#### Standalone Redis
```bash
CACHE_ADDRESSES="localhost:6379"
```

#### Redis Cluster
```bash
CACHE_ADDRESSES="node1:6379,node2:6379,node3:6379"
```

#### Redis Sentinel
```bash
CACHE_SENTINEL_ADDRESSES="sentinel1:26379,sentinel2:26379"
CACHE_SENTINEL_MASTER="mymaster"
CACHE_SENTINEL_PASSWORD="sentinel_password"
```

## API Reference

### Cache Operations

#### Storing Data

```go
// Basic store with default TTL (1 minute)
_, err := cache.Set(ctx, "key").Put("value")

// Store with custom TTL
_, err := cache.Set(ctx, "key").
    SetTTL(10 * time.Minute).
    Put("value")

// Store with prefix
_, err := cache.Set(ctx, "user_id").
    SetPrefix("session").
    Put(userData)

// Store forever (no expiration)
_, err := cache.Set(ctx, "key").PutForever("value")
```

#### Retrieving Data

```go
// Basic retrieval
var result string
err := cache.Get(ctx, "key").Pull(&result)

// Retrieval with prefix
var userData User
err := cache.Get(ctx, "user_id", "session").Pull(&userData)

// Remember pattern (get or set default)
var user User
err := cache.Get(ctx, "user:123").Remember(&user, func() (forever bool, val any, err error) {
    // Fetch from database
    user, err := db.GetUser(123)
    return false, user, err
})
```

#### Other Operations

```go
// Check if key exists
exists := cache.Has(ctx, "key", "prefix")

// Delete key
count, err := cache.Delete(ctx, "key", "prefix").Perform()

// Get all keys with prefix
keys, err := cache.GetAllKeys(ctx, "session")
```

### Rate Limiting

```go
// Rate limit: allow once per minute
allowed, err := cache.Set(ctx, "user:123").
    RateLimitOnce(time.Minute)

if !allowed {
    log.Println("Rate limited!")
}
```

### Supported Data Types

The cache supports automatic encoding/decoding for:

- **Primitives**: `string`, `int`, `uint`, `float32`, `float64`, `bool`
- **Byte slices**: `[]byte`
- **Complex types**: Any struct/slice/map (using MessagePack)

```go
// String
cache.Set(ctx, "key").Put("hello")

// Integer
cache.Set(ctx, "counter").Put(42)

// Struct
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
user := User{ID: 1, Name: "John"}
cache.Set(ctx, "user").Put(user)

// Slice
numbers := []int{1, 2, 3, 4, 5}
cache.Set(ctx, "numbers").Put(numbers)
```

## Health Monitoring

```go
// Get health statistics
stats := cache.HealthCheck(ctx)
fmt.Printf("Uptime: %s\n", stats.UptimeHuman)
fmt.Printf("Ping Latency: %dms\n", stats.PINGLatencyMillis)
fmt.Printf("Ping Response: %s\n", stats.PINGResponse)
```

## Dependency Management

The cache implements the `qore.Dependency` interface:

```go
// As a dependency
var deps []qore.Dependency
deps = append(deps, cache.New())

// Dependency info
fmt.Println("Name:", cache.Name())
fmt.Println("Priority:", cache.Priority())
```

## Error Handling

The package defines specific error types:

```go
var (
    ErrClientNil        = errors.New("redis client is null")
    ErrClientNotCluster = errors.New("redis set to cluster mode, but unfortunately the client is not cluster client")
    ErrEmptyKey         = errors.New("cache key cannot be empty")
    ErrEmptyPrefix      = errors.New("prefix cannot be empty")
    ErrOutNonPointer    = errors.New("out type non-pointer")
)
```

## Configuration Reference

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `CACHE_DEPENDENCY_PRIORITY` | Dependency priority for open/close order | `10` |
| `CACHE_NAMESPACE` | Cache key prefix | `"cache-app"` |
| `CACHE_DB` | Redis logical database | `0` |
| `CACHE_USERNAME` | Redis username | `""` |
| `CACHE_PASSWORD` | Redis password | `""` |
| `CACHE_ADDRESSES` | Redis addresses (comma-separated for cluster) | `""` |
| `CACHE_SENTINEL_ADDRESSES` | Sentinel addresses (comma-separated) | `""` |
| `CACHE_SENTINEL_MASTER` | Sentinel master name | `""` |
| `CACHE_SENTINEL_USERNAME` | Sentinel username | `""` |
| `CACHE_SENTINEL_PASSWORD` | Sentinel password | `""` |
| `CACHE_SENTINEL_CLUSTER` | Whether sentinel backend uses cluster | `false` |

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestGet ./...
```

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/qoinlyid/cache"
)

type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}

func main() {
    // Initialize cache
    cache := cache.New()
    if err := cache.Open(); err != nil {
        log.Fatal("Failed to open cache:", err)
    }
    defer cache.Close()
    
    ctx := context.Background()
    
    // Example user
    user := User{
        ID:      123,
        Name:    "John Doe",
        Email:   "john@example.com",
        Created: time.Now(),
    }
    
    // Store user data
    _, err := cache.Set(ctx, fmt.Sprintf("%d", user.ID)).
        SetPrefix("user").
        SetTTL(1 * time.Hour).
        Put(user)
    if err != nil {
        log.Fatal("Failed to store user:", err)
    }
    
    // Retrieve user data
    var retrievedUser User
    err = cache.Get(ctx, "123", "user").Pull(&retrievedUser)
    if err != nil {
        log.Fatal("Failed to retrieve user:", err)
    }
    
    fmt.Printf("Retrieved user: %+v\n", retrievedUser)
    
    // Remember pattern example
    var cachedUser User
    err = cache.Get(ctx, "456", "user").Remember(&cachedUser, func() (forever bool, val any, err error) {
        // Simulate database fetch
        return false, User{
            ID:      456,
            Name:    "Jane Doe",
            Email:   "jane@example.com",
            Created: time.Now(),
        }, nil
    })
    if err != nil {
        log.Fatal("Failed to remember user:", err)
    }
    
    fmt.Printf("Remembered user: %+v\n", cachedUser)
    
    // Rate limiting example
    allowed, err := cache.Set(ctx, "api_call:user:123").
        RateLimitOnce(1 * time.Minute)
    if err != nil {
        log.Fatal("Rate limit error:", err)
    }
    
    if allowed {
        fmt.Println("API call allowed")
    } else {
        fmt.Println("API call rate limited")
    }
    
    // Health check
    stats := cache.HealthCheck(ctx)
    fmt.Printf("Cache uptime: %s\n", stats.UptimeHuman)
    fmt.Printf("Ping latency: %dms\n", stats.PINGLatencyMillis)
}
```