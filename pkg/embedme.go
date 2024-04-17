package embedme

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/romnn/embedme/internal"
)

var (
	backtickRegex = regexp.MustCompile("^```")
	// Magenta ...
	Magenta = color.New(color.FgMagenta).FprintfFunc()
	// Info ...
	Info = color.New(color.FgBlue).FprintfFunc()
	// Warning ...
	Warning = color.New(color.FgYellow).FprintfFunc()
	// Error ...
	Error = color.New(color.FgRed).FprintfFunc()
	// Log ...
	Log = color.New(color.FgWhite).FprintfFunc()
)

func embedBlock(
	path string,
	relPath string,
	options *Options,
	block *CodeBlock,
	newline string,
	startLine int,
) (string, error) {
	// relPath, err := filepath.Rel(options.Cwd, path)
	// if err != nil {
	// panic(err)
	//
	endLine := len(internal.Lines(block.Code)) + startLine
	logPrefix := fmt.Sprintf("  %s#L%d-L%d", relPath, startLine, endLine)
	log.SetPrefix(logPrefix)

	if block.Ignore {
		// log
		// Magenta(log.Writer(), "embedme v%s\n", versionString())
		Info(log.Writer(), "Ignore next comment detected, skipping ...\n")
		// color.Blue(`"Ignore next" comment detected, skipping ...`)
		return block.Code, nil
	}
	Info(log.Writer(), logPrefix+"\n")

	// let commentedFilename: string | null;
	// if (commentEmbedOverrideFilepath) {
	//   commentedFilename = commentEmbedOverrideFilepath;
	// } else {
	//   if (!codeExtension) {
	//     log({ returnSnippet: substr }, chalk => chalk.blue(

	if block.Language == "" {
		Info(log.Writer(), "No code extension detected, skipping ...\n")
		return block.Code, nil
	}

	// language, err := block.Language()
	// languagerr := block.Language()
	// if err != nil {
	// 	color.Blue(err.Error())
	// 	return block.Code, nil
	// }

	// if false {
	// 	color.Blue(
	// 		"Code block is empty and no preceding embedme comment, skipping...",
	// 	)
	// 	return block.Code, nil
	// }

	// var supported []string
	// for _, languages := range languageComments {
	// 	for _, language := range languages {
	// 		supported = append(supported, string(language))
	// 	}
	// }
	// if comment, ok :=
	// language := Language(block.Language)
	// comment, ok := CommentForLanguage[language]

	if _, err := block.CommentType(); err != nil {
		Warning(log.Writer(), err.Error()+"\n")
		// color.Yellow(err.Error())
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
	comment, command, err := block.EmbedCommand(options)
	if err != nil {
		Error(log.Writer(), err.Error()+"\n")
		// color.Red(err.Error())
		return block.Code, nil
	}
	if command == nil {
		color.White(
			"No command detected in first line for block with extension %q",
			block.Language,
		)
		return block.Code, nil
	}

	// get the output of the command
	lines, err := command.Output()
	if err != nil {
		return block.Code, err
	}

	output := strings.Join(lines, newline)

	// todo: diff here now
	if output == block.Code {
		color.White("No changes required, already up to date")
		return block.Code, nil
	}

	// fmt.Println(backtickRegex.FindAllString(output, 1))
	if len(backtickRegex.FindAllIndex([]byte(output), 1)) > 0 {
		// invalid
		return block.Code, fmt.Errorf(
			"refusing to embed:\n%s\n as it contains a code block",
			strings.Join(internal.PreviewLines(lines, 3), newline),
		)
	}

	var replacement string
	replacement += "```"
	replacement += string(block.Language)
	replacement += newline
	if !(options.StripEmbedComment || block.EmbedComment != "") {
		replacement += block.Comment()
		replacement += comment // .Command
		replacement += newline
		replacement += newline
	}
	replacement += output
	replacement += newline
	replacement += "```"

	// indent
	replacementLines := internal.Lines(replacement)
	for i, line := range replacementLines {
		replacementLines[i] = block.Indent + line
	}
	replacement = strings.Join(replacementLines, newline)

	// fmt.Println(replacement)

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

	return replacement, nil
}

// Embed embeds a document
func Embed(
	markdown []byte,
	path string,
	relPath string,
	options *Options,
) (string, error) {
	color.Magenta("Analysing %s ...", relPath)

	var partials []string
	previousEnd := 0
	newline := internal.DetectNewline(markdown)

	blocks := ExtractCodeBlocks(string(markdown))

	limit := 2
	for _, block := range blocks {

		// color.Magenta("\n%+v\n", block)

		// color.Magenta("\n%v\n", block)
		// color.Magenta("\n%+v\n", group)

		startLine := 0
		if options.DryRun || options.Stdout || options.Verify {
			startLine = block.StartLine
		} else {
			startLine = len(internal.Lines(strings.Join(partials, newline))) - 1
		}
		// text.substring(0, index).split(lineEnding).length;

		//     return getLineNumber(sourceText.substring(0, result.index), result.index, lineEnding);
		//   }
		//   const startingLineNumber = docPartials.join('').split(lineEnding).length - 1;
		//   return (
		//     startingLineNumber + getLineNumber(sourceText.substring(previousEnd, result.index), result.index, lineEnding)
		//   );

		// const commentInsertion = start.match(/<!--\s*?embedme[ ]+?(\S+?)\s*?-->/);

		// /<!--\s*?embedme[ -]ignore-next\s*?-->/g.test(start),

		embedded, err := embedBlock(
			path,
			relPath,
			options,
			// log,
			&block,
			newline,
			startLine,
			// leadingSpaces,
			// lineEnding,
			// infoString,
			// codeExtension as SupportedFileType,
			// firstLine || '',
			// startLineNumber,
			// commentInsertion ? commentInsertion[1] : undefined,
		)
		if err != nil {
			return "", err
		}
		partials = append(partials, string(markdown)[previousEnd:block.Start])
		partials = append(partials, embedded)
		previousEnd = block.End
		if false {
			color.Magenta("\n%s\n", embedded)
		}
		limit--
		if limit <= 0 {
			break
		}
	}

	final := strings.Join(partials, newline)
	if false {
		fmt.Println(final)
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
