package runtime

import (
	"fmt"
	"reflect"
	"strconv"
)

func makeInt(val interface{}) (int64, bool) {
	rval := reflect.ValueOf(val)
	switch rval.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		return rval.Int(), true
	case reflect.String:
		if ival, err := strconv.ParseInt(rval.String(), 10, 64); err == nil {
			return ival, true
		}
	}

	return 0, false
}

func makeFloat(val interface{}) (float64, bool) {
	rval := reflect.ValueOf(val)
	switch rval.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		return float64(rval.Int()), true
	case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint16, reflect.Uint8:
		return float64(rval.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rval.Float(), true
	case reflect.String:
		if fval, err := strconv.ParseFloat(rval.String(), 64); err == nil {
			return fval, true
		}
	}

	return 0, false
}

func makeString(val interface{}) (string, bool) {
	rval := reflect.ValueOf(val)
	switch rval.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		return fmt.Sprintf("%d", rval.Int()), true
	case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint16, reflect.Uint8:
		return fmt.Sprintf("%d", rval.Uint()), true
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", rval.Float()), true
	case reflect.String:
		return rval.String(), true
	}

	return "", false
}

func makeBool(val interface{}) bool {
	rval := reflect.ValueOf(val)

	if rval.CanAddr() && rval.IsNil() {
		return false
	}

	if rval.IsValid() {
		switch rval.Kind() {
		case reflect.Bool:
			return rval.Bool()
		case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
			return rval.Int() != 0
		case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint16, reflect.Uint8:
			return rval.Uint() != 0
		case reflect.Float32, reflect.Float64:
			return rval.Float() != 0
		case reflect.Complex64, reflect.Complex128:
			return rval.Complex() != 0
		case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
			return !rval.IsNil()
		case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
			return rval.Len() > 0
		case reflect.Struct:
			return true
		}
	}

	return false
}
