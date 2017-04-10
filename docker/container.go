package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	dockerInspectFormatter = `{{ .ID }}` +
		` {{ index .Config.Labels "com.docker.compose.project" }}` +
		`:{{ index .Config.Labels "com.docker.compose.service" }}` +
		`:{{ index .Config.Labels "com.docker.compose.container-number" }}`
)

// FindContainerByService returns a container ID for service
func FindContainerByService(projectName, service string, num int) (string, error) {
	projectName = NormalizeProjectName(projectName)
	composeID := fmt.Sprintf("%s:%s:%d", projectName, service, num)

	psStream := new(bytes.Buffer)
	psCmd := exec.Command("docker", "ps", "-q")
	psCmd.Stdout = psStream
	if err := psCmd.Run(); err != nil {
		return "", err
	}

	inspectStream := new(bytes.Buffer)
	inspectCmd := exec.Command("xargs", "docker", "inspect", "--format", dockerInspectFormatter)
	inspectCmd.Stdin = psStream
	inspectCmd.Stdout = inspectStream
	if err := inspectCmd.Run(); err != nil {
		return "", err
	}

	r := bufio.NewReader(inspectStream)

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		s := strings.SplitN(string(line[:]), " ", 2)
		containerID, label := s[0], s[1]
		if label == composeID {
			return containerID, nil
		}
	}

	return "", fmt.Errorf("No container for projectName=%s, service=%s, number=%d", projectName, service, num)
}
