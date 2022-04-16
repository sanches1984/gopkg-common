package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_validateFieldEmpty(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("field_empty", ValidateFieldEmpty)
	if err != nil {
		t.Fatal(err)
	}
	type check struct {
		First  string
		Second string `validate:"field_empty=First|oneof=val1 val2"`
	}

	t.Run("FirstEmpty_SecondEmpty", func(t *testing.T) {
		s := check{
			First:  "",
			Second: "",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})

	t.Run("FirstEmpty_SecondAny", func(t *testing.T) {
		s := check{
			First:  "",
			Second: "xyz",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})

	t.Run("FirstEmpty_SecondValid", func(t *testing.T) {
		s := check{
			First:  "",
			Second: "val2",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})

	t.Run("FirstFilled_SecondEmpty", func(t *testing.T) {
		s := check{
			First:  "val1",
			Second: "",
		}

		err := v.Struct(s)
		assert.NotNil(t, err)
	})

	t.Run("FirstFilled_SecondFilled", func(t *testing.T) {
		s := check{
			First:  "val1",
			Second: "xyz",
		}

		err := v.Struct(s)
		assert.NotNil(t, err)
	})

	t.Run("FirstFilled_SecondValid", func(t *testing.T) {
		s := check{
			First:  "val1",
			Second: "val2",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})
}

func Test_validateFieldNotUpdated(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("not_up", ValidateFieldNotUpdated)
	if err != nil {
		t.Fatal(err)
	}
	type check struct {
		FirstValue   string `validate:"not_up=First|min=1"`
		SecondValue  string
		UpdateFields []string
	}

	cs := []struct {
		message string
		value   check
		valid   bool
	}{
		{
			message: "FirstValid, UpdateFieldsEmpty, Success",
			value: check{
				FirstValue:   "valid",
				UpdateFields: []string{},
			},
			valid: true,
		},
		{
			message: "FirstNotValid, UpdateFieldsEmpty, Failed",
			value: check{
				FirstValue:   "",
				UpdateFields: []string{},
			},
			valid: false,
		},
		{
			message: "FirstValid, UpdateFieldsThat, Success",
			value: check{
				FirstValue:   "valid",
				UpdateFields: []string{"first_value"},
			},
			valid: true,
		},
		{
			message: "FirstValid, UpdateFieldsThatAndOther, Success",
			value: check{
				FirstValue:   "valid",
				UpdateFields: []string{"zero", "first_value", "second"},
			},
			valid: true,
		},
		{
			message: "FirstNotValid, UpdateFieldsThat, Failed",
			value: check{
				FirstValue:   "",
				UpdateFields: []string{"first_value"},
			},
			valid: false,
		},
		{
			message: "FirstNotValid, UpdateFieldsOther, Success",
			value: check{
				FirstValue:   "",
				UpdateFields: []string{"second"},
			},
			valid: true,
		},
	}

	for _, c := range cs {
		err := v.Struct(c.value)
		if c.valid {
			assert.Nil(t, err, c.message)
		} else {
			assert.NotNil(t, err, c.message)
		}
	}
}
