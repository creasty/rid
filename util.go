package main

import (
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	newlinePattern = regexp.MustCompile("\r\n|\r|\n")
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}

		if ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() == nil {
			continue
		}

		return ipnet.IP.String()
	}

	return ""
}

func removePrefix(prefix, str string) (string, bool) {
	had := strings.HasPrefix(str, prefix)
	return strings.TrimPrefix(str, prefix), had
}

func loadHelpFile(file string) (summary, description string) {
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
