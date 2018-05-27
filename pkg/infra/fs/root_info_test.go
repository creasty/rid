package fs

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_LocateRoot(t *testing.T) {
	t.Run("no rid dir", func(t *testing.T) {
		fs := New(afero.NewMemMapFs())

		f, ok := fs.LocateRoot("/app")
		assert.Equal(t, false, ok)
		assert.Equal(t, (*RootInfo)(nil), f)
	})

	t.Run("no docker compose file", func(t *testing.T) {
		fs := New(afero.NewMemMapFs())

		fs.MkdirAll("/app/rid", os.ModeDir)

		f, ok := fs.LocateRoot("/app")
		assert.Equal(t, false, ok)
		assert.Equal(t, (*RootInfo)(nil), f)
	})

	t.Run("find at the current dir", func(t *testing.T) {
		fs := New(afero.NewMemMapFs())

		fs.MkdirAll("/app/rid", os.ModeDir)
		fs.WriteFile("/app/rid/docker-compose.yml", []byte(""), 0666)

		f, ok := fs.LocateRoot("/app")
		assert.Equal(t, true, ok)
		assert.Equal(t, &RootInfo{
			RootDir:     "/app",
			RidDir:      "/app/rid",
			ComposeFile: "/app/rid/docker-compose.yml",
		}, f)
	})

	t.Run("find at a parent dir", func(t *testing.T) {
		fs := New(afero.NewMemMapFs())

		fs.MkdirAll("/app/rid", os.ModeDir)
		fs.WriteFile("/app/rid/docker-compose.yml", []byte(""), 0666)

		fs.MkdirAll("/app/foo/bar/baz", os.ModeDir)

		f, ok := fs.LocateRoot("/app/foo/bar/baz")
		assert.Equal(t, true, ok)
		assert.Equal(t, &RootInfo{
			RootDir:     "/app",
			RidDir:      "/app/rid",
			ComposeFile: "/app/rid/docker-compose.yml",
		}, f)
	})

	t.Run("a verbose extension (.yaml)", func(t *testing.T) {
		fs := New(afero.NewMemMapFs())

		fs.MkdirAll("/app/rid", os.ModeDir)
		fs.WriteFile("/app/rid/docker-compose.yaml", []byte(""), 0666)

		fs.MkdirAll("/app/foo/bar/baz", os.ModeDir)

		f, ok := fs.LocateRoot("/app/foo/bar/baz")
		assert.Equal(t, true, ok)
		assert.Equal(t, &RootInfo{
			RootDir:     "/app",
			RidDir:      "/app/rid",
			ComposeFile: "/app/rid/docker-compose.yaml",
		}, f)
	})
}
