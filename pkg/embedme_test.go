package embedme

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp/v3"
)

func pretty(m interface{}) string {
	return pp.Sprint(m)
}

// todo: test first comment for hash comment (python)
// todo: test first comment for single quote comment (?)
// todo: test first comment for double slash quote comment (C)

func TestFirstCommentHash(t *testing.T) {
	for _, c := range []struct {
		description string
		source      string
		language    []Language
		expected    string
	}{
		{
			description: "single valid python comment",
			language:    Python,
			source: `

# valid python comment

def test():
  pass
      `,
			expected: " valid python comment",
		},
		{
			description: "single invalid python comment",
			language:    Python,
			source: `

// invalid python comment

def test():
  pass
      `,
			expected: "",
		},
		{
			description: "2 python comments: valid then invalid",
			language:    Python,
			source: `

  # valid python comment
// invalid python comment

def test():
  pass
      `,
			expected: " valid python comment",
		},
		{
			description: "2 python comments: invalid then valid",
			language:    Python,
			source: `

// invalid python comment
  # valid python comment

def test():
  pass
      `,
			expected: "",
		},
		{
			description: "2 valid python comments",
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

		options := []cmp.Option{
			cmpopts.IgnoreUnexported(CodeBlock{}),
		}
		equal := cmp.Equal(comment, c.expected, options...)
		diff := cmp.Diff(comment, c.expected, options...)
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
			description: "2 python code blocks with valid and invalid comments",
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

		options := []cmp.Option{
			cmpopts.IgnoreUnexported(CodeBlock{}),
		}
		equal := cmp.Equal(blocks, c.expected, options...)
		diff := cmp.Diff(blocks, c.expected, options...)
		if !equal {
			t.Log(pretty(blocks))
			t.Log(pretty(c.expected))
			t.Log(c.source)
			t.Fatalf("%s: unexpected code blocks: %s", c.description, diff)
		}
	}
}
