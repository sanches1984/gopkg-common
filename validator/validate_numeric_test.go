package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-playground/validator/v10"
)

func Test_validateSGTE(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("str_gte", ValidateStrGTE)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Greater", func(t *testing.T) {
		s := struct {
			Val string `validate:"str_gte=1.2345"`
		}{
			Val: "1.3",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})

	t.Run("Equal", func(t *testing.T) {
		s := struct {
			Val string `validate:"str_gte=1.2345"`
		}{
			Val: "1.2345",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})

	t.Run("Less", func(t *testing.T) {
		s := struct {
			Val string `validate:"str_gte=1.2345"`
		}{
			Val: "1.2",
		}

		err := v.Struct(s)
		assert.NotNil(t, err)
	})
}

func Test_validateStrLTE(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("str_lte", ValidateStrLTE)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Greater", func(t *testing.T) {
		s := struct {
			Val string `validate:"str_lte=1.2345"`
		}{
			Val: "1.3",
		}

		err := v.Struct(s)
		assert.NotNil(t, err)
	})

	t.Run("Equal", func(t *testing.T) {
		s := struct {
			Val string `validate:"str_lte=1.2345"`
		}{
			Val: "1.2345",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})

	t.Run("Less", func(t *testing.T) {
		s := struct {
			Val string `validate:"str_lte=1.2345"`
		}{
			Val: "1.2",
		}

		err := v.Struct(s)
		assert.Nil(t, err)
	})
}

func Test_validateStrInt(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("str_int", ValidateStrInt)
	if err != nil {
		t.Fatal(err)
	}

	type request struct {
		Val string `validate:"str_int"`
	}

	cs := []struct {
		val string
		ok  bool
	}{
		{val: "", ok: true},
		{val: "12", ok: true},
		{val: "1278314567823648756236745621345621575647855926873", ok: true},
		{val: "12.", ok: false},
		{val: "1000 00", ok: false},
	}

	for _, c := range cs {
		s := request{Val: c.val}
		err := v.Struct(s)
		if c.ok {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
