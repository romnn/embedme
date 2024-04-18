package embedme

import (
	"fmt"
	"log"
	"regexp"

	"github.com/romnn/embedme/internal"
	"github.com/romnn/embedme/pkg/commands"
	"github.com/spf13/afero"
)

var (
	// Match a code block
	// optional: capture groups around the file extension
	// optional: capture first line starting with //
	blockRegex = regexp.MustCompile(
		// multiline mode
		// "(?m:)(?P<all>" +
		"(?m:(?P<all>" +
			// optional embed comment
			"(<!--\\s*?embedme[ ]+?(?P<embedComment>\\S+?)\\s*?-->)?" +
			// optional ignore next comment
			"(?P<embedIgnore><!--\\s*?embedme[ -]ignore-next\\s*?-->)?" +
			// [\\s\\S]*?)?" +
			// todo: ensure this is not another code block here
			"[\\s\\S]*?" +
			// start of block and language
			// indent counts the number of whitespace and tabs for indentation
			"^(?P<indent>[ \t]*?)```(?P<language>\\w*)?\\s*\n" +
			// "^(?P<indent>\\s*```(?P<language>\\w*)?.*\n" +
			// inside code block
			// "(?P<block>[\\s\\S]*?)^[ \t]*?```" +
			"(?P<block>[\\s\\S]*?)```" +
			// end of multiline mode
			"))",
	)
)

// CodeBlock ...
type CodeBlock struct {
	Start        int
	End          int
	StartLine    int
	EndLine      int
	Indent       string
	Code         string
	EmbedComment string
	Ignore       bool
	Language     LanguageID
}

// Comment ...
func (b *CodeBlock) Comment() string {
	typ, err := b.CommentType()
	if err != nil {
		return ""
	}
	return commentString(*typ)
}

// CommentType ...
func (b *CodeBlock) CommentType() (*CommentType, error) {
	comment, ok := CommentForLanguage[b.Language]
	if !ok {
		return nil, fmt.Errorf(
			"unsupported file extension %q (must be one of %+v)",
			b.Language,
			SupportedLanguages,
		)
	}
	return &comment, nil
}

// EmbedCommand ...
func (b *CodeBlock) EmbedCommand(fs afero.Fs, options *Options) (string, commands.Command, error) {
	typ, err := b.CommentType()
	if err != nil {
		return "", nil, err
	}

	embedComment := b.EmbedComment
	if embedComment == "" {
		if firstEmbedComment, ok := FirstComment(b.Code, *typ); ok {
			embedComment = firstEmbedComment
		} else {
			return "", nil, fmt.Errorf(
				"no comment starting with %s in first line of %q block",
				commentString(*typ),
				b.Language,
			)
		}
	}

	baseDirs := []string{options.Base, options.WorkingDir}
	commands := []commands.Command{
		commands.NewEmbedFileCommand(fs, baseDirs...),
		commands.NewEmbedCommandOutputCommand(fs, options.WorkingDir),
	}
	for _, cmd := range commands {
		if err := cmd.Parse(embedComment); err == nil {
			return embedComment, cmd, nil
		}
	}
	return embedComment, nil, fmt.Errorf(
		"%q is not a valid command", embedComment,
	)
}

// ExtractCodeBlocks ...
func ExtractCodeBlocks(source string) []CodeBlock {
	var blocks []CodeBlock

	matches := internal.GetMatches(blockRegex, string(source))
	for _, match := range matches {
		log.Printf("match: %+v\n\n", match)
		var block CodeBlock
		if b, ok := match["block"]; ok {
			block.Start = b.Start
			block.End = b.End

			block.Code = b.Text
			block.StartLine = internal.LineNumber(source, b.Start)
			block.EndLine = internal.LineNumber(source, b.End)
		}

		if _, ok := match["embedIgnore"]; ok {
			block.Ignore = true
		}
		// if all, ok := match["all"]; ok {
		// 	block.Start = all.Start
		// 	block.End = all.End
		// }
		if indent, ok := match["indent"]; ok {
			block.Indent = indent.Text
		}

		if comment, ok := match["embedComment"]; ok {
			block.EmbedComment = comment.Text
		}
		if language, ok := match["language"]; ok {
			block.Language = LanguageID(language.Text)
		}

		blocks = append(blocks, block)
	}

	return blocks
}

func commentString(typ CommentType) string {
	switch typ {
	case CommentDoubleSlash:
		return "//"
	case CommentXML:
		return "<!-- ... -->"
	case CommentHash:
		return "#"
	case CommentSingleQuote:
		return "'"
	case CommentDoublePercent:
		return "%%"
	case CommentDoubleHyphens:
		return "--"
	}
	return ""
}

func commentPrefixRegex(comment string) *regexp.Regexp {
	// independent of newline type
	return regexp.MustCompile(`^\s*` + comment + `([\s\S]*?)\r?\n`)
}

// FirstComment ...
func FirstComment(source string, typ CommentType) (string, bool) {
	var re *regexp.Regexp
	switch typ {
	case CommentNone:
		return "", true
	case CommentDoubleSlash:
		re = commentPrefixRegex("//")
	case CommentXML:
		re = regexp.MustCompile(`<!--\s*?(\S*?)\s*?-->`)
	case CommentHash:
		re = commentPrefixRegex("#")
	case CommentSingleQuote:
		re = commentPrefixRegex("'")
	case CommentDoublePercent:
		re = commentPrefixRegex("%%")
	case CommentDoubleHyphens:
		re = commentPrefixRegex("--")
	}
	match := re.FindStringSubmatch(source)
	if len(match) > 0 {
		return match[1], true
	}
	return "", false
}
