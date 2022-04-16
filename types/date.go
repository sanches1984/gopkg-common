package types

import (
	"errors"
	"time"
)

const DateFormatYMD = "2006-01-02"

var errDateParseFormat = errors.New(`TimeParseError: should be a string formatted as "2006-01-02"`)
var errDateParseEmpty = errors.New("TimeParseError: should be not empty")

func NewDate(t *time.Time) Date {
	return Date{t: t}
}

func NewNotEmptyDateFromString(s string) (Date, error) {
	if s == "" {
		return Date{}, errDateParseEmpty
	}
	d, err := NewDateFromString(s)
	if err != nil {
		return Date{}, err
	}
	return d, nil
}

func NewDateFromString(s string) (Date, error) {
	d := Date{}
	if s == "" {
		return d, nil
	}
	if err := d.decode(s); err != nil {
		return d, err
	}
	return d, nil
}

type Date struct {
	t *time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.t == nil {
		return []byte(`""`), nil
	}
	return []byte(`"` + d.t.Format(DateFormatYMD) + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == `""` {
		return nil
	}
	return d.decode(s[1:11])
}

func (d *Date) decode(s string) error {
	if len(s) != 10 {
		return errDateParseFormat
	}
	ret, err := time.Parse(DateFormatYMD, s)
	if err != nil {
		return err
	}
	d.t = &ret
	return nil
}

func (d Date) ToYMD() string {
	if d.t == nil {
		return ""
	}
	return d.t.Format(DateFormatYMD)
}

func (d Date) ToTime() *time.Time {
	return d.t
}

func (d Date) Equal(dt *time.Time) bool {
	if d.t == nil {
		return dt == nil
	}
	if dt == nil {
		return d.t == nil
	}
	return d.t.Equal(*dt)
}

func (d Date) After(dt *time.Time) bool {
	if d.t == nil {
		return false
	}
	if dt == nil {
		return true
	}
	return d.t.After(*dt)
}

func (d Date) Before(dt *time.Time) bool {
	if d.t == nil {
		return false
	}
	if dt == nil {
		return false
	}
	return d.t.Before(*dt)
}
