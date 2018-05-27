//go:generate mockgen -source=fs.go -package fs -destination=fs_mock.go

package fs

import (
	"os"
	"time"

	"github.com/spf13/afero"
)

// FileSystem is the filesystem interface.
type FileSystem interface {
	// afero.Fs
	Create(name string) (afero.File, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (afero.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (afero.File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(name string) (os.FileInfo, error)
	Name() string
	Chmod(name string, mode os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error

	// ioutil functions in afero.Afero
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	TempFile(dir, prefix string) (f afero.File, err error)
	TempDir(dir, prefix string) (name string, err error)

	// LocateRoot looks for the application's root directory upward from the specified directory.
	// Returns false when it couldn't locate the directory.
	LocateRoot(baseDir string) (*RootInfo, bool)
}

// New creates a file system instance
func New(fs afero.Fs) FileSystem {
	return &fileSystem{
		Afero: afero.Afero{Fs: fs},
	}
}

// NewTest creates a file system instance for testing
func NewTest() FileSystem {
	return New(afero.NewMemMapFs())
}

type fileSystem struct {
	afero.Afero
}
