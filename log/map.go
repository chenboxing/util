package log

import (
	"encoding/json"
)

func Map(val map[string]string) Fields {
	tmp := make(map[string]interface{})
	for k, v := range val {
		tmp[k] = v
	}
	return Fields(tmp)
}

// func Struct(st interface{}) Fields {
// 	val := reflect.ValueOf(st)
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 	}
// 	typ := val.Type()
// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 	}

// 	fields := make(map[string]interface{})
// 	for i := 0; i < val.NumField(); i++ {
// 		fields[typ.Field(i).Name] = val.Field(i).Interface()
// 	}

// 	return fields
// }

func Struct(st interface{}) Fields {
	body, _ := json.Marshal(st)

	return Fields{
		"struct_json": string(body),
	}
}
