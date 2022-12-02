package internal

import (
  "os"
  "fmt"
)

func EnsureFile(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file %s does not exist", path)
	}
	if stat.IsDir() {
		return fmt.Errorf("file %s is a directory", path)
	}
	return nil
}
