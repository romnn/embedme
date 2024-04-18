package embedme

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/romnn/embedme/internal"
	"github.com/romnn/embedme/pkg/fs"
	"github.com/spf13/afero"
)

var (
	backtickRegex = regexp.MustCompile("^```")
	// Magenta ...
	Magenta = color.New(color.FgMagenta).FprintfFunc()
	// Info ...
	Success = color.New(color.FgGreen).FprintfFunc()
	// Info ...
	Info = color.New(color.FgBlue).FprintfFunc()
	// Warning ...
	Warning = color.New(color.FgYellow).FprintfFunc()
	// Error ...
	Error = color.New(color.FgRed).FprintfFunc()
	// Log ...
	Log = color.New(color.FgWhite).FprintfFunc()
)

// Embedder ...
type Embedder struct {
	Options Options
	// FS      fs.FileSystem
	FS afero.Fs
}

func NewEmbedder(options Options) (Embedder, error) {
	return Embedder{
		Options: options,
		FS:      afero.OsFs{},
	}, nil
}

func (e *Embedder) ProcessSource(i int, absSource string) error {
	if !filepath.IsAbs(absSource) {
		absSource = filepath.Join(e.Options.WorkingDir, absSource)
	}

	if !filepath.IsAbs(absSource) {
		log.Panicf("expected absolute source path, but got %s", absSource)
	}

	relSource, err := filepath.Rel(e.Options.WorkingDir, absSource)
	if err != nil {
		log.Panicf(
			"source path %s is outside of working dir %s",
			absSource, e.Options.WorkingDir,
		)
	}

	if i > 0 {
		Log(log.Writer(), "---")
	}
	log.SetPrefix("test")

	if err := fs.EnsureFile(e.FS, absSource); err != nil {
		return fmt.Errorf("file %s does not exist: %v", relSource, err)
	}

	markdown, err := afero.ReadFile(e.FS, absSource)
	if err != nil {
		return fmt.Errorf("file %s could not be read: %v", relSource, err)
	}

	embedded, err := e.Embed(markdown, absSource, relSource)
	if err != nil {
		return fmt.Errorf("failed to embed %s: %v", relSource, err)
	}

	diff := string(markdown) != embedded
	if e.Options.Verify {
		if diff {
			return fmt.Errorf("Difference detected, exiting 1\n")
		}
	} else if e.Options.Stdout {
		fmt.Print(embedded)
	} else if !e.Options.DryRun {
		if diff {
			Magenta(
				log.Writer(),
				"Writing %s with embedded changes.\n", relSource,
			)
			file, err := e.FS.OpenFile(
				absSource,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to open %s: %v", relSource, err)
			}
			defer file.Close()
			if _, err := file.Write([]byte(embedded)); err != nil {
				return fmt.Errorf("failed to write %s: %v", relSource, err)
			}
		} else {
			Success(log.Writer(), "No changes to write for %s\n", relSource)
		}
	}
	return nil
}

func (e *Embedder) embedBlock(
	absPath string,
	relPath string,
	block *CodeBlock,
	newline string,
	// startLine int,
) (string, error) {
	// relPath, err := filepath.Rel(options.Cwd, path)
	// if err != nil {
	// panic(err)
	// if startLine != block.StartLine {
	// 	log.Panicf(
	// 		"different start lines %d != %d",
	// 		startLine, block.StartLine,
	// 	)
	// }
	if !filepath.IsAbs(absPath) {
		log.Panicf("expected absolute code path, but got %s", absPath)
	}

	startLine := block.StartLine
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

	log.Printf("block code: %s\n", block.Code)

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
	commandComment, command, err := block.EmbedCommand(e.FS, &e.Options)
	if err != nil {
		Error(log.Writer(), err.Error()+"\n")
		// color.Red(err.Error())
		return block.Code, nil
	}
	if command == nil {
		color.White(
			"No command detected in first line for block with extension %q (comment %s)",
			block.Language,
			commandComment,
		)
		return block.Code, nil
	}

	// get the output of the command
	lines, err := command.Output()
	if err != nil {
		return block.Code, err
	}

	// log.Printf("command lines: %+v\n", lines)
	output := strings.Join(lines, newline)
	// output = strings.TrimRight(output, " \n\r\t")
	output = strings.TrimRightFunc(output, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	output += newline
	// log.Printf("command: %+v\n", output)

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
	// replacement += "```"
	// replacement += string(block.Language)
	// replacement += newline
	if !(e.Options.StripEmbedComment || block.EmbedComment != "") {
		// comment for langugage
		replacement += block.Comment() + " "
		// embed command
		replacement += strings.TrimSpace(commandComment)
		replacement += newline
		replacement += newline
	}
	replacement += output
	// replacement += newline
	// replacement += "```"

	// replacement = strings.TrimSpace(replacement)

	// indent
	replacementLines := internal.Lines(replacement)
	for i, line := range replacementLines {
		replacementLines[i] = block.Indent + line
	}
	replacement = strings.Join(replacementLines, newline)

	// fmt.Println(replacement)

	if e.Options.Verify {
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
func (e *Embedder) Embed(
	markdown []byte,
	absPath string,
	relPath string,
) (string, error) {
	color.Magenta("Analysing %s ...", relPath)

	var partials []string
	previousEnd := 0
	newline := internal.DetectNewline(markdown)

	blocks := ExtractCodeBlocks(string(markdown))

	// log.Printf("blocks: %v\n", blocks)

	// limit := 2
	for bi, block := range blocks {

		log.Printf("\n%+v\n", block)

		// color.Magenta("\n%v\n", block)
		// color.Magenta("\n%+v\n", group)

		// add the partial here
		between := string(markdown)[previousEnd:block.Start]
		// between = strings.TrimRight(between, ")
		log.Printf("PREV:\n%q\n", between)
		partials = append(partials, between)

		newStartLine := 0
		if e.Options.DryRun || e.Options.Stdout || e.Options.Verify {
			newStartLine = block.StartLine
		} else {
			currentLines := internal.Lines(strings.Join(partials, newline))
			newStartLine = len(currentLines) - 1
			log.Printf("block %d starts in line %d\n", bi, newStartLine)
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

		embedded, err := e.embedBlock(
			absPath,
			relPath,
			// log,
			&block,
			newline,
			// startLine,
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

		partials = append(partials, embedded)
		previousEnd = block.End

		log.Printf("EMBEDDED:\n%q\n", embedded)
		if false {
			color.Magenta("\n%s\n", embedded)
		}
		// limit--
		// if limit <= 0 {
		// 	break
		// }
	}

	// add the final partial here
	between := string(markdown)[previousEnd:]
	partials = append(partials, between)

	// final := strings.Join(partials, newline)
	final := strings.Join(partials, "")
	// if true {
	// fmt.Printf("==== final\n%s\n====\n", final)
	// }
	// color.Magenta("\n%+v\n", block)

	//   const [codeFence, leadingSpaces] = result;
	//   const start = sourceText.substring(previousEnd, result.index);
	// const infoStringMatch = codeFence.match("/```(.*)/")
	// const infoString = infoStringMatch ? infoStringMatch[1] : '';
	// const codeExtension = infoString !== '' ? infoString.trim().split(/\s/)[0] : null;
	// }
	return final, nil
}
