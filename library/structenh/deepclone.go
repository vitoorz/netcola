package structenh

import (
	"reflect"
)

func DeepClone(origin interface{}) interface{} {
	if origin == nil {
		return nil
	}

	original := reflect.ValueOf(origin)

	cpy := reflect.New(original.Type()).Elem()

	copyRecursive(original, cpy)

	//fmt.Println(InterfacePresentation(cpy.Interface()))
	return cpy.Interface()
}

func copyRecursive(original, cpy reflect.Value) {
	// handle according to original's Kind
	if !cpy.CanSet() {
		return
	}

	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()
		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, cpy.Elem())
	case reflect.Interface:
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}
		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)
	case reflect.Struct:
		// Go through each field of the struct and copy it.
		oriType := original.Type()
		for i := 0; i < original.NumField(); i++ {
			if cpy.Field(i).CanSet() &&
				(oriType.Field(i).Tag.Get("bson") != "-" || oriType.Field(i).Tag.Get("cpy") != "") {
				copyRecursive(original.Field(i), cpy.Field(i))
			}
		}
	case reflect.Slice:
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i))
		}
	case reflect.Map:
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			cpy.SetMapIndex(key, copyValue)
		}
	// Set the actual values from here on.
	default:
		cpy.Set(original)
	}
}
