package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	initOnce sync.Once
	validate *validator.Validate
)

// New create new validator instance with custom tag
func New() *validator.Validate {
	initOnce.Do(func() {
		validate = validator.New()
		customValidations := map[string]validator.Func{
			"date_ymd":       ValidateDateYMD,
			"date_rfc3339":   ValidateDateRfc3339,
			"str_gte":        ValidateStrGTE,
			"str_lte":        ValidateStrLTE,
			"str_int":        ValidateStrInt,
			"field_empty":    ValidateFieldEmpty,
			"field_required": ValidateFieldRequired,
			"not_up":         ValidateFieldNotUpdated,
			"time_hhmm":      ValidateTimeHHMM,
			"time_hhmmss":    ValidateTimeHHMMSS,
			"min_letter":     ValidateMinLetter,
			"email_list":     ValidateEmailList,
		}
		for tag, fn := range customValidations {
			_ = validate.RegisterValidation(tag, fn)
		}
	})

	return validate
}
