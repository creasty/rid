package util

import (
	"strings"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	if GetLocalIP() == "" {
		t.Fatal("it should retrieve a local IP")
	}
}

func TestRemovePrefix(t *testing.T) {
	t.Run("without prefix", func(t *testing.T) {
		str, ok := RemovePrefix("prefix-", "aaa-foo")
		if ok {
			t.Fatal("it should return false")
		}
		if str != "aaa-foo" {
			t.Fatal("it should return the given string as is")
		}
	})

	t.Run("with prefix", func(t *testing.T) {
		str, ok := RemovePrefix("prefix-", "prefix-foo")
		if !ok {
			t.Fatal("it should return true")
		}
		if str != "foo" {
			t.Fatal("it should trim a prefix")
		}
	})
}

func TestLoadHelpFile(t *testing.T) {
	summary, description := LoadHelpFile("../testdata/help.txt")

	if summary != "Summary line" {
		t.Fatal("a summary line should be parsed")
	}

	if !strings.Contains(description, summary) {
		t.Fatal("description should contains the summary line")
	}

	if !strings.Contains(description, "Description goes here") {
		t.Fatal("description should have full contents")
	}
}
