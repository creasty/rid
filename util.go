package main

import (
	"os/exec"
	"strings"
)

const (
	WrapperPrefix = "devc-"
)

func getLocalIP() string {
	for _, i := range []string{"en0", "en1", "en2"} {
		cmd := exec.Command("ipconfig", "getifaddr", i)
		b, err := cmd.Output()
		if err != nil {
			continue
		}

		if len(b) > 0 {
			return strings.Trim(string(b[:]), "\n")
		}
	}

	return ""
}

func removeWrapperPrefix(str string) (string, bool) {
	had := strings.HasPrefix(str, WrapperPrefix)
	return strings.TrimPrefix(str, WrapperPrefix), had
}
