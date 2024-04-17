package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/romnn/embedme/internal"
	"github.com/romnn/embedme/pkg"
	"github.com/sabhiram/go-gitignore"
	"github.com/urfave/cli/v3"
)

// Version is set during build
var Version = ""

// Rev is set during build
var Rev = ""

var (
	Magenta = color.New(color.FgMagenta).FprintfFunc()
	Info    = color.New(color.FgBlue).FprintfFunc()
	Warning = color.New(color.FgYellow).FprintfFunc()
	Error   = color.New(color.FgRed).FprintfFunc()
	Log     = color.New(color.FgWhite).FprintfFunc()
)

func versionString() string {
	return fmt.Sprintf("%s (%s)", Version, Rev)
}

func run(ctx context.Context, cmd *cli.Command) error {
	start := time.Now()

	stdout := cmd.Bool(stdoutFlag.Name)
	silent := cmd.Bool(silentFlag.Name)
	useColor := cmd.Bool(colorFlag.Name)
	forceColor := cmd.Bool(forceColorFlag.Name)

	if silent {
		log.SetOutput(io.Discard)
	} else if stdout {
		// the result will be written to stdout,
		// so we redirect the logs to stderr
		log.SetOutput(os.Stderr)
	}

	if !useColor {
		color.NoColor = true
	}
	if forceColor {
		color.NoColor = false
	}

	Magenta(log.Writer(), "embedme v%s\n", versionString())

	realCwd, err := os.Getwd()
	if err != nil {
		return err
	}

	output := cmd.String(outputFlag.Name)

	glob := cmd.Bool(globFlag.Name)
	cwd := cmd.String(cwdFlag.Name)
	base := cmd.String(sourceBaseFlag.Name)
	if cwd == "" {
		cwd = realCwd
	}
	if base == "" {
		base = cwd
	}

	options := embedme.Options{
		StripEmbedComment: cmd.Bool(stripEmbedCommentFlag.Name),
		Stdout:            stdout,
		Verify:            cmd.Bool(verifyFlag.Name),
		DryRun:            cmd.Bool(dryRunFlag.Name),
		Cwd:               cwd,
		Base:              base,
	}

	allFlags := make(map[string]bool)
	for _, flag := range cmd.Flags {
		for _, name := range flag.Names() {
			allFlags[name] = true
		}
	}

	sources := make(sourceMap)
	for _, arg := range cmd.Args().Slice() {
		flag := strings.TrimLeft(strings.TrimSpace(arg), "-")
		if _, ok := allFlags[flag]; ok {
			// is a flag, skip
			continue
		}
		if glob {
			matches, err := fs.Glob(os.DirFS(cwd), arg)
			if err != nil {
				return fmt.Errorf("failed to glob pattern %q in %s: %v", arg, cwd, err)
			}
			sources.Add(matches...)
		} else {
			if !filepath.IsAbs(arg) {
				arg = filepath.Join(cwd, arg)
			}
			sources.Add(arg)
		}
	}

	if len(sources) > 1 && (options.Stdout || output != "") {
		Warning(log.Writer(), "more than one file matched: results will be concatenated")
	}
	if len(sources) == 0 {
		Warning(log.Writer(), "no files matched your input")
		return nil
	}

	if options.StripEmbedComment && !options.Stdout {
		Error(log.Writer(), `Invalid use of --strip-embed-comment.
If you use the --strip-embed-comment flag, you must use the --stdout flag
and redirect the result to your destination file, otherwise your source
file(s) will be overwritten and the comment source is lost.`)
		os.Exit(1)
	}

	if options.Verify {
		Info(log.Writer(), "Verifying...\n")
	} else if options.DryRun {
		Info(log.Writer(), "Doing a dry run...\n")
	} else if options.Stdout {
		Info(log.Writer(), "Writing to stdout...\n")
	} else {
		Info(log.Writer(), "Embedding...\n")
	}

	for _, ignoreFile := range []string{".embedmeignore", ".gitignore"} {
		ignorePath := filepath.Join(cwd, ignoreFile)
		if _, err := os.Stat(ignorePath); err != nil {
			continue
		}

		ignore, err := ignore.CompileIgnoreFile(ignorePath)
		if err != nil {
			continue
		}
		ignored := sources.Ignore(ignore)

		if ignored > 0 {
			Info(log.Writer(), "Skipped %d files ignored in %s\n", ignored, ignoreFile)
		}
	}

	if len(sources) == 0 {
		Warning(log.Writer(), "All matching files were ignored\n")
		return nil
	}

	for i, source := range sources.Valid() {
		relSource := source
		if rel, err := filepath.Rel(cwd, source); err == nil {
			relSource = rel
		}

		if i > 0 {
			Log(log.Writer(), "---")
		}
		log.SetPrefix("test")

		if err := internal.EnsureFile(source); err != nil {
			return fmt.Errorf("file %s does not exist: %v", relSource, err)
		}

		markdown, err := os.ReadFile(source)
		if err != nil {
			return fmt.Errorf("file %s could not be read: %v", relSource, err)
		}

		embedded, err := embedme.Embed(markdown, source, relSource, &options)
		if err != nil {
			return fmt.Errorf("failed to embed %s: %v", relSource, err)
		}

		diff := string(markdown) != embedded
		if options.Verify {
			if diff {
				return fmt.Errorf("Difference detected, exiting 1\n")
			}
		} else if options.Stdout {
			fmt.Print(embedded)
		} else if !options.DryRun {
			if diff {
				Magenta(log.Writer(), "Writing %s with embedded changes.\n", relSource)
				f, err := os.Open(source)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				// if _, err := f.Write([]byte(embedded)); err != nil {
				//   panic(err)
				// }
			} else {
				Magenta(log.Writer(), "No changes to write for %s\n", relSource)
			}
		}
	}

	Magenta(log.Writer(), "done in %v\n", time.Now().Sub(start))
	return nil
}

func main() {
	app := &cli.Command{
		Name:        "embedme",
		Description: "utility for embedding code snippets into markdown documents",
		Version:     versionString(),
		Flags: []cli.Flag{
			&verifyFlag,
			&dryRunFlag,
			&forceColorFlag,
			&colorFlag,
			&cwdFlag,
			&sourceBaseFlag,
			&globFlag,
			&silentFlag,
			&stdoutFlag,
			&outputFlag,
			&stripEmbedCommentFlag,
		},
		Action: run,
	}
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		Error(log.Writer(), "error: %v\n", err)
		os.Exit(1)
	}
}
