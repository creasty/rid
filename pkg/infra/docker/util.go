package docker

import (
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

var (
	projectNormalizePatternNewVersion *semver.Constraints

	// @see https://github.com/docker/compose/blob/1.21.0/compose/cli/command.py#L132
	projectNormalizePatternNew = regexp.MustCompile(`[^-_a-z0-9]`)

	// @see https://github.com/docker/compose/blob/1.20.1/compose/cli/command.py#L130
	projectNormalizePatternOld = regexp.MustCompile(`[^a-z0-9]`)
)

func init() {
	if c, err := semver.NewConstraint(">= 1.21.0"); err == nil {
		projectNormalizePatternNewVersion = c
	} else {
		panic(err)
	}
}

func (d *docker) NormalizeProjectName(str string) string {
	str = strings.ToLower(str)
	str = d.projectNormalizePattern().ReplaceAllString(str, "")
	return str
}

func (d *docker) projectNormalizePattern() (pattern *regexp.Regexp) {
	pattern = projectNormalizePatternNew

	verStr, err := d.GetDockerComposeVersion()
	if err != nil {
		return
	}

	ver, err := semver.NewVersion(verStr)
	if err != nil {
		return
	}

	if projectNormalizePatternNewVersion.Check(ver) {
		return
	}

	pattern = projectNormalizePatternOld
	return
}
