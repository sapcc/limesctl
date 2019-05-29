# limesctl

[![GitHub release](https://img.shields.io/github/release/sapcc/limesctl.svg)](https://github.com/sapcc/limesctl/releases/latest)
[![Build Status](https://travis-ci.org/sapcc/limesctl.svg?branch=master)](https://travis-ci.org/sapcc/limesctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/limesctl)](https://goreportcard.com/report/github.com/sapcc/limesctl)

`limesctl` is the command-line client for [Limes](https://github.com/sapcc/limes).

## Installation

### Installer script

The simplest way to install `limesctl` on Linux or macOS is to run:

```
$ sh -c "$(curl -sL git.io/limesctl)"
```

This will put the binary in `/usr/local/bin/limesctl`

### Pre-compiled binaries

Pre-compiled binaries for Linux and macOS are avaiable on the
[releases page](https://github.com/sapcc/limesctl/releases/latest).

The binaries are static executables.

### Building from source

The only required build dependency is Go 1.11 or above.

```
$ git clone https://github.com/sapcc/limesctl.git
$ cd limesctl
$ make install
```

This will put the binary in `/usr/local/bin/limesctl`

Alternatively, you can also build `limesctl` directly with the `go get` command
without manually cloning the repository:

```
$ go get -u github.com/sapcc/limesctl
```

This will put the binary in `$GOPATH/bin/limesctl`

## Usage

To get an overview of all the operations that are available:

```
$ limesctl help
```

**Note**: `limesctl` requires a valid Keystone Token for all its operation.
Your mileage may vary depending on your authorization scope.
