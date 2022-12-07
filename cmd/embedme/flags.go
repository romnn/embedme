package main

import (
	"github.com/urfave/cli/v3"
)

const ENV_PREFIX = "EMBEDME"

var (
	verifyFlag = cli.BoolFlag{
		Name:    "verify",
		EnvVars: []string{ENV_PREFIX + "_VERIFY"},
		Usage:   "verify that running embedme would result in no changes. Useful for CI",
	}
	dryRunFlag = cli.BoolFlag{
		Name:    "dry-run",
		EnvVars: []string{ENV_PREFIX + "_DRY_RUN"},
		Usage:   "run embedme as usual, but don't write",
	}
	noColorFlag = cli.BoolFlag{
		Name:    "no-color",
		EnvVars: []string{ENV_PREFIX + "_NO_COLOR"},
		Usage:   "disable colored output",
	}
	forceColorFlag = cli.BoolFlag{
		Name:    "force-color",
		EnvVars: []string{ENV_PREFIX + "_FORCE_COLOR"},
		Usage:   "force colored output",
	}
	sourceBaseFlag = cli.StringFlag{
		Name:    "base",
		Aliases: []string{"source", "root"},
		EnvVars: []string{ENV_PREFIX + "_BASE_PATH"},
		Usage:   "source files directory prefix for shorter code block comments",
	}
	cwdFlag = cli.StringFlag{
		Name:    "directory",
		Aliases: []string{"C"},
		EnvVars: []string{ENV_PREFIX + "_DIRECTORY"},
		Usage:   "run embedme from this directory",
	}
	silentFlag = cli.BoolFlag{
		Name:    "silent",
		EnvVars: []string{ENV_PREFIX + "_SILENT"},
		Usage:   "no console output",
	}
	globFlag = cli.BoolFlag{
		Name:    "glob",
		EnvVars: []string{ENV_PREFIX + "_GLOB"},
		Usage:   "treat arguments as patterns and glob from current directory",
	}
	stdoutFlag = cli.BoolFlag{
		Name:    "stdout",
		EnvVars: []string{ENV_PREFIX + "_STDOUT"},
		Usage:   "output resulting file to stdout without writing",
	}
	outputFlag = cli.StringFlag{
		Name:    "output",
		Aliases: []string{"dest", "destination"},
		EnvVars: []string{ENV_PREFIX + "_OUTPUT"},
		Usage:   "output file to write the result into",
	}
	stripEmbedCommentFlag = cli.BoolFlag{
		Name:    "strip-embed-comment",
		EnvVars: []string{ENV_PREFIX + "_STRIP_EMBED_COMMENT"},
		Usage:   "remove the comment from the code block (only works with --stdout )",
	}
)
