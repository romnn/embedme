package main

import (
	"github.com/urfave/cli/v2"
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
		EnvVars: []string{"EMBEDME_DRY_RUN"},
		Usage:   "run embedme as usual, but don't write",
	}
	sourceBaseFlag = cli.StringFlag{
		Name:    "base",
		Aliases: []string{"source", "root"},
		EnvVars: []string{"EMBEDME_BASE_PATH"},
		Usage:   "source files directory prefix for shorter code block comments",
	}
	cwdFlag = cli.StringFlag{
		Name:    "directory",
		Aliases: []string{"C"},
		EnvVars: []string{"EMBEDME_DIRECTORY"},
		Usage:   "run embedme from this directory",
	}
	silentFlag = cli.BoolFlag{
		Name:    "silent",
		EnvVars: []string{"EMBEDME_SILENT"},
		Usage:   "no console output",
	}
	globFlag = cli.BoolFlag{
		Name:    "glob",
		EnvVars: []string{"EMBEDME_GLOB"},
		Usage:   "treat arguments as patterns and glob from current directory",
	}
	stdoutFlag = cli.BoolFlag{
		Name:    "stdout",
		EnvVars: []string{"EMBEDME_STDOUT"},
		Usage:   "output resulting file to stdout without writing to original",
	}
	stripEmbedCommentFlag = cli.BoolFlag{
		Name:    "strip-embed-comment",
		EnvVars: []string{"EMBEDME_STRIP_EMBED_COMMENT"},
		Usage:   "remove the comment from the code block (only works with --stdout )",
	}
)
