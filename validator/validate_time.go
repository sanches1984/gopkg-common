package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	timeHHMMFormat   = "15:04"
	timeHHMMSSFormat = "15:04:05"
)

// ValidateTimeHHMM validate field time. format hh:mm
func ValidateTimeHHMM(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) == 0 {
		return true
	}
	if len(val) != 5 {
		return false
	}
	_, err := time.Parse(timeHHMMFormat, val)
	return err == nil
}

// ValidateTimeHHMMSS validate field time. format hh:mm:ss
func ValidateTimeHHMMSS(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) == 0 {
		return true
	}
	if len(val) != 8 {
		return false
	}
	_, err := time.Parse(timeHHMMSSFormat, val)
	return err == nil
}
