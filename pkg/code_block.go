package embedme

import (
	"errors"
	"fmt"
	"github.com/romnn/embedme/internal"
	"github.com/romnn/embedme/pkg/commands"
	"regexp"
	"strings"
)

var (
	// Match a code block
	// optional: capture groups around the file extension
	// optional: capture first line starting with //
	blockRe = regexp.MustCompile(
		// multiline mode
		"(?m:" +
			// optional embed comment
			"(<!--\\s*?embedme[ ]+?(?P<embedComment>\\S+?)\\s*?-->)?" +
			// [\\s\\S]*?)?" +
			// start of block and language
			"[\\s\\S]*?```(?P<language>\\w*)?.*\n" +
			// inside code block
			"(?P<block>[\\s\\S]*?)^[ \t]*?```" +
			// end of multiline mode
			")",
	)
)

type CodeBlock struct {
	start        int
	end          int
	StartLine    int
	EndLine      int
	Code         string
	embedComment string
	language     string
}

func (b *CodeBlock) CommentType() (*CommentType, error) {
	language, err := b.Language()
	if err != nil {
		return nil, err
	}
	comment, ok := CommentForLanguage[language]
	if !ok {
		return nil, fmt.Errorf(
			`Unsupported file extension %q.
Supported extensions are %s, skipping...`,
			language,
			strings.Join(SupportedLanguages, ", "),
		)
	}
	return &comment, nil
}

func (b *CodeBlock) Language() (Language, error) {
	if b.language == "" {
		return Language(""), errors.New("No code extension detected, skipping ...")
	}
	return Language(b.language), nil
}

func (b *CodeBlock) EmbedCommand(options *Options) (commands.Command, error) {
	// language, err := b.Language()
	// if err != nil {
	// 	return nil, err
	// }
	typ, err := b.CommentType()
	if err != nil {
		return nil, err
	}
	embedComment := b.embedComment
	if embedComment == "" {
		embedComment, err = FirstComment(b.Code, *typ)
	}
	if embedComment == "" {
		return nil, nil
		// return nil, fmt.Errorf(
		// 	"No command detected in first line for block with extension %q",
		// 	language,
		// )
	}
	fmt.Println(embedComment)

	commands := []commands.Command{
		commands.NewEmbedFileCommand(options.Base),
		commands.NewEmbedCommandOutputCommand(options.Cwd),
	}
	for _, cmd := range commands {
		if err := cmd.Parse(embedComment); err == nil {
			return cmd, nil
		}
	}
	// matches := getMatches(embedPathRegex, embedComment)
	// fmt.Println(matches)
	// check if is filename
	// todo: match embedCommand against different regexes

	return nil, fmt.Errorf("%q is not a valid command", embedComment)
}

func ExtractCodeBlocks(source string) []CodeBlock {
	var blocks []CodeBlock

	newline := internal.DetectNewline([]byte(source))

	matches := internal.GetMatches(blockRe, string(source))
	for _, match := range matches {
		var block CodeBlock
		if b, ok := match["block"]; ok {
			block.start = b.Start
			block.end = b.End
			block.Code = b.Text
			block.StartLine = internal.GetLineNumber(source, b.Start, newline)
			block.EndLine = internal.GetLineNumber(source, b.End, newline)
		}

		// block.start = match.Start
		// block.end = match.End
		// block.Code = match.Text
		// block.StartLine = LineNumber(source, match.Start, newline)
		// block.EndLine = LineNumber(source, match.End, newline)

		if comment, ok := match["embedComment"]; ok {
			block.embedComment = comment.Text
		}
		if language, ok := match["language"]; ok {
			block.language = language.Text
		}
		blocks = append(blocks, block)
	}
	return blocks

	// matches := blockRe.FindAllStringSubmatchIndex(string(source), -1)
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

	// 	// fmt.Printf("%v\n", blockRe.SubexpNames())
	// 	// fmt.Printf("%v\n", pos)

	// 	for _, name := range blockRe.SubexpNames() {
	// 		i := blockRe.SubexpIndex(name)
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

func commentPrefixRegex(comment string) *regexp.Regexp {
	return regexp.MustCompile(`^\s*` + comment + `([\s\S]*?)\r?\n`)
}

func FirstComment(source string, typ CommentType) (string, error) {
	var re *regexp.Regexp
	switch typ {
	case COMMENT_NONE:
		return "", nil
	case COMMENT_DOUBLE_SLASH:
		re = commentPrefixRegex("//")
		// return "", nil
	case COMMENT_XML:
		re = regexp.MustCompile(`<!--\s*?(\S*?)\s*?-->`)
		// const match = line.match(/<!--\s*?(\S*?)\s*?-->/);
		// return "", nil
	case COMMENT_HASH:
		re = commentPrefixRegex("#")
		// return "", nil
	case COMMENT_SINGLE_QUOTE:
		re = commentPrefixRegex("'")
		// return "", nil
	case COMMENT_DOUBLE_PERCENT:
		re = commentPrefixRegex("%%")
		// return "", nil
	case COMMENT_DOUBLE_HYPHENS:
		re = commentPrefixRegex("--")
		// return "", nil
	}
	// fmt.Printf("comment: %s\n", *typ)
	// fmt.Println(source)
	// fmt.Println(re.String())
	match := re.FindStringSubmatch(source)
	// fmt.Println(match)
	if len(match) > 0 {
		return match[1], nil
	}
	// i := blockRe.SubexpIndex(name)
	// if (!match) {
	//   return null;
	// }

	// return match[1];

	return "", nil
}
