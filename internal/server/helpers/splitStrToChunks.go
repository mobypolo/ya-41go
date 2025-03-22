package helpers

import "strings"

func SplitStrToChunks(str string) []string {
	return strings.FieldsFunc(str, func(r rune) bool { return r == '/' })
}
