package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/qoinlyid/qore"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

// open is helper function to open redis connection based on appropriate client.
func (i *Instance) open() error {
	var addrs []string

	// Read field "Addresses" first from config, means the firt & priority is using redis standalone
	// or cluster mode.
	if !qore.ValidationIsEmpty(i.cfg.Addresses) {
		i.clustering = strings.Contains(i.cfg.Addresses, ",")
		for addr := range strings.SplitSeq(i.cfg.Addresses, ",") {
			if len(strings.TrimSpace(addr)) > 0 {
				addrs = append(addrs, addr)
			}
		}
		if len(addrs) == 0 {
			return errors.New("[cache] failed to open connection: redis.Addresses is empty")
		}

		// New client.
		if i.clustering {
			// Redis client cluster.
			i.client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:      addrs,
				Username:   i.cfg.Username,
				Password:   i.cfg.Password,
				ClientName: i.cfg.Namespace,
			})
		} else {
			i.client = redis.NewClient(&redis.Options{
				Addr:       addrs[0],
				ClientName: i.cfg.Namespace,
				Username:   i.cfg.Username,
				Password:   i.cfg.Password,
				DB:         i.cfg.DB,
			})
		}
		return nil
	}

	// If field "Addresses" empty read through field "SentinelAddresses", if any open redis using sentinel.
	if !qore.ValidationIsEmpty(i.cfg.SentinelAddresses) {
		for addr := range strings.SplitSeq(i.cfg.SentinelAddresses, ",") {
			if len(strings.TrimSpace(addr)) > 0 {
				addrs = append(addrs, addr)
			}
		}
		if len(addrs) == 0 {
			return errors.New("[cache] failed to open connection: redis.SentinelAddresses is empty")
		}

		i.clustering = i.cfg.SentinelCluster
		sentOpts := &redis.FailoverOptions{
			// Sentinel.
			MasterName:       i.cfg.SentinelMaster,
			SentinelAddrs:    addrs,
			SentinelUsername: i.cfg.SentinelUsername,
			SentinelPassword: i.cfg.SentinelPassword,
			ClientName:       i.cfg.Namespace,

			// Redis.
			Username: i.cfg.Username,
			Password: i.cfg.Password,
			DB:       i.cfg.DB,
		}

		// Sentinel cluster?.
		if i.clustering {
			sentOpts.RouteByLatency = true
			i.client = redis.NewFailoverClusterClient(sentOpts)
		} else {
			i.client = redis.NewFailoverClient(sentOpts)
		}
		return nil
	}

	// Fallback error.
	return errors.New("at least one of redis addresses and sentinel addresses must be defines")
}

type base struct {
	ctx    context.Context
	cancel context.CancelFunc
	key    string
}

func (i *Instance) validateClient() error {
	if i.client == nil {
		return ErrClientNil
	}
	return nil
}

func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// encoder return encoded bytes from the `val` based on the type.
func encoder(val any) ([]byte, error) {
	switch v := val.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case int8, int16, int32, int64, int,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		return fmt.Appendf(nil, "%v", v), nil
	default:
		data, err := msgpack.Marshal(v)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

// decoder parses data to the out.
func decoder(data []byte, out any) error {
	// Validate.
	if out == nil {
		return errors.New("output parameter is nil")
	}

	// Out type checker.
	switch v := out.(type) {
	case *string:
		*v = string(data)
		return nil
	case *[]byte:
		*v = append((*v)[:0], data...)
		return nil

	// Primitive int.
	case *int:
		i, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}
		*v = i
		return nil
	case *int8:
		i64, err := strconv.ParseInt(string(data), 10, 8)
		if err != nil {
			return err
		}
		*v = int8(i64)
		return nil
	case *int16:
		i64, err := strconv.ParseInt(string(data), 10, 16)
		if err != nil {
			return err
		}
		*v = int16(i64)
		return nil
	case *int32:
		i64, err := strconv.ParseInt(string(data), 10, 32)
		if err != nil {
			return err
		}
		*v = int32(i64)
		return nil
	case *int64:
		i64, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*v = i64
		return nil

	// Primitive unsigned int.
	case *uint:
		u64, err := strconv.ParseUint(string(data), 10, 0)
		if err != nil {
			return err
		}
		*v = uint(u64)
		return nil
	case *uint8:
		u64, err := strconv.ParseUint(string(data), 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(u64)
		return nil
	case *uint16:
		u64, err := strconv.ParseUint(string(data), 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(u64)
		return nil
	case *uint32:
		u64, err := strconv.ParseUint(string(data), 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(u64)
		return nil
	case *uint64:
		u64, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return err
		}
		*v = u64
		return nil

	// Primitive floating number.
	case *float32:
		f64, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return err
		}
		*v = float32(f64)
		return nil
	case *float64:
		f64, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*v = f64
		return nil

	// Primitive boolean.
	case *bool:
		b, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		*v = b
		return nil

	// Fallback using msgpack.
	default:
		ov := reflect.ValueOf(out)
		if ov.Kind() != reflect.Ptr {
			return fmt.Errorf("%w %T", ErrOutNonPointer, out)
		}
		return msgpack.Unmarshal(data, out)
	}
}

// set helper to store the value into redis storage.
func (i *Instance) set(set *setter, val any, borrow ...bool) (redis.UniversalClient, error) {
	// Validate.
	if err := i.validateClient(); err != nil {
		set.cleanup()
		return nil, err
	}
	if qore.ValidationIsEmpty(set.key) {
		set.cleanup()
		return nil, ErrEmptyKey
	}
	set.key = i.cfg.Namespace + DefaultKeySeparator + set.key

	// Want borrow client?
	if len(borrow) > 0 {
		if borrow[0] {
			return i.client, nil
		}
	}
	defer set.cleanup()

	// Exec.
	encoded, err := encoder(val)
	if err != nil {
		return nil, fmt.Errorf("failed to encode value %T: %w", val, err)
	}
	return nil, i.client.Set(set.ctx, set.key, encoded, set.ttl).Err()
}

// getRemember helper to get default value and set it to the redis storage.
func (i *Instance) getRemember(
	get *getter,
	out any,
	rem RememberFn,
) (next bool, err error) {
	// Check is exist.
	exist := i.Has(get.ctx, get.key)
	if exist {
		return true, nil
	}

	// Check out must be pointer.
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr || outVal.IsNil() {
		return false, errors.New("out must be a non-nil pointer")
	}

	// Call closure function to get the default value.
	forever, val, err := rem()
	if err != nil {
		return false, err
	}

	// Make sure val can be assigned to the out.
	valVal := reflect.ValueOf(val)
	if !valVal.Type().AssignableTo(outVal.Elem().Type()) {
		return false, fmt.Errorf(
			"cannot assign value of type %s to out of type %s",
			valVal.Type(), outVal.Elem().Type(),
		)
	}

	// Store value into cache.
	set := i.Set(get.ctx, get.key)
	if forever {
		if _, err := set.PutForever(val); err != nil {
			return false, err
		}
	} else {
		if _, err := set.Put(val); err != nil {
			return false, err
		}
	}

	// Ok.
	outVal.Elem().Set(valVal)
	return false, nil
}

// get helper to retrieve the value from redis storage.
func (i *Instance) get(get *getter, out any, rem ...RememberFn) error {
	defer func() {
		if get.cancel != nil {
			get.cancel()
		}
		// Zero out all fields to help GC or prepare for reuse
		*get = getter{}
	}()

	// Validate.
	if qore.ValidationIsEmpty(get.key) {
		return ErrEmptyKey
	}
	if err := i.validateClient(); err != nil {
		return err
	}

	// Remember: get default value and store it to the redis storage.
	if len(rem) > 0 {
		if rem[0] != nil {
			next, err := i.getRemember(get, out, rem[0])
			if err != nil {
				return err
			} else if !next {
				return nil
			}
		}
	}

	// Exec.
	get.key = i.cfg.Namespace + DefaultKeySeparator + get.key
	b, err := i.client.Get(get.ctx, get.key).Bytes()
	if err != nil {
		return err
	}
	err = decoder(b, out)
	if err != nil {
		return fmt.Errorf("failed to decode value to %T: %w", out, err)
	}
	return nil
}

func (i *Instance) del(del *deleter) (int64, error) {
	defer func() {
		if del.cancel != nil {
			del.cancel()
		}
		// Zero out all fields to help GC or prepare for reuse
		*del = deleter{}
	}()

	// Validate.
	if qore.ValidationIsEmpty(del.key) {
		return 0, ErrEmptyKey
	}
	if err := i.validateClient(); err != nil {
		return 0, err
	}

	// Exec.
	del.key = i.cfg.Namespace + DefaultKeySeparator + del.key
	return i.client.Del(context.Background(), del.key).Result()
}
