package embedme

import (
	// "errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/romnn/embedme/internal"
	"github.com/romnn/embedme/pkg/commands"
)

var (
	// Match a code block
	// optional: capture groups around the file extension
	// optional: capture first line starting with //
	blockRegex = regexp.MustCompile(
		// multiline mode
		"(?m:" +
			// optional embed comment
			"(<!--\\s*?embedme[ ]+?(?P<embedComment>\\S+?)\\s*?-->)?" +
			// optional ignore next comment
			"(?P<embedIgnore><!--\\s*?embedme[ -]ignore-next\\s*?-->)?" +
			// [\\s\\S]*?)?" +
			// todo: ensure this is not another code block here
			"[\\s\\S]*?" +
			// start of block and language
			"^(?P<indent>[ \t]*?)```(?P<language>\\w*)?.*\n" +
			// inside code block
			"(?P<block>[\\s\\S]*?)^[ \t]*?```" +
			// end of multiline mode
			")",
	)
)

type CodeBlock struct {
	Start     int
	End       int
	StartLine int
	EndLine   int
	Indent    string
	Code      string
	// embedComment *EmbedComment
	EmbedComment string
	Ignore       bool
	Language     Language
}

func (b *CodeBlock) Comment() string {
	typ, err := b.CommentType()
	if err != nil {
		return ""
	}
	return commentString(*typ)
}

func (b *CodeBlock) CommentType() (*CommentType, error) {
	// language, err := b.Language
	// language, err := b.Language
	// if err != nil {
	// 	return nil, err
	// }
	comment, ok := CommentForLanguage[b.Language]
	if !ok {
		return nil, fmt.Errorf(
			`Unsupported file extension %q.
Supported extensions are %s, skipping...`,
			b.Language,
			strings.Join(SupportedLanguages, ", "),
		)
	}
	return &comment, nil
}

// make a none language and avoid those errors here
// func (b *CodeBlock) Language() (Language, error) {
// func (b *CodeBlock) Language() Language {
//   return Language(b.language)
// 	// if b.language == "" {
// 	// 	// return NONE[0]
// 	// 	// return NONE[0], errors.New("No code extension detected, skipping ...")
// 	// }
// 	// return Language(b.language)
// }

// type EmbedComment struct {
// 	Original string
// 	Command  string
// }

// func (b *CodeBlock) EmbedCommand(options *Options) (*EmbedComment, commands.Command, error) {
func (b *CodeBlock) EmbedCommand(options *Options) (string, commands.Command, error) {
	// language, err := b.Language()
	// if err != nil {
	// 	return nil, err
	// }
	typ, err := b.CommentType()
	if err != nil {
		return "", nil, err
	}
	// language, _ := b.Language()
	embedComment := b.EmbedComment
	if embedComment == "" {
		var ok bool
		embedComment, ok = FirstComment(b.Code, *typ)
		if !ok {
			return "", nil, fmt.Errorf(
				"no comment starting with %s in first line of %q block",
				commentString(*typ),
				b.Language,
			)
		}
	}
	// if embedComment.Command == "" {
	// 	return nil, nil, fmt.Errorf(
	// 		"no command in first commentwith %s in first line of %q block",
	// 		language,
	// 	)
	// }
	fmt.Println(embedComment)

	commands := []commands.Command{
		commands.NewEmbedFileCommand(options.Base),
		commands.NewEmbedCommandOutputCommand(options.Cwd),
	}
	for _, cmd := range commands {
		// if err := cmd.Parse(embedComment.Command); err == nil {
		if err := cmd.Parse(embedComment); err == nil {
			return embedComment, cmd, nil
		}
	}
	// matches := getMatches(embedPathRegex, embedComment)
	// fmt.Println(matches)
	// check if is filename
	// todo: match embedCommand against different regexes

	return embedComment, nil, fmt.Errorf("%q is not a valid command", embedComment)
}

func ExtractCodeBlocks(source string) []CodeBlock {
	var blocks []CodeBlock

	// newline := internal.DetectNewline([]byte(source))

	matches := internal.GetMatches(blockRegex, string(source))
	for _, match := range matches {
		var block CodeBlock
		if b, ok := match["block"]; ok {
			block.Start = b.Start
			block.End = b.End
			block.Code = b.Text
			// newline)
			block.StartLine = internal.LineNumber(source, b.Start)
			// newline)
			block.EndLine = internal.LineNumber(source, b.End)
		}

		// block.start = match.Start
		// block.end = match.End
		// block.Code = match.Text
		// block.StartLine = LineNumber(source, match.Start, newline)
		// block.EndLine = LineNumber(source, match.End, newline)

		if _, ok := match["embedIgnore"]; ok {
			block.Ignore = true
		}

		if indent, ok := match["indent"]; ok {
			block.Indent = indent.Text
		}

		if comment, ok := match["embedComment"]; ok {
			block.EmbedComment = comment.Text
			// block.embedComment = &EmbedComment{
			// 	Command:  comment.Text,
			// 	Original: comment.Text,
			// }
		}
		if language, ok := match["language"]; ok {
			block.Language = Language(language.Text)
		}
		blocks = append(blocks, block)
	}
	return blocks

	// matches := blockRegex.FindAllStringSubmatchIndex(string(source), -1)
	// for _, pos := range matches {
	// 	// block := match[0]
	// 	// group := make(map[string]string)

	// 	// start := pos[0]
	// 	// end := pos[1]

	// 	var block CodeBlock
	// 	// block := CodeBlock{
	// 	// 	Start: start,
	// 	// 	End:   end,
	// 	// }

	// 	// fmt.Printf("%v\n", blockRegex.SubexpNames())
	// 	// fmt.Printf("%v\n", pos)

	// 	for _, name := range blockRegex.SubexpNames() {
	// 		i := blockRegex.SubexpIndex(name)
	// 		// fmt.Printf("%s %d\n", name, i)
	// 		if i < 0 {
	// 			// no match
	// 			continue
	// 		}
	// 		start := pos[i*2+0]
	// 		end := pos[i*2+1]
	// 		// fmt.Printf("%d %d\n", start, end)
	// 		if start < 0 || end < 0 {
	// 			// no match
	// 			continue
	// 		}
	// 		match := source[start:end]

	// 		switch name {
	// 		case "block":
	// 			block.start = start
	// 			block.end = end
	// 			block.Code = match
	// 			block.StartLine = LineNumber(source, start, newline)
	// 			block.EndLine = LineNumber(source, end, newline)
	// 		case "embedComment":
	// 			block.embedComment = match
	// 		case "language":
	// 			block.language = match
	// 		default:
	// 		}
	// 		// if i != 0 && name != "" {
	// 		// 	// group[name] = match[i]
	// 		// }
	// 	}
	// 	blocks = append(blocks, block)
	// }
	// return blocks
}

// var (
// 	newlineRe = regexp.MustCompile("\r\n")
// 	// strings.Split(strings.ReplaceAll(windows, "\r\n", "\n"), "\n")
// )

// const leadingSymbol = (symbol: string): FilenameFromCommentReader => line => {
//   const regex = new RegExp(k

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

// func FirstComment(source string, typ CommentType) (*EmbedComment, bool) {
func FirstComment(source string, typ CommentType) (string, bool) {
	var re *regexp.Regexp
	switch typ {
	case CommentNone:
		// return &EmbedComment{Command: "", Original: ""}, true
		return "", true
	case CommentDoubleSlash:
		re = commentPrefixRegex("//")
		// return "", nil
	case CommentXML:
		re = regexp.MustCompile(`<!--\s*?(\S*?)\s*?-->`)
		// const match = line.match(/<!--\s*?(\S*?)\s*?-->/);
		// return "", nil
	case CommentHash:
		re = commentPrefixRegex("#")
		// return "", nil
	case CommentSingleQuote:
		re = commentPrefixRegex("'")
		// return "", nil
	case CommentDoublePercent:
		re = commentPrefixRegex("%%")
		// return "", nil
	case CommentDoubleHyphens:
		re = commentPrefixRegex("--")
		// return "", nil
	}
	// fmt.Printf("comment: %s\n", *typ)
	// fmt.Println(source)
	// fmt.Println(re.String())
	match := re.FindStringSubmatch(source)
	// fmt.Println(match)
	if len(match) > 0 {
		return match[1], true
		// &EmbedComment{
		// Original: match[0],
		// Command:  match[1],
		// }, true
		// return match[1], true
	}
	// i := blockRegex.SubexpIndex(name)
	// if (!match) {
	//   return null;
	// }

	// return match[1];

	return "", false
}
