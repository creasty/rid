package fs

import (
	"github.com/spf13/afero"
)

func (fs *fileSystem) Exists(path string) bool {
	ok, err := afero.Exists(fs.Fs, path)
	return ok && err == nil
}

func (fs *fileSystem) DirExists(path string) bool {
	ok, err := afero.DirExists(fs.Fs, path)
	return ok && err == nil
}
