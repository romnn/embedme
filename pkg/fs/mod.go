package fs

import (
	"fmt"

	"github.com/spf13/afero"
)

// DirFS returns a filesystem that is rooted at a given path
func DirFS(fs afero.Fs, path string) afero.Fs {
	return afero.NewBasePathFs(fs, path)
}

// EnsureFile ensures the presence of a file for a path
func EnsureFile(fs afero.Fs, path string) error {
	stat, err := fs.Stat(path)
	if err != nil {
		return fmt.Errorf("file %s does not exist", path)
	}
	if stat.IsDir() {
		return fmt.Errorf("file %s is a directory", path)
	}
	return nil
}
