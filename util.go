package main

import (
	"io/ioutil"
	"os"
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

func loadHelpFile(file string) (summary, description string) {
	if file == "" {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	description = string(b[:])
	summary = strings.SplitN(description, "\n", 2)[0] // FIXME: consider other newline chars
	return
}
