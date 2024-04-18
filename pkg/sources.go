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

// SourceMap is a map that tracks if a source path is valid (non-ignored)
type SourceMap map[string]bool

// ValidCount counts the number of valid (non-ignored) sources
func (s SourceMap) ValidCount() int {
	valid := 0
	for _, ok := range s {
		if ok {
			valid++
		}
	}
	return valid
}

// Ignore ignores all sources that match the entries of the Gitignore
//
// Note: Ignore should be called after all sources were added
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

// Add adds a source to the soure map
func (s SourceMap) Add(sources ...string) {
	for _, source := range sources {
		s[source] = true
	}
}

// Valid returns all valid (non-ignored) sources
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
	WorkingDir  string
	Glob        bool
	IgnoreFiles []string
}

// NewSourceFinder ...
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
	for _, pattern := range patterns {
		if f.Glob {
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

	for _, ignoreFile := range f.IgnoreFiles {
		ignorePath := filepath.Join(f.WorkingDir, ignoreFile)
		if _, err := fs.Stat(ignorePath); err != nil {
			return sources, err
		}

		ignoreContent, err := afero.ReadFile(fs, ignorePath)
		if err != nil {
			return sources, err
		}
		ignoreLines := strings.Split(string(ignoreContent), "\n")

		ignore := ignore.CompileIgnoreLines(ignoreLines...)
		ignored := sources.Ignore(ignore)

		if ignored > 0 {
			Info(
				log.Writer(),
				"Skipped %d file(s) ignored in %s\n",
				ignored, ignoreFile,
			)
		}
	}
	return sources, nil
}
