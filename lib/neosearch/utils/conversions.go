package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func BoolFromInterface(value interface{}, kind reflect.Kind) (bool, error) {
	switch kind {
	case reflect.String:
		sval := value.(string)

		if sval == "true" {
			return true, nil
		} else if sval == "false" {
			return false, nil
		} else {
			return false, fmt.Errorf("Invalid boolean: %s", value)
		}
	case reflect.Int:
		ival := value.(int64)

		if ival == 0 {
			return false, nil
		}

		return true, nil
	case reflect.Uint:
		uval := value.(uint64)

		if uval == 0 {
			return false, nil
		}

		return true, nil
	case reflect.Float32, reflect.Float64:
		fval := value.(float64)

		if fval == 0.0 {
			return false, nil
		}

		return true, nil
	default:
		return false, fmt.Errorf("Impossible to convert '%s' to boolean", value)
	}
}

func Uint64FromInterface(value interface{}, kind reflect.Kind) (uint64, error) {
	var (
		uret uint64
		err  error
	)

	switch kind {
	case reflect.String:
		intval, err := strconv.Atoi(value.(string))

		if err != nil {
			return uret, fmt.Errorf("Failed to convert '%s' to uint64", value)
		}

		uret = uint64(intval)
	case reflect.Float32, reflect.Float64:
		uret = uint64(value.(float64))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		uret = uint64(value.(int64))
	default:
		err = fmt.Errorf("Impossible to convert '%s' to uint64", value)
	}

	return uret, err
}

func Int64FromInterface(value interface{}, kind reflect.Kind) (int64, error) {
	var (
		ret int64
		err error
	)

	switch kind {
	case reflect.String:
		intval, err := strconv.Atoi(value.(string))

		if err != nil {
			return ret, fmt.Errorf("Failed to convert '%s' to int64", value)
		}

		ret = int64(intval)
	case reflect.Float32, reflect.Float64:
		ret = int64(value.(float64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret = int64(value.(uint64))
	default:
		err = fmt.Errorf("Impossible to convert '%s' to int64", value)
	}

	return ret, err
}
