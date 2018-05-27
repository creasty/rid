package fs

import (
	"path/filepath"
)

const (
	ridDirName = "rid"
)

var (
	requiredFiles = []string{
		"docker-compose.yml",
		"docker-compose.yaml",
	}
)

// RootInfo holds essential paths which rid interacts with
type RootInfo struct {
	// RootDir is a path for the application's root directory
	RootDir string
	// RidDir is a path for the `rid/` directory that contains Docker related files
	RidDir string
	// ComposeFile is a path for a configuration file of docker-compose
	ComposeFile string
}

func (fs *fileSystem) LocateRoot(baseDir string) (*RootInfo, bool) {
	if baseDir == "." || baseDir == "/" {
		return nil, false
	}

	ridDir := filepath.Join(baseDir, ridDirName)
	if fs.DirExists(ridDir) {
		for _, f := range requiredFiles {
			composeFile := filepath.Join(ridDir, f)

			if fs.Exists(composeFile) {
				return &RootInfo{
					RootDir:     baseDir,
					RidDir:      ridDir,
					ComposeFile: composeFile,
				}, true
			}
		}
	}

	return fs.LocateRoot(filepath.Dir(baseDir))
}
