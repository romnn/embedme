package fs

import (
	iofs "io/fs"
	"regexp"

	"github.com/spf13/afero"
)

var replaces = regexp.MustCompile(`(\.)|(\*\*\/)|(\*)|([^\/\*]+)|(\/)`)

func toRegexp(pattern string) string {
	pat := replaces.ReplaceAllStringFunc(pattern, func(s string) string {
		switch s {
		case "/":
			return "\\/"
		case ".":
			return "\\."
		case "**/":
			return ".*"
		case "*":
			return "[^/]*"
		default:
			return s
		}
	})
	return "^" + pat + "$"
}

// Glob returns a list of files matching the pattern.
// The pattern can include **/ to match any number of directories.
func Glob(fs afero.Fs, pattern string) ([]string, error) {
	files := []string{}

	regexpPat := regexp.MustCompile(toRegexp(pattern))

	err := afero.Walk(fs, ".", func(path string, d iofs.FileInfo, err error) error {
		if d.IsDir() || err != nil {
			return nil
		}
		if regexpPat.MatchString(path) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
