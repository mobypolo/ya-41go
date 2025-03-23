package utils

import "strings"

func ExistInArrayWithWildCard(arr []string, s string) bool {
	for _, pattern := range arr {
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(s, prefix) {
				return true
			}
		} else {
			if pattern == s {
				return true
			}
		}
	}
	return false
}
