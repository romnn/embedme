package fs

import (
	iofs "io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

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

// if the filesystem supports it, use Lstat, else use fs.Stat
func lstatIfPossible(fs afero.Fs, path string) (os.FileInfo, error) {
	if lfs, ok := fs.(afero.Lstater); ok {
		fi, _, err := lfs.LstatIfPossible(path)
		return fi, err
	}
	return fs.Stat(path)
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

func legacyGlob(fs afero.Fs, pattern string) (matches []string, err error) {
	log.Printf("glob: pattern=%s hasMeta=%v\n", pattern, hasMeta(pattern))
	if !hasMeta(pattern) {
		// Lstat not supported by all filesystems.
		if _, err = lstatIfPossible(fs, pattern); err != nil {
			return nil, nil
		}
		return []string{pattern}, nil
	}

	dir, file := filepath.Split(pattern)
	switch dir {
	case "":
		dir = "."
	case string(filepath.Separator):
	// nothing
	default:
		dir = dir[0 : len(dir)-1] // chop off trailing separator
	}

	if !hasMeta(dir) {
		return glob(fs, dir, file, nil)
	}

	var m []string
	m, err = Glob(fs, dir)
	if err != nil {
		return
	}
	for _, d := range m {
		matches, err = glob(fs, d, file, matches)
		if err != nil {
			return
		}
	}
	return
}

// glob searches for files matching pattern in the directory dir
// and appends them to matches. If the directory cannot be
// opened, it returns the existing matches. New matches are
// added in lexicographical order.
func glob(fs afero.Fs, dir, pattern string, matches []string) (m []string, e error) {
	m = matches
	fi, err := fs.Stat(dir)
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return
	}
	d, err := fs.Open(dir)
	if err != nil {
		return
	}
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	sort.Strings(names)

	for _, n := range names {
		matched, err := filepath.Match(pattern, n)
		if err != nil {
			return m, err
		}
		if matched {
			m = append(m, filepath.Join(dir, n))
		}
	}
	return
}

// hasMeta reports whether path contains any of the magic characters
// recognized by Match.
func hasMeta(path string) bool {
	return strings.ContainsAny(path, "*?[")
}
