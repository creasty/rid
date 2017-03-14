package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	WrapperPrefix = "devc-"
)

var (
	newlinePattern = regexp.MustCompile("\r\n|\r|\n")
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
	summary = newlinePattern.Split(description, 2)[0]
	return
}
