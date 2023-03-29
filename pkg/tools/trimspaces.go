package tools

import (
	"fmt"
	"reflect"
	"strings"
)

func TrimSpaces(iface interface{}) error {
	if !(reflect.ValueOf(iface).Type().Kind() == reflect.Ptr) {
		return fmt.Errorf("TrimSpaces: iface is not a pointer")
	}

	ifv := reflect.ValueOf(iface)
	ind := reflect.Indirect(ifv)
	switch ind.Kind() {

	case reflect.String:
		ind.SetString(strings.TrimSpace(ind.String()))

	case reflect.Struct:
		for i := 0; i < ind.NumField(); i++ {
			v := ind.Field(i)
			switch v.Kind() {

			case reflect.Struct, reflect.Slice:
				TrimSpaces(v.Addr().Interface())

			case reflect.String:
				v.SetString(strings.TrimSpace(v.String()))

			}
		}

	case reflect.Slice:
		for i := 0; i < ind.Len(); i++ {
			v := ind.Index(i)
			switch v.Kind() {
			case reflect.String:
				v.SetString(strings.TrimSpace(v.String()))
			}
		}

	case reflect.Map:
		for _, k := range ind.MapKeys() {
			v := ind.MapIndex(k)
			switch v.Kind() {
			case reflect.String:
				ind.SetMapIndex(k, reflect.ValueOf(strings.TrimSpace(v.String())))
			}
		}
	}

	return nil
}
