package compiler

import (
	"reflect"
)

func reflectTypeAndCompare(value interface{}, expected map[string]interface{}) bool {
	v := reflect.ValueOf(value)
	e := v.Elem()

	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varValue := e.Field(i).Interface()

		if val, ok := expected[varName]; ok {
			if !reflect.DeepEqual(val, varValue) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}
