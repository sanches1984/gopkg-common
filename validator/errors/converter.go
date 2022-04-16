package errors

import (
	"context"
	"github.com/sanches1984/gopkg-common/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/sanches1984/gopkg-errors"

	"github.com/go-playground/validator/v10"
)

const ValidationMsg = "Invalid request payload"

type Code string

// Global validation codes
const (
	CommonCode       Code = "INVALID"
	RequiredCode     Code = "REQUIRED"
	EqLenCode        Code = "LEN"
	MinLenCode       Code = "MIN_LEN"
	MaxLenCode       Code = "MAX_LEN"
	MinCode          Code = "MIN"
	MaxCode          Code = "MAX"
	TooFewItemsCode  Code = "TOO_FEW_ELEM"
	TooManyItemsCode Code = "TOO_MANY_ELEM"
)

// Converter converts error from validation errors
func Converter() errors.ErrorConverter {
	return func(ctx context.Context, err error) (*errors.Error, bool) {
		for {
			if errValidation, ok := err.(validator.ValidationErrors); ok {
				return convertValidationError(ctx, errValidation), true
			}
			errC, ok := err.(errors.Causer)
			if !ok {
				return nil, false
			}
			err = errC.Cause()
		}
	}
}

func BuildCode(code Code, param interface{}) string {
	ret := removeTrashOrCondition(string(code))
	strParam := ""
	switch v := param.(type) {
	case string:
		strParam = v
	case int:
		strParam = strconv.Itoa(v)
	case int8:
		strParam = strconv.FormatInt(int64(v), 10)
	case int32:
		strParam = strconv.FormatInt(int64(v), 10)
	case int64:
		strParam = strconv.FormatInt(v, 10)
	}
	if param != "" {
		ret += ":" + strParam
	}
	return ret
}

var trashOrConditionRE = regexp.MustCompile(`\|?NOT_UP\|?`)

func removeTrashOrCondition(str string) string {
	return trashOrConditionRE.ReplaceAllString(str, "")
}

func convertValidationError(ctx context.Context, err validator.ValidationErrors) *errors.Error {
	result := errors.BadRequest.ErrWrap(ctx, ValidationMsg, err)
	// https://godoc.org/github.com/go-playground/validator
	for _, fieldErr := range err {
		value := Code(strings.ToUpper(fieldErr.Tag()))
		switch fieldErr.Tag() {
		case "min":
			if fieldErr.Kind() == reflect.Slice {
				value = TooFewItemsCode
			}
			if fieldErr.Kind() == reflect.String {
				value = MinLenCode
			}
		case "max":
			if fieldErr.Kind() == reflect.Slice {
				value = TooManyItemsCode
			}
			if fieldErr.Kind() == reflect.String {
				value = MaxLenCode
			}
		}

		field := types.CamelToSnakeCase(fieldErr.Field())
		result = result.WithPayloadKV(field, BuildCode(value, fieldErr.Param()))
	}

	return result
}
