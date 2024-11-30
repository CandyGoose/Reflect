package main

import (
	"errors"
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	val := reflect.ValueOf(out)

	if val.Kind() != reflect.Pointer {
		return errors.New("out must be a pointer")
	}
	if val.IsNil() {
		return errors.New("out is nil")
	}

	val = val.Elem()

	switch val.Kind() {
	case reflect.Struct:
		d, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map[string]interface{}, got %T", data)
		}

		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			fieldName := val.Type().Field(i).Name
			dataValue, ok := d[fieldName]
			if !ok {
				continue
			}

			err := i2s(dataValue, valueField.Addr().Interface())
			if err != nil {
				return fmt.Errorf("failed to set field %s: %v", fieldName, err)
			}
		}

	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("expected []interface{}, got %T", data)
		}

		v := reflect.MakeSlice(val.Type(), len(d), len(d))
		for idx, elem := range d {
			err := i2s(elem, v.Index(idx).Addr().Interface())
			if err != nil {
				return fmt.Errorf("failed to set slice element at index %d: %v", idx, err)
			}
		}
		val.Set(v)

	case reflect.Bool:
		d, ok := data.(bool)
		if !ok {
			return fmt.Errorf("expected bool, got %T", data)
		}
		val.SetBool(d)

	case reflect.Int:
		d, ok := data.(float64)
		if !ok {
			return fmt.Errorf("expected float64, got %T for int", data)
		}
		val.SetInt(int64(d))

	case reflect.String:
		d, ok := data.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", data)
		}
		val.SetString(d)

	default:
		return fmt.Errorf("unsupported type %s", val.Kind())
	}

	return nil
}
