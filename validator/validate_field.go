package validator

import (
	"github.com/sanches1984/gopkg-common/types"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// ValidateFieldRequired ...
func ValidateFieldRequired(fl validator.FieldLevel) bool {
	field, _, _, ok := fl.GetStructFieldOK2()
	if !ok {
		return true
	}
	return !isEmptyValue(field)
}

// ValidateFieldEmpty ...
func ValidateFieldEmpty(fl validator.FieldLevel) bool {
	field, _, _, ok := fl.GetStructFieldOK2()
	if !ok {
		return true
	}
	return isEmptyValue(field)
}

// ValidateFieldNotUpdated checks if the field is participating in the update,
// this happens if UpdateFields is empty or the current field is listed in UpdateFields
func ValidateFieldNotUpdated(fl validator.FieldLevel) bool {
	top := fl.Top()
	if top.Kind() == reflect.Ptr {
		top = top.Elem()
	}
	updateFields := top.FieldByName("UpdateFields")
	if updateFields.Kind() != reflect.Slice {
		return true
	}
	if updateFields.Len() == 0 {
		return false
	}
	if updateFields.Index(0).Kind() != reflect.String {
		return true
	}
	fieldName := types.CamelToSnakeCase(fl.FieldName())
	for i := 0; i < updateFields.Len(); i++ {
		if updateFields.Index(i).String() == fieldName {
			return false
		}
	}
	return true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
