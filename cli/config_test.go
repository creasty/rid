package cli

import (
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		_, err := NewConfig("./testdata/config_invalid.yml")
		if err == nil || !strings.Contains(err.Error(), "ProjectName") {
			t.Fatal("it should validate a presence of ProjectName")
		}
	})

	t.Run("regular config", func(t *testing.T) {
		c, err := NewConfig("./testdata/config_valid.yml")
		if err != nil {
			t.Fatal("it should be valid")
		}

		if c.ProjectName != "foo/bar" {
			t.Fatal("ProjectName should be set from the file")
		}

		if c.MainService != DefaultMainService {
			t.Fatal("MainService should be a default value")
		}
	})

	t.Run("custom config", func(t *testing.T) {
		c, err := NewConfig("./testdata/config_custom.yml")
		if err != nil {
			t.Fatal("it should be valid")
		}

		if c.ProjectName != "foo/bar" {
			t.Fatal("ProjectName should be set from the file")
		}

		if c.MainService != "apple" {
			t.Fatal("MainService should be set from the file")
		}
	})
}
