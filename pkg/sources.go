package embedme

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	fsutil "github.com/romnn/embedme/pkg/fs"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

type SourceMap map[string]bool

func (s SourceMap) ValidCount() int {
	valid := 0
	for _, ok := range s {
		if ok {
			valid++
		}
	}
	return valid
}

func (s SourceMap) Ignore(gi *ignore.GitIgnore) int {
	before := s.ValidCount()
	for source := range s {
		if gi.MatchesPath(source) {
			s[source] = false
		}
	}
	after := s.ValidCount()
	return before - after
}

func (s SourceMap) Add(sources ...string) {
	for _, source := range sources {
		s[source] = true
	}
}

func (s SourceMap) Valid() []string {
	var valid []string
	for source, ok := range s {
		if ok {
			valid = append(valid, source)
		}
	}
	return valid
}

// SourceFinder ...
type SourceFinder struct {
	WorkingDir string
	Glob       bool
	// Ignore      bool
	IgnoreFiles []string
	// FS afero.Fs
}

func NewSourceFinder() SourceFinder {
	return SourceFinder{
		WorkingDir:  "",
		Glob:        true,
		IgnoreFiles: []string{},
	}
}

// GlobFiles ...
func GlobFiles(fs afero.Fs, workingDir string, patterns ...string) ([]string, error) {
	files := []string{}
	dirFS := fsutil.DirFS(fs, workingDir)
	for _, pattern := range patterns {
		matches, err := fsutil.Glob(dirFS, pattern)
		if err != nil {
			return files, fmt.Errorf(
				"failed to glob pattern %q in %s: %v",
				pattern, workingDir, err,
			)
		}
		files = append(files, matches...)
	}
	return files, nil
}

// FindSources ...
func (f *SourceFinder) FindSources(fs afero.Fs, patterns ...string) (SourceMap, error) {
	sources := make(SourceMap)
	// dirFS := fsutil.DirFS(fs, f.WorkingDir)
	for _, pattern := range patterns {
		if f.Glob {
			// matches, err := afero.Glob(dirFS, pattern)
			// if err != nil {
			// 	return sources, fmt.Errorf(
			// 		"failed to glob pattern %q in %s: %v",
			// 		pattern, f.WorkingDir, err,
			// 	)
			// }
			matches, err := GlobFiles(fs, f.WorkingDir, pattern)
			if err != nil {
				return sources, err
			}
			sources.Add(matches...)
		} else {
			if !filepath.IsAbs(pattern) {
				pattern = filepath.Join(f.WorkingDir, pattern)
			}
			sources.Add(pattern)
		}
	}

	// if !f.Ignore {
	// 	return sources, nil
	// }

	// find all the ignore files
	// allIgnoreFiles := []string{}
	// for _, ignorePattern := range []string{".embedmeignore", ".gitignore"} {
	// 	ignoreFiles, err := afero.Glob(dirFS, ignorePattern)
	// 	if err != nil {
	// 		return sources, fmt.Errorf(
	// 			"failed to glob pattern %q in %s: %v",
	// 			ignoreFiles, f.WorkingDir, err,
	// 		)
	// 	}
	// 	allIgnoreFiles = append(allIgnoreFiles, ignoreFiles...)
	// }

	// Log(log.Writer(), "found ignore files: %v\n", allIgnoreFiles)

	for _, ignoreFile := range f.IgnoreFiles {
		ignorePath := filepath.Join(f.WorkingDir, ignoreFile)
		if _, err := fs.Stat(ignorePath); err != nil {
			return sources, err
			// continue
		}

		ignoreContent, err := afero.ReadFile(fs, ignorePath)
		if err != nil {
			return sources, err
		}
		ignoreLines := strings.Split(string(ignoreContent), "\n")

		ignore := ignore.CompileIgnoreLines(ignoreLines...)
		// if err != nil {
		// 	return sources, err
		// 	// continue
		// }
		ignored := sources.Ignore(ignore)

		if ignored > 0 {
			Info(
				log.Writer(), "Skipped %d file(s) ignored in %s\n",
				ignored, ignoreFile,
			)
		}
	}
	return sources, nil
}
