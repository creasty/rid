package util

import (
	"regexp"
	"strings"
)

var (
	newlinePattern = regexp.MustCompile("\r\n|\r|\n")
)

// RemovePrefix trims prefix and returns true if there was a match
func RemovePrefix(prefix, str string) (string, bool) {
	had := strings.HasPrefix(str, prefix)
	return strings.TrimPrefix(str, prefix), had
}

// ParseHelpFile parses a help file
func ParseHelpFile(data []byte) (summary, description string) {
	description = string(data[:])
	summary = newlinePattern.Split(description, 2)[0]
	return
}
