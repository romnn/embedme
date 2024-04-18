package fs

import (
	"fmt"
	// "os"

	"github.com/spf13/afero"
	// "io"
)

// func ReadFile(file File) ([]byte, error) {
// 	defer file.Close()
// 	return io.ReadAll(file)
// }

//	func ReadFile(fs FileSystem, path string) ([]byte, error) {
//		file, err := fs.Open(path)
//		if err != nil {
//			return nil, err
//		}
//		defer file.Close()
//		return io.ReadAll(file)
//	}

// bp := afero.NewBasePathFs(afero.NewOsFs(), "/base/path")

func DirFS(fs afero.Fs, path string) afero.Fs {
	return afero.NewBasePathFs(fs, path)
	// return os.DirFS(path)
}

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
