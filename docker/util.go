package docker

import (
	"regexp"
	"strings"
)

var (
	// https://github.com/docker/compose/blob/f55c9d42013e8fbb5285bc402d8248a846485217/compose/cli/command.py#L105
	projectNormalizePattern = regexp.MustCompile(`[^a-z0-9]`)
)

// NormalizeProjectName normalizes a project name
func NormalizeProjectName(str string) string {
	str = strings.ToLower(str)
	str = projectNormalizePattern.ReplaceAllString(str, "")
	return str
}
