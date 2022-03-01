package docs

import (
	"fmt"
	"reflect"
)

type ConfigField struct {
	Name         string
	DefaultValue string
	Type         string
	Description  string
}

func Display(s interface{}) []ConfigField {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	var fields []ConfigField
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		switch value.Kind() {
		default:
			desc := field.Tag.Get("desc")
			env, ok := field.Tag.Lookup("env")
			if !ok {
				continue
			}
			v := fmt.Sprintf("%v", value.Interface())
			fields = append(fields, ConfigField{Name: env, DefaultValue: v, Description: desc, Type: value.Type().Name()})
		case reflect.Struct:
			fields = append(fields, Display(value.Interface())...)
		}
	}
	return fields
}
