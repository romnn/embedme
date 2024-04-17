package main

import (
	"github.com/urfave/cli/v3"
)

// EnvPrefix ...
const EnvPrefix = "EMBEDME"

var (
	verifyFlag = cli.BoolFlag{
		Name:    "verify",
		Sources: cli.EnvVars(EnvPrefix + "_VERIFY"),
		Usage:   "verify that running embedme would result in no changes. Useful for CI",
	}
	dryRunFlag = cli.BoolFlag{
		Name:    "dry-run",
		Sources: cli.EnvVars(EnvPrefix + "_DRY_RUN"),
		Usage:   "run embedme as usual, but don't write",
	}
	colorFlag = cli.BoolFlag{
		Name:    "color",
		Sources: cli.EnvVars(EnvPrefix + "_COLOR"),
		Value:   true,
		Usage:   "disable colored output",
	}
	forceColorFlag = cli.BoolFlag{
		Name:    "force-color",
		Sources: cli.EnvVars(EnvPrefix + "_FORCE_COLOR"),
		Usage:   "force colored output",
	}
	sourceBaseFlag = cli.StringFlag{
		Name:    "base",
		Aliases: []string{"source", "root"},
		Sources: cli.EnvVars(EnvPrefix + "_BASE_PATH"),
		Usage:   "source files directory prefix for shorter code block comments",
	}
	cwdFlag = cli.StringFlag{
		Name:    "directory",
		Aliases: []string{"C"},
		Sources: cli.EnvVars(EnvPrefix + "_DIRECTORY"),
		Usage:   "run embedme from this directory",
	}
	silentFlag = cli.BoolFlag{
		Name:    "silent",
		Sources: cli.EnvVars(EnvPrefix + "_SILENT"),
		Usage:   "disable console output",
	}
	globFlag = cli.BoolFlag{
		Name:    "glob",
		Sources: cli.EnvVars(EnvPrefix + "_GLOB"),
		Usage:   "treat arguments as patterns and glob from current directory",
	}
	stdoutFlag = cli.BoolFlag{
		Name:    "stdout",
		Sources: cli.EnvVars(EnvPrefix + "_STDOUT"),
		Usage:   "output resulting file to stdout without writing",
	}
	outputFlag = cli.StringFlag{
		Name:    "output",
		Aliases: []string{"dest", "destination"},
		Sources: cli.EnvVars(EnvPrefix + "_OUTPUT"),
		Usage:   "output file to write the result into",
	}
	stripEmbedCommentFlag = cli.BoolFlag{
		Name:    "strip-embed-comment",
		Sources: cli.EnvVars(EnvPrefix + "_STRIP_EMBED_COMMENT"),
		Usage:   "remove the comment from the code block (only works with --stdout)",
	}
)
