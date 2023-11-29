package lhdiff

import (
	"regexp"
	"strings"
)

func ConvertToLinesWithoutNewLine(text string) []string {
	if text == "" {
		return make([]string, 0)
	}
	lines := strings.SplitAfter(text, "\n")
	return conv(lines, removeMultipleSpaceAndTrim)
}

func conv(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		res := f(v)
		vsm[i] = res
	}
	return vsm
}

func removeMultipleSpaceAndTrim(s string) string {
	re := regexp.MustCompile("[ \t]+")
	return strings.TrimSpace(re.ReplaceAllString(s, " ")) + "\n"
}
