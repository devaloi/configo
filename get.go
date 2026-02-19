package configo

import (
	"fmt"
	"strconv"
	"time"
)

// Get retrieves a typed value from the config.
func Get[T any](c *Config, key string) (T, error) {
	c.mu.RLock()
	val, ok := c.data[key]
	c.mu.RUnlock()

	var zero T
	if !ok {
		return zero, &KeyNotFoundError{Key: key}
	}

	result, err := coerce[T](val)
	if err != nil {
		return zero, &TypeMismatchError{Key: key, Expected: fmt.Sprintf("%T", zero), Actual: val}
	}
	return result, nil
}

// GetOr retrieves a typed value, returning fallback if the key is missing.
func GetOr[T any](c *Config, key string, fallback T) T {
	val, err := Get[T](c, key)
	if err != nil {
		return fallback
	}
	return val
}

// MustGet retrieves a typed value and panics if the key is missing or conversion fails.
func MustGet[T any](c *Config, key string) T {
	val, err := Get[T](c, key)
	if err != nil {
		panic(fmt.Sprintf("configo: %v", err))
	}
	return val
}

func coerce[T any](val any) (T, error) {
	var zero T
	target := any(&zero)

	switch ptr := target.(type) {
	case *string:
		*ptr = fmt.Sprintf("%v", val)
		return zero, nil
	case *int:
		v, err := toInt64(val)
		if err != nil {
			return zero, err
		}
		*ptr = int(v)
		return zero, nil
	case *int64:
		v, err := toInt64(val)
		if err != nil {
			return zero, err
		}
		*ptr = v
		return zero, nil
	case *float64:
		v, err := toFloat64(val)
		if err != nil {
			return zero, err
		}
		*ptr = v
		return zero, nil
	case *bool:
		v, err := toBool(val)
		if err != nil {
			return zero, err
		}
		*ptr = v
		return zero, nil
	case *time.Duration:
		v, err := toDuration(val)
		if err != nil {
			return zero, err
		}
		*ptr = v
		return zero, nil
	case *[]string:
		v, err := toStringSlice(val)
		if err != nil {
			return zero, err
		}
		*ptr = v
		return zero, nil
	case *[]int:
		v, err := toIntSlice(val)
		if err != nil {
			return zero, err
		}
		*ptr = v
		return zero, nil
	default:
		if typed, ok := val.(T); ok {
			return typed, nil
		}
		return zero, fmt.Errorf("unsupported type %T", zero)
	}
}

func toInt64(val any) (int64, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", val)
	}
}

func toFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", val)
	}
}

func toBool(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", val)
	}
}

func toDuration(val any) (time.Duration, error) {
	switch v := val.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(v)
	case int:
		return time.Duration(v) * time.Millisecond, nil
	case int64:
		return time.Duration(v) * time.Millisecond, nil
	case float64:
		return time.Duration(v) * time.Millisecond, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to duration", val)
	}
}

func toStringSlice(val any) ([]string, error) {
	switch v := val.(type) {
	case []string:
		return v, nil
	case []any:
		out := make([]string, len(v))
		for i, item := range v {
			out[i] = fmt.Sprintf("%v", item)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []string", val)
	}
}

func toIntSlice(val any) ([]int, error) {
	switch v := val.(type) {
	case []int:
		return v, nil
	case []any:
		out := make([]int, len(v))
		for i, item := range v {
			n, err := toInt64(item)
			if err != nil {
				return nil, err
			}
			out[i] = int(n)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []int", val)
	}
}
