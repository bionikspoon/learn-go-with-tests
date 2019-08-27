package reflection

import (
	"reflect"
)

func Walk(x interface{}, fn func(input string)) {
	value := getValue(x)

	switch value.Kind() {
	case reflect.String:
		fn(value.String())
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			Walk(value.Field(i).Interface(), fn)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			Walk(value.Index(i).Interface(), fn)
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			Walk(value.MapIndex(key).Interface(), fn)
		}
	}
}

func getValue(x interface{}) reflect.Value {
	value := reflect.ValueOf(x)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value
}
