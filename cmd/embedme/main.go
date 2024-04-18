package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	embedme "github.com/romnn/embedme/pkg"
	"github.com/urfave/cli/v3"
)

// Version is set during build
var Version = ""

// Rev is set during build
var Rev = ""

func versionString() string {
	return fmt.Sprintf("%s (%s)", Version, Rev)
}

// config contains all embedme CLI options
type config struct {
	Stdout     bool
	Silent     bool
	UseColor   bool
	ForceColor bool
	Verify     bool
	DryRun     bool
	Output     string
	Glob       bool
	WorkingDir string
	Base       string
}

func parseConfig(cmd *cli.Command) (config, error) {
	config := config{
		Stdout:     cmd.Bool(stdoutFlag.Name),
		Silent:     cmd.Bool(silentFlag.Name),
		UseColor:   cmd.Bool(colorFlag.Name),
		ForceColor: cmd.Bool(forceColorFlag.Name),
		Verify:     cmd.Bool(verifyFlag.Name),
		DryRun:     cmd.Bool(dryRunFlag.Name),
		Output:     cmd.String(outputFlag.Name),
		Glob:       cmd.Bool(globFlag.Name),
		WorkingDir: cmd.String(cwdFlag.Name),
		Base:       cmd.String(sourceBaseFlag.Name),
	}
	realWorkingDir, err := os.Getwd()
	if err != nil {
		return config, err
	}
	if config.WorkingDir == "" {
		config.WorkingDir = realWorkingDir
	}
	if config.Base == "" {
		config.Base = config.WorkingDir
	}
	return config, nil
}

func sourcePatterns(cmd *cli.Command) []string {
	allFlags := make(map[string]bool)
	for _, flag := range cmd.Flags {
		for _, name := range flag.Names() {
			allFlags[name] = true
		}
	}
	patterns := []string{}
	for _, arg := range cmd.Args().Slice() {
		flag := strings.TrimLeft(strings.TrimSpace(arg), "-")
		if _, ok := allFlags[flag]; !ok {
			// is not flag
			patterns = append(patterns, arg)
		}
	}
	return patterns
}

func logOperation(options *embedme.Options) {
	if options.Verify {
		embedme.Info(log.Writer(), "Verifying...\n")
	} else if options.DryRun {
		embedme.Info(log.Writer(), "Doing a dry run...\n")
	} else if options.Stdout {
		embedme.Info(log.Writer(), "Writing to stdout...\n")
	} else {
		embedme.Info(log.Writer(), "Embedding...\n")
	}
}

func configureOutput(config config) {
	if config.Silent {
		log.SetOutput(io.Discard)
	} else if config.Stdout {
		// the result will be written to stdout,
		// so we redirect the logs to stderr
		log.SetOutput(os.Stderr)
	}

	if !config.UseColor {
		color.NoColor = true
	}
	if config.ForceColor {
		color.NoColor = false
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	start := time.Now()

	config, err := parseConfig(cmd)
	if err != nil {
		return err
	}

	configureOutput(config)

	embedme.Magenta(log.Writer(), "embedme v%s\n", versionString())

	options := embedme.Options{
		StripEmbedComment: cmd.Bool(stripEmbedCommentFlag.Name),
		Stdout:            config.Stdout,
		Verify:            config.Verify,
		DryRun:            config.DryRun,
		WorkingDir:        config.WorkingDir,
		Base:              config.Base,
	}

	embedder, err := embedme.NewEmbedder(options)

	ignoreFiles, err := embedme.GlobFiles(
		embedder.FS,
		options.WorkingDir,
		".embedmeignore", ".gitignore",
	)
	if err != nil {
		return fmt.Errorf("failed to find ignore files: %v", err)
	}
	finder := embedme.SourceFinder{
		WorkingDir:  options.WorkingDir,
		Glob:        config.Glob,
		IgnoreFiles: ignoreFiles,
	}

	srcPatterns := sourcePatterns(cmd)
	sources, err := finder.FindSources(embedder.FS, srcPatterns...)
	if err != nil {
		return err
	}
	validSources := sources.Valid()

	if len(sources) == 0 {
		embedme.Warning(log.Writer(), "no files matched your input")
		return nil
	}

	if len(sources) > 1 && (options.Stdout || config.Output != "") {
		embedme.Warning(log.Writer(), "more than one file matched: results will be concatenated")
	}

	if options.StripEmbedComment && !options.Stdout {
		return fmt.Errorf(`invalid use of --strip-embed-comment.
If you use the --strip-embed-comment flag, you must use the --stdout flag
and redirect the result to your destination file, otherwise your source
file(s) will be overwritten and the comment source is lost`)
	}

	logOperation(&options)

	if len(validSources) == 0 {
		embedme.Warning(log.Writer(), "All matching files were ignored\n")
		return nil
	}

	for i, source := range validSources {
		if err := embedder.ProcessSource(i, source); err != nil {
			return err
		}
	}

	embedme.Magenta(log.Writer(), "done in %v\n", time.Now().Sub(start))
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
		embedme.Error(log.Writer(), "error: %v\n", err)
		os.Exit(1)
	}
}
