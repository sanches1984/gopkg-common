package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateMinLetter(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("min_letter", ValidateMinLetter)
	if err != nil {
		t.Fatal(err)
	}

	type field struct {
		Val string `validate:"min_letter=2"`
	}

	cs := map[string]bool{
		"":    false,
		"1":   false,
		"я":   false,
		"d":   false,
		"12":  true,
		"яв":  true,
		"gh":  true,
		"123": true,
		"ява": true,
		"jfg": true,
	}

	for value, ok := range cs {
		err := v.Struct(field{Val: value})
		if ok {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func Test_validateEmailList(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("email_list", ValidateEmailList)
	if err != nil {
		t.Fatal(err)
	}

	type field struct {
		Val string `validate:"email_list"`
	}

	cs := map[string]bool{
		"a@a.a":               true,
		"a@a.a, b@b.b":        true,
		"a@a.a,b@b.b , c@c.c": true,
		"a.d.v.b@d.d.f.g":     true,
		"notemail":            false,
		"a@a.a b@b.b":         false,
		"a@aa, b@b.b":         false,
	}

	for value, ok := range cs {
		err := v.Struct(field{Val: value})
		if ok {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
