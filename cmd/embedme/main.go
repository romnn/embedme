package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/romnn/embedme/pkg"
	"github.com/romnn/embedme/internal"
	"github.com/sabhiram/go-gitignore"
	"github.com/urfave/cli/v2"
)

// Version is set during build
var Version = ""

// Rev is set during build
var Rev = ""

func versionString() string {
	return fmt.Sprintf("%s (%s)", Version, Rev)
}

func run(cliCtx *cli.Context) error {
	// todo: cli arg
	color.NoColor = false

	realCwd, err := os.Getwd()
	if err != nil {
		return err
	}

	glob := cliCtx.Bool(globFlag.Name)
	cwd := cliCtx.String(cwdFlag.Name)
	base := cliCtx.String(sourceBaseFlag.Name)
	if cwd == "" {
		cwd = realCwd
	}
	if base == "" {
		base = cwd
	}

	options := embedme.Options{
		StripEmbedComment: cliCtx.Bool(stripEmbedCommentFlag.Name),
		Stdout:            cliCtx.Bool(stdoutFlag.Name),
		Verify:            cliCtx.Bool(verifyFlag.Name),
		DryRun:            cliCtx.Bool(dryRunFlag.Name),
		Cwd:               cwd,
		Base:              base,
	}

	sources := make(sourceMap)
	for _, arg := range cliCtx.Args().Slice() {
		if glob {
			fmt.Printf("globbing for %s\n", arg)
			matches, err := fs.Glob(os.DirFS(cwd), arg)
			if err != nil {
				color.Red("failed to glob pattern %q in %s: %v", arg, cwd, err)
				return nil
			}
			sources.Add(matches...)
		} else {
			sources.Add(arg)
		}
	}

	if len(sources) > 1 {
		color.Yellow("more than one file matched your input, results will be concatenated in stdout")
	} else if len(sources) == 0 {
		color.Yellow("no files matched your input")
		return nil
	}

	if options.StripEmbedComment && !options.Stdout {
		color.Red("If you use the --strip-embed-comment flag, you must use the --stdout flag and redirect the result to your destination file, otherwise your source file(s) will be rewritten and the comment source is lost.")
		return nil
	}

	if options.Verify {
		color.Blue("Verifying...")
	} else if options.DryRun {
		color.Blue("Doing a dry run...")
	} else if options.Stdout {
		color.Blue("Outputting to stdout...")
	} else {
		color.Blue("Embedding...")
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
			color.Blue("Skipped %d files ignored in %s", ignored, ignoreFile)
		}
	}

	if len(sources) == 0 {
		color.Yellow("All matching files were ignored")
		return nil
	}

	for i, source := range sources.Valid() {
		if i > 0 {
			color.White("---")
		}

		// if _, err := os.Stat(source); err != nil {
		// 	color.Red("File %s does not exist", source)
		// 	return nil
		// }
		if err := internal.EnsureFile(source); err != nil {
			return fmt.Errorf("file %s does not exist: %v", source, err)
		}

		markdown, err := os.ReadFile(source)
		if err != nil {
			color.Red("File %s could not be read: %v", source, err)
			return nil
		}

		embedded, err := embedme.Embed(markdown, source, &options)
		if err != nil {
			color.Red("failed to embed %s: %v", source, err)
			return nil
		}

		diff := string(markdown) != embedded
		if options.Verify {
			if diff {
				color.Red("Diff detected, exiting 1")
				return errors.New("diff")
			}
		} else if options.Stdout {
			fmt.Print(embedded)
		} else if !options.DryRun {
			if diff {
				color.Magenta("Writing %s with embedded changes.", source)
				f, err := os.Open(source)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				// if _, err := f.Write([]byte(embedded)); err != nil {
				//   panic(err)
				// }
				// writeFileSync(source, outText)
			} else {
				color.Magenta("No changes to write for %s", source)
			}
		}
	}

	return nil
}

func main() {
	// fmt.Printf("embedme v%s\n", versionString())

	app := &cli.App{
		Name:    "embedme",
		Usage:   "utility for embedding code snippets into markdown documents",
		Version: versionString(),
		Flags: []cli.Flag{
			&verifyFlag,
			&dryRunFlag,
			&cwdFlag,
			&sourceBaseFlag,
			&globFlag,
			&silentFlag,
			&stdoutFlag,
			&stripEmbedCommentFlag,
		},
		Action: run,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
