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

func (d *docker) Prepare(dir string) error {
	cmd := exec.Command("docker-compose", "up", "-d", "--remove-orphans")
	cmd.Dir = dir
	cmd.Stdin = d.Stdin
	cmd.Stdout = d.Stdout
	cmd.Stderr = d.Stderr
	return cmd.Run()
}

func (d *docker) Exec(cid string, envs []string, name string, args ...string) error {
	dockerArgs := []string{"exec", "-it"}
	{
		for _, e := range envs {
			dockerArgs = append(dockerArgs, "-e", e)
		}

		dockerArgs = append(dockerArgs, cid)
	}

	args = append([]string{name}, args...)
	args = append(dockerArgs, args...)

	cmd := exec.Command("docker", args...)
	cmd.Stdin = d.Stdin
	cmd.Stdout = d.Stdout
	cmd.Stderr = d.Stderr
	return cmd.Run()
}
