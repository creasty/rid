package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContext_findConfigFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("on a root dir", func(t *testing.T) {
		c := &Context{}
		if err := c.findConfigFile(wd); err != nil {
			t.Fatalf("it should find: %v", err)
		}
	})

	t.Run("on a sub-dir", func(t *testing.T) {
		c := &Context{}
		if err := c.findConfigFile(filepath.Join(wd, "foo")); err != nil {
			t.Fatalf("it should find: %v", err)
		}
	})

	t.Run("otherwise", func(t *testing.T) {
		c := &Context{}
		if err := c.findConfigFile("/foo/bar"); err == nil {
			t.Fatalf("it should not find any file")
		}
	})
}

func TestContext_findSubstitutions(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	c := &Context{Command: make(map[string]*Command)}
	c.findConfigFile(wd)

	if err := c.findSubstitutions(); err != nil {
		t.Fatalf("it should return an error: %v", err)
	}

	t.Run("non RIC command", func(t *testing.T) {
		cmd := c.Command["sample"]
		if cmd.Name != filepath.Join(c.BaseDir, "libexec", "rid-sample") {
			t.Fatal("it should be located")
		}
		if cmd.RunInContainer {
			t.Fatal("it should not be run in container")
		}
	})

	t.Run("RIC command", func(t *testing.T) {
		cmd := c.Command["sample2"]
		if cmd.Name != filepath.Join("rid", "libexec", "sample2") {
			t.Fatal("it should be located")
		}
		if !cmd.RunInContainer {
			t.Fatal("it should be run in container")
		}
	})
}
