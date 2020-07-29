# limesctl

[![GitHub release](https://img.shields.io/github/release/sapcc/limesctl.svg)](https://github.com/sapcc/limesctl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/limesctl)](https://goreportcard.com/report/github.com/sapcc/limesctl)

`limesctl` is the command-line client for [Limes](https://github.com/sapcc/limes).

## Installation

### Pre-compiled binaries

Pre-compiled binaries for Linux and macOS are avaiable on the
[releases page](https://github.com/sapcc/limesctl/releases/latest).

### Building from source

The only required build dependency is [Go](https://golang.org/).

```
$ git clone https://github.com/sapcc/limesctl.git
$ cd limesctl
$ make install
```

This will put the binary in `/usr/bin/limesctl` on Linux and
`/usr/local/bin/limesctl` for macOS.

Alternatively, you can also build `limesctl` directly with the `go get` command
without manually cloning the repository:

```
$ go get -u github.com/sapcc/limesctl
```

This will put the binary in `$GOPATH/bin/limesctl`.

## Usage

To get an overview of all the commands:

```
$ limesctl --help
```

**Note**: `limesctl` requires a valid Keystone token for all its operations.
