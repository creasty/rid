package docker

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

func init() {
	if !doesAllowDashAndUnderscore() {
		// https://github.com/docker/compose/blob/f55c9d42013e8fbb5285bc402d8248a846485217/compose/cli/command.py#L105
		projectNormalizePattern = regexp.MustCompile(`[^a-z0-9]`)
	}
}

var (
	// https://github.com/docker/compose/blob/1.21.0/compose/cli/command.py#L132
	projectNormalizePattern = regexp.MustCompile(`[^-_a-z0-9]`)
)

// NormalizeProjectName normalizes a project name
func NormalizeProjectName(str string) string {
	str = strings.ToLower(str)
	str = projectNormalizePattern.ReplaceAllString(str, "")
	return str
}

func doesAllowDashAndUnderscore() (allowed bool) {
	allowed = true

	data, err := exec.Command("docker-compose", "version", "--short").Output()
	if err != nil {
		return
	}

	currVer, err := semver.NewVersion(strings.TrimRight(string(data), "\n"))
	if err != nil {
		return
	}

	allowedVer, err := semver.NewVersion("1.21.0")
	if err != nil {
		return
	}

	return currVer.Equal(allowedVer) || currVer.GreaterThan(allowedVer)
}
