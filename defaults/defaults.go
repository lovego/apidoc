// Set default values for struct
package defaults

import (
	"reflect"
)

func Set(ptr interface{}) interface{} {
	v := reflect.ValueOf(ptr)
	if reflect.TypeOf(ptr).Kind() == reflect.Ptr {
		v = reflect.ValueOf(ptr).Elem()
	}
	setField(v)
	return v.Interface()
}

func setField(field reflect.Value) {
	if !field.CanSet() {
		// fmt.Println(`can not set :`, field.Type().Name())
		return
	}
	switch field.Kind() {
	case reflect.Array:
		for j := 0; j < field.Len(); j++ {
			setField(field.Index(j))
		}
	case reflect.Map:
		mType := field.Type()
		field.Set(reflect.MakeMap(mType))

		key := reflect.New(mType.Key()).Elem()
		val := reflect.New(mType.Elem()).Elem()
		setField(val)
		field.SetMapIndex(key, val)
	case reflect.Slice:
		field.Set(reflect.MakeSlice(field.Type(), 1, 1))
		setField(field.Index(0))
	case reflect.Ptr:
		field.Set(reflect.New(field.Type().Elem()))
		setField(field.Elem())
	case reflect.Interface:
		setField(field.Elem())
	case reflect.Struct:
		for i := 0; i < field.NumField(); i++ {
			setField(field.Field(i))
		}
	}
}
