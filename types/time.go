package types

import (
	"errors"
	"time"
)

const TimeLayout = "15:04"

var errTimeParseFormat = errors.New(`TimeParseError: should be a string formatted as "15:04"`)
var errTimeParseEmpty = errors.New("TimeParseError: should be not empty")

func NewTime(t *time.Time) Time {
	return Time{t: t}
}

func NewNotEmptyTimeFromString(s string) (Time, error) {
	if s == "" {
		return Time{}, errTimeParseEmpty
	}
	t, err := NewTimeFromString(s)
	if err != nil {
		return Time{}, err
	}
	return t, nil
}

func NewTimeFromString(s string) (Time, error) {
	t := Time{}
	if s == "" {
		return t, nil
	}
	if err := t.decodeStr(s); err != nil {
		return t, err
	}
	return t, nil
}

func NewTimeFromInt(min int16) (Time, error) {
	t := Time{}
	if min == 0 {
		return t, nil
	}
	if err := t.decodeInt(min); err != nil {
		return t, err
	}
	return t, nil
}

type Time struct {
	t *time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.t == nil {
		return []byte(`""`), nil
	}
	return []byte(`"` + t.t.Format(TimeLayout) + `"`), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == `""` {
		return nil
	}
	return t.decodeStr(s[1:6])
}

func (t *Time) decodeStr(s string) error {
	if len(s) != 5 {
		return errTimeParseFormat
	}
	ret, err := time.Parse(TimeLayout, s)
	if err != nil {
		return err
	}
	t.t = &ret
	return nil
}

func (t *Time) decodeInt(min int16) error {
	if min > 1440 {
		return errTimeMinMax
	}
	dt := time.Date(0, 0, 0, int(min)/60, int(min)%60, 0, 0, time.UTC)
	t.t = &dt
	return nil
}

func (t Time) ToHHMM() string {
	if t.t == nil {
		return ""
	}
	return t.t.Format(TimeLayout)
}

func (t Time) ToMin() int16 {
	if t.t == nil {
		return 0
	}
	return int16(t.t.Hour()*60 + t.t.Minute())
}

func (t Time) ToTime() time.Time {
	if t.t == nil {
		return time.Time{}
	}
	return *t.t
}

func (t Time) ToDuration() time.Duration {
	if t.t == nil {
		return 0
	}
	return time.Minute * time.Duration(t.ToMin())
}

func (t Time) After(dt time.Time) bool {
	if t.t == nil {
		return false
	}
	return int16(dt.Hour()*60+dt.Minute()) < t.ToMin()
}

func (t Time) Before(dt time.Time) bool {
	if t.t == nil {
		return true
	}
	return int16(dt.Hour()*60+dt.Minute()) >= t.ToMin()
}

func (t Time) Equal(dt *time.Time) bool {
	if t.t == nil {
		return dt == nil
	}
	if dt == nil {
		return t.t == nil
	}
	return t.t.Equal(*dt)
}
