## embedme

[![Build Status](https://github.com/romnn/embedme/workflows/test/badge.svg)](https://github.com/romnn/embedme/actions)
[![GitHub](https://img.shields.io/github/license/romnn/embedme)](https://github.com/romnn/embedme)
[![Test Coverage](https://codecov.io/gh/romnn/embedme/branch/master/graph/badge.svg)](https://codecov.io/gh/romnn/embedme)
[![Release](https://img.shields.io/github/release/romnn/embedme)](https://github.com/romnn/embedme/releases/latest)

t.b.a

```bash
go install github.com/romnn/embedme/cmd/embedme
embedme serve --generate
```

You can also download pre-built [release binaries](https://github.com/romnn/embedme/releases).

For a list of options, run with `--help`.

### Development

#### Tools

Before you get started, make sure you have installed the following tools:

```bash
$ python3 -m pip install pre-commit bump2version invoke
$ go install github.com/kyoh86/richgo@latest
$ go install golang.org/x/tools/cmd/goimports@latest
$ go install golang.org/x/lint/golint@latest
$ go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
```

Please always make sure code checks pass:

```bash
inv pre-commit
```

#### TODO

- use different color library?
- refactor the code
- add a lot of tests!
- fix the source line numbers

- done
  - basic parsing
  - add regex that can identify commands
