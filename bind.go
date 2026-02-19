package configo

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Bind populates a struct from config values using `config` and `default` struct tags.
func (c *Config) Bind(target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("bind: target must be a pointer to a struct")
	}
	return c.bindStruct(v.Elem())
}

func (c *Config) bindStruct(v reflect.Value) error {
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		fv := v.Field(i)

		if !fv.CanSet() {
			continue
		}

		// Handle embedded/nested structs
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Duration(0)) {
			if err := c.bindStruct(fv); err != nil {
				return err
			}
			continue
		}

		key := field.Tag.Get("config")
		if key == "" {
			continue
		}

		c.mu.RLock()
		val, ok := c.data[key]
		c.mu.RUnlock()

		if !ok {
			defStr := field.Tag.Get("default")
			if defStr != "" {
				if err := setFieldFromString(fv, defStr); err != nil {
					return fmt.Errorf("bind: field %s default: %w", field.Name, err)
				}
			}
			continue
		}

		if err := setField(fv, val); err != nil {
			return fmt.Errorf("bind: field %s: %w", field.Name, err)
		}
	}
	return nil
}

func setField(fv reflect.Value, val any) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(fmt.Sprintf("%v", val))
	case reflect.Int, reflect.Int64:
		if fv.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := toDuration(val)
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(d))
			return nil
		}
		n, err := toInt64(val)
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Float64:
		n, err := toFloat64(val)
		if err != nil {
			return err
		}
		fv.SetFloat(n)
	case reflect.Bool:
		b, err := toBool(val)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Slice:
		return setSliceField(fv, val)
	default:
		return fmt.Errorf("unsupported field type %s", fv.Kind())
	}
	return nil
}

func setSliceField(fv reflect.Value, val any) error {
	switch fv.Type().Elem().Kind() {
	case reflect.String:
		s, err := toStringSlice(val)
		if err != nil {
			return err
		}
		fv.Set(reflect.ValueOf(s))
	case reflect.Int:
		s, err := toIntSlice(val)
		if err != nil {
			return err
		}
		fv.Set(reflect.ValueOf(s))
	default:
		return fmt.Errorf("unsupported slice type %s", fv.Type())
	}
	return nil
}

func setFieldFromString(fv reflect.Value, s string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(s)
	case reflect.Int, reflect.Int64:
		if fv.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(d))
			return nil
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	default:
		return fmt.Errorf("unsupported default type %s", fv.Kind())
	}
	return nil
}
