package types

import (
	"regexp"
	"strings"
)

var fileRE = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9 \(\)\.,+!]`)
var nl2brReplacer = strings.NewReplacer("\n", "<br/>")

func NlToBr(txt string) string {
	return nl2brReplacer.Replace(txt)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func CamelToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func SanitizeBaseName(file string) string {
	return fileRE.ReplaceAllString(file, "")
}
