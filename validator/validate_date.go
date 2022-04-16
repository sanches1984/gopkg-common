package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const rfcDateFormat = "2006-01-02"

// ValidateDateYMD ...
func ValidateDateYMD(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) == 0 {
		return true
	}
	if len(val) != 10 {
		return false
	}
	_, err := time.Parse(rfcDateFormat, val)
	return err == nil
}

//ValidateDateRfc3339 validate field against RFC3339 format
func ValidateDateRfc3339(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) == 0 {
		return true
	}

	if len(val) < 20 {
		return false
	}

	_, err := time.Parse(time.RFC3339, val)
	return err == nil
}
