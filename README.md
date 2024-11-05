# Porty

![Go Tests](https://github.com/softwarespot/porty/actions/workflows/go.yml/badge.svg)

## Prerequisites

-   go 1.23.0 or above
-   make (if you want to use the `Makefile` provided)

## Installation

Build the binary `porty` executable to the directory `./bin` i.e. `./bin/porty`.

```bash
make
```

Install i.e. copy the executable `./bin/porty` to `$HOME/bin` (if it exists).

```bash
make install
```

## Dependencies

## Linting

Docker

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v --tests=false --disable-all -E durationcheck,errorlint,exhaustive,gocritic,gosimple,ineffassign,misspell,predeclared,revive,staticcheck,unparam,unused,whitespace --max-issues-per-linter=10000 --max-same-issues=10000
```

Local

```bash
golangci-lint run --tests=false --disable-all -E durationcheck,errorlint,exhaustive,gocritic,gosimple,ineffassign,misspell,predeclared,revive,staticcheck,unparam,unused,whitespace --max-issues-per-linter=10000 --max-same-issues=10000
```

## License

The code has been licensed under the [MIT](https://opensource.org/license/mit) license.
