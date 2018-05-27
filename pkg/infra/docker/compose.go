package docker

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func (d *docker) GetDockerComposeVersion() (string, error) {
	data, err := exec.Command("docker-compose", "version", "--short").Output()
	if err != nil {
		return "", errors.Wrap(err, "fail")
	}
	return strings.TrimRight(string(data[:]), "\n"), nil
}
