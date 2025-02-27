# Porty

![Go Tests](https://github.com/softwarespot/porty/actions/workflows/go.yml/badge.svg)

This application is a ports manager, for registering an application name and getting the next available port. The port range is 8000-9000.

## Prerequisites

- go 1.24.0 or above
- make (if you want to use the `Makefile` provided)

## Installation

### Binaries

`porty` is available on Linux, macOS and Windows platforms.
Binaries for Linux, Windows and Mac are available as tarballs in the [release page](https://github.com/softwarespot/porty/releases).

### Go Install

Install to the Go `bin` directory e.g. `$HOME/go/bin/`.

```bash
go install github.com/softwarespot/porty@latest
```

### Usage

### Init

Initialize the ports database.

```bash
# As text
porty init

# As JSON
porty init --json
```

### Get

Get or assign the next available port number for the app name.

```bash
# As text
porty get myapp

# As JSON
porty get myapp --json
```

### List

List all port numbers.

```bash
# As text
porty list

# As JSON
porty list --json
```

### Next

Get the next available port number, without assigning to an app name.

```bash
# As text
porty next

# As JSON
porty next --json
```

### Remove

Remove a port number for the app name.

```bash
# As text
porty remove myapp

# As JSON
porty remove myapp --json
```

### Who

Get who has been assigned to a port number.

```bash
# As text
porty who 8001

# As JSON
porty who 8001 --json
```

### Version

Display the version of the application.

```bash
# As text
porty version

# As JSON
porty version --text
```

### Help

Display the help text and exit.

```bash
porty --help
```

## Autocompletion

Adds autocompletion for the application's commands.

Add the following line to the `.bashrc` file; otherwise, refer to [kubectl's](https://kubernetes.io/docs/tasks/tools/install-kubectl/#optional-kubectl-configurations) documentation about locations for other OSs

- For `bash`

```bash
source <(porty completion bash)
```

- For `fish`

```bash
source <(porty completion fish)
```

- For `zsh`

```bash
source <(porty completion zsh)
```

### Example (bash)

Example of adding autocompletion for `bash` to the `.bashrc` file.

```bash
echo 'source <(porty completion bash)' >> ~/.bashrc
source ~/.bashrc
```

## Dependencies

**IMPORTANT:** 3rd party dependencies are used.

I only ever use dependencies when it's say an adapter for
an external service e.g. Redis, MySQL or Prometheus.

## Linting

Docker

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run --tests=false --default=none -E durationcheck,errorlint,exhaustive,gocritic,ineffassign,misspell,predeclared,revive,staticcheck,unparam,unused,whitespace --max-issues-per-linter=10000 --max-same-issues=10000
```

Local

```bash
golangci-lint run --tests=false --default=none -E durationcheck,errorlint,exhaustive,gocritic,ineffassign,misspell,predeclared,revive,staticcheck,unparam,unused,whitespace --max-issues-per-linter=10000 --max-same-issues=10000
```

## License

The code has been licensed under the [MIT](https://opensource.org/license/mit) license.
