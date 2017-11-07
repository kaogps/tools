package utils

import (
	"regexp"
)

// TrimPathName //seg1////seg2 --> /seg1/seg2
func TrimPathName(s string) string {
	if len(s) == 0 {
		return ""
	}
	var reg = regexp.MustCompile(`/+`)
	return reg.ReplaceAllString(s, "/")
}
