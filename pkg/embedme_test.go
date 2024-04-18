package embedme

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/afero"
)

func pretty(m interface{}) string {
	return pp.Sprint(m)
}

func TestFirstCommentHash(t *testing.T) {
	for _, c := range []struct {
		description string
		source      string
		language    []LanguageID
		expected    string
	}{
		{
			description: "single valid Python comment",
			language:    Python,
			source: `

# valid Python comment

def test():
  pass
      `,
			expected: " valid Python comment",
		},
		{
			description: "single invalid Python comment",
			language:    Python,
			source: `

// invalid Python comment

def test():
  pass
      `,
			expected: "",
		},
		{
			description: "2 Python comments: valid then invalid",
			language:    Python,
			source: `

  # valid Python comment
// invalid Python comment

def test():
  pass
      `,
			expected: " valid Python comment",
		},
		{
			description: "2 Python comments: invalid then valid",
			language:    Python,
			source: `

// invalid Python comment
  # valid Python comment

def test():
  pass
      `,
			expected: "",
		},
		{
			description: "2 valid Python comments",
			language:    Python,
			source: `
  # valid 1
  # valid 2

def test():
  pass
      `,
			expected: " valid 1",
		},
	} {
		lang := c.language[0]
		commentTyp, ok := CommentForLanguage[lang]
		if !ok {
			t.Fatalf("Unsupported file extension %q", lang)
		}
		comment, _ := FirstComment(c.source, commentTyp)

		equal := cmp.Equal(comment, c.expected)
		diff := cmp.Diff(comment, c.expected)
		if !equal {
			t.Log(pretty(comment))
			t.Log(pretty(c.expected))
			t.Log(c.source)
			t.Fatalf("%s: unexpected first comment: %s", c.description, diff)
		}
	}
}

func TestExtractCodeBlocks(t *testing.T) {
	for _, c := range []struct {
		description string
		source      string
		expected    []CodeBlock
	}{
		{
			description: "2 code blocks without language",
			source: `
This is a regular readme

` + "```" + `
block1
` + "```" + `

` + "```" + `
block2
` + "```" + `

      `,
			expected: []CodeBlock{
				{
					Code:      "block1\n",
					Start:     31,
					End:       38,
					StartLine: 5,
					EndLine:   6,
				},
				{
					Code:      "block2\n",
					Start:     47,
					End:       54,
					StartLine: 9,
					EndLine:   10,
				},
			},
		},
		{
			description: "2 blocks with language",
			source: `
This is a regular readme

` + "```" + `python
block1
` + "```" + `

` + "```" + `python
block2
` + "```" + `

      `,
			expected: []CodeBlock{
				{
					Code:      "block1\n",
					Language:  "python",
					Start:     37,
					End:       44,
					StartLine: 5,
					EndLine:   6,
				},
				{
					Code:      "block2\n",
					Language:  "python",
					Start:     59,
					End:       66,
					StartLine: 9,
					EndLine:   10,
				},
			},
		},
		{
			description: "2 Python code blocks with valid and invalid comments",
			source: `
This is a regular readme

` + "```" + `python
// bad comment
` + "```" + `

` + "```" + `python
# good comment
` + "```" + `

      `,
			expected: []CodeBlock{
				{
					Code:      "// bad comment\n",
					Language:  "python",
					Start:     37,
					End:       52,
					StartLine: 5,
					EndLine:   6,
				},
				{
					Code:      "# good comment\n",
					Language:  "python",
					Start:     67,
					End:       82,
					StartLine: 9,
					EndLine:   10,
				},
			},
		},
	} {
		blocks := ExtractCodeBlocks(c.source)

		equal := cmp.Equal(blocks, c.expected)
		diff := cmp.Diff(blocks, c.expected)
		if !equal {
			t.Log(pretty(blocks))
			t.Log(pretty(c.expected))
			t.Log(c.source)
			t.Fatalf("%s: unexpected code blocks: %s", c.description, diff)
		}
	}
}

func TestEmbedFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	workingDir := "/work/dir"
	fs.MkdirAll(workingDir, 0755)

	readmePath := filepath.Join(workingDir, "readme.md")
	readme := strings.TrimSpace(`
## This is a readme

#### Python code embedding
` + "```" + `python
# code/python.py
` + "```" + `
#### Rust code embedding
` + "```" + `rust
// code/rust.rs
` + "```" + `
#### Go code embedding
` + "```" + `rust
// go.go
` + "```" + `
	`)
	afero.WriteFile(fs, readmePath, []byte(readme), 0644)

	afero.WriteFile(fs, "/work/dir/.gitignore", []byte(strings.TrimSpace(`
ignored/*
	`)), 0644)

	afero.WriteFile(fs, "/work/dir/nested/.gitignore", []byte(strings.TrimSpace(`
values.txt
	`)), 0644)
	afero.WriteFile(fs, "/work/dir/nested/values.txt", []byte(""), 0644)

	// add code files
	afero.WriteFile(fs, "/work/dir/code/python.py", []byte(strings.TrimSpace(`
def greet(name):
	print(f"Hello {name}")
	`)), 0644)
	afero.WriteFile(fs, "/work/dir/code/rust.rs", []byte(strings.TrimSpace(`
pub fn greet(name: &str) {
	println!("Hello {}", name);
}
	`)), 0644)
	afero.WriteFile(fs, "/work/dir/code/go.go", []byte(strings.TrimSpace(`
package main
import (
	"fmt"
)
func greet(name string) {
	fmt.Printf("Hello %s\n", name)
}
	`)), 0644)

	options := NewDefaultOptions()
	options.Base = filepath.Join(workingDir, "code")
	options.WorkingDir = workingDir
	embedder := Embedder{
		Options: options,
		FS:      fs,
	}
	ignoreFiles, err := GlobFiles(
		fs, workingDir, "**/.gitignore",
	)
	if err != nil {
		t.Fatalf("failed to find ignore files: %v", err)
	}
	t.Logf("ignore files: %v", pretty(ignoreFiles))
	expectedIgnoreFiles := []string{".gitignore", "nested/.gitignore"}
	equal := cmp.Equal(ignoreFiles, expectedIgnoreFiles)
	diff := cmp.Diff(ignoreFiles, expectedIgnoreFiles)
	if !equal {
		t.Log(pretty(expectedIgnoreFiles))
		t.Fatalf("unexpected ignore files: %s", diff)
	}

	finder := SourceFinder{
		WorkingDir:  workingDir,
		Glob:        true,
		IgnoreFiles: ignoreFiles,
	}
	sources, err := finder.FindSources(fs, "**/*.md", "**/*.txt")
	if err != nil {
		t.Fatalf("failed to find sources: %v", err)
	}
	t.Logf("sources: %v", pretty(sources))
	expectedSources := SourceMap{
		"nested/values.txt": false,
		"readme.md":         true,
	}
	equal = cmp.Equal(sources, expectedSources)
	diff = cmp.Diff(sources, expectedSources)
	if !equal {
		t.Log(pretty(expectedSources))
		t.Fatalf("unexpected sources: %s", diff)
	}

	for i, source := range sources.Valid() {
		if err := embedder.ProcessSource(i, source); err != nil {
			t.Fatalf("%v", err)
		}
	}
	embeddedReadme, err := afero.ReadFile(fs, readmePath)
	if err != nil {
		t.Fatalf("failed to read embedded file %s: %v", readmePath, err)
	}
	t.Logf("original readme: %s", readme)
	t.Logf("embedded readme: %s", string(embeddedReadme))

	expectedReadme := strings.TrimSpace(`
## This is a readme

#### Python code embedding
` + "```" + `python
# code/python.py

def greet(name):
	print(f"Hello {name}")
` + "```" + `
#### Rust code embedding
` + "```" + `rust
// code/rust.rs

pub fn greet(name: &str) {
	println!("Hello {}", name);
}
` + "```" + `
#### Go code embedding
` + "```" + `rust
// go.go

package main
import (
	"fmt"
)
func greet(name string) {
	fmt.Printf("Hello %s\n", name)
}
` + "```" + `
	`)

	equal = cmp.Equal(string(embeddedReadme), expectedReadme)
	diff = cmp.Diff(string(embeddedReadme), expectedReadme)
	if !equal {
		t.Logf("expected readme: %s", expectedReadme)
		t.Fatalf("unexpected readme: %s", diff)
	}

}

func TestEmbedCommands(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("src/a", 0755)
	afero.WriteFile(appFS, "src/a/b", []byte("file b"), 0644)
	afero.WriteFile(appFS, "src/c", []byte("file c"), 0644)
}
