package embedme

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/romnn/embedme/internal"
)

// var (
// // Detects info string following the beginning of a code block
// // infoStringRe = regexp.MustCompile("```(.*)")

// // Detects line ending to use based on whether CRLF is used in the source
// // const lineEnding = detectLineEnding(sourceText);
// // leRe = regexp.MustCompile("/\r\n/")
// )

var (
	backtickRegex = regexp.MustCompile("^```")
)

func getReplacement(path string, options *Options, block *CodeBlock, newline string) (string, error) {
	if false {
		color.Blue(`"Ignore next" comment detected, skipping ...`)
		return block.Code, nil
	}

	// let commentedFilename: string | null;
	// if (commentEmbedOverrideFilepath) {
	//   commentedFilename = commentEmbedOverrideFilepath;
	// } else {
	//   if (!codeExtension) {
	//     log({ returnSnippet: substr }, chalk => chalk.blue(
	// if block.Language == "" {
	// 	color.Blue("No code extension detected, skipping ...")
	// 	return block.Code, nil
	// }
	language, err := block.Language()
	if err != nil {
		color.Blue(err.Error())
		return block.Code, nil
	}

	if false {
		color.Blue(
			"Code block is empty and no preceding embedme comment, skipping...",
		)
		return block.Code, nil
	}

	// var supported []string
	// for _, languages := range languageComments {
	// 	for _, language := range languages {
	// 		supported = append(supported, string(language))
	// 	}
	// }
	// if comment, ok :=
	// language := Language(block.Language)
	// comment, ok := CommentForLanguage[language]

	_, err = block.CommentType()
	if err != nil {
		color.Yellow(err.Error())
		return block.Code, nil
	}

	// if !ok {
	// 	color.Yellow(
	// 		"Unsupported file extension %q, supported extensions are %s, skipping...",
	// 		block.Language,
	// 		strings.Join(SupportedLanguages, ", "),
	// 	)
	// 	return block.Code, nil
	// }

	// comment, ok := CommentForLanguage[language]
	// if !ok {
	// 	color.Red(
	// 		"File extension %q supported, ",
	// 		"but comment family could not be determined.",
	// 	)
	// }
	// color.Yellow("%s", comment)

	// var err error
	// embedCommand := block.embedCommand
	// if embedCommand == "" {
	// 	embedCommand, err = EmbedCommand(block.Code, comment)
	// 	color.Yellow("%s", firstComment)
	// 	color.Yellow("%v", err)
	// }

	// if embedCommand == "" {
	// 	color.White(
	// 		"No command detected in first line for block with extension %q",
	// 		block.Language,
	// 	)
	// }
	command, err := block.EmbedCommand(options)
	if err != nil {
		color.Red(err.Error())
		return block.Code, nil
	}
	if command == nil {
		color.White(
			"No command detected in first line for block with extension %q",
			language,
		)
		return block.Code, nil
	}

	// get the output of the command
	lines, err := command.Output()
	if err != nil {
		return block.Code, err
	}
	// todo: use the newline characters from the readme!!
	output := strings.Join(lines, newline)

	// fmt.Println(backtickRegex.FindAllString(output, 1))
	if len(backtickRegex.FindAllIndex([]byte(output), 1)) > 0 {
		// invalid
		return block.Code, fmt.Errorf(
			"refusing to embed:\n%s\n as it contains a code block",
			strings.Join(internal.PreviewLines(lines, 3), newline),
		)
	}
	// todo: diff here now
	if output == block.Code {
		color.White("No changes required, already up to date")
		return block.Code, nil
	}

	// const chalkColour = options.verify ? 'yellow' : 'green';
	if options.Verify {
		color.Yellow("Embedded %d lines", len(lines))
	} else {
		color.Green("Embedded %d lines", len(lines))
	}
	// todo: type switch on the command for better logging
	// const [, filename, , lineNumbering, startLine, endLine] = matches;

	// cli command: ^\s*\$\s*([\s\S]+?)$
	// file path with lines: ^\s*(\S+?)((#L(\d+)-L(\d+))|$)
	// const matches = commentedFilename.match(/\s?(\S+?)((#L(\d+)-L(\d+))|$)/m);

	return "", nil
}

// Embed ...
func Embed(source []byte, path string, options *Options) (string, error) {
	color.Magenta(" Analysing %s ...", path)

	// le := "\n"
	// if leRe.Match(source) {
	// 	le = "\r\n"
	// }
	// if false {
	// 	color.Magenta(" le %s ...", le)
	// }

	// const docPartials = [];
	// let previousEnd = 0;
	// let result: RegExpExecArray | null;
	// let replacementError = false;
	newline := internal.DetectNewline(source)

	// fmt.Print(string(source))
	blocks := ExtractCodeBlocks(string(source))

	// matches := blockRe.FindAllStringSubmatch(string(source), -1)
	// for _, match := range matches {
	// 	block := match[0]
	// 	group := make(map[string]string)
	// 	for i, name := range blockRe.SubexpNames() {
	// 		if i != 0 && name != "" {
	// 			group[name] = match[i]
	// 		}
	// 	}
	limit := 2
	for _, block := range blocks {

		// color.Magenta("\n%+v\n", block)

		// color.Magenta("\n%v\n", block)
		// color.Magenta("\n%+v\n", group)

		// if (options.dryRun || options.stdout || options.verify) {
		// text.substring(0, index).split(lineEnding).length;

		//     return getLineNumber(sourceText.substring(0, result.index), result.index, lineEnding);
		//   }
		//   const startingLineNumber = docPartials.join('').split(lineEnding).length - 1;
		//   return (
		//     startingLineNumber + getLineNumber(sourceText.substring(previousEnd, result.index), result.index, lineEnding)
		//   );

		// const commentInsertion = start.match(/<!--\s*?embedme[ ]+?(\S+?)\s*?-->/);

		// /<!--\s*?embedme[ -]ignore-next\s*?-->/g.test(start),

		text, err := getReplacement(
			path,
			options,
			// log,
			&block,
			newline,
			// leadingSpaces,
			// lineEnding,
			// infoString,
			// codeExtension as SupportedFileType,
			// firstLine || '',
			// startLineNumber,
			// commentInsertion ? commentInsertion[1] : undefined,
		)
		if err != nil {
			panic(err)
		}
		if false {
			color.Magenta("\n%s\n", text)
		}
		limit--
		if limit <= 0 {
			break
		}
	}
	// color.Magenta("\n%+v\n", block)

	//   const [codeFence, leadingSpaces] = result;
	//   const start = sourceText.substring(previousEnd, result.index);
	// const infoStringMatch = codeFence.match("/```(.*)/")
	// const infoString = infoStringMatch ? infoStringMatch[1] : '';
	// const codeExtension = infoString !== '' ? infoString.trim().split(/\s/)[0] : null;
	// }
	return "", nil
}
