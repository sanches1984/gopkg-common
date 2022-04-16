package types

import "time"

const dateTimeDBLayout = "2006-01-02 15:04:05.999999999 -0700 MST"
const dateTimeYMDThmsLayout = "20060102T150405"
const dateTimeYMDhmLayout = "2006-01-02 15:04"

func DateTimeToDBString(date time.Time) string {
	return date.Format(dateTimeDBLayout)
}

func DateTimeToYMDTHms(date time.Time) string {
	return date.Format(dateTimeYMDThmsLayout)
}

func DateTimeToYMDHm(date time.Time) string {
	return date.Format(dateTimeYMDhmLayout)
}
