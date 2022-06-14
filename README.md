# limesctl

[![GitHub Release](https://img.shields.io/github/v/release/sapcc/limesctl)](https://github.com/sapcc/limesctl/releases/latest)
[![CI](https://github.com/sapcc/limesctl/actions/workflows/ci.yaml/badge.svg)](https://github.com/sapcc/limesctl/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/limesctl)](https://goreportcard.com/report/github.com/sapcc/limesctl)

`limesctl` is the command-line interface for [Limes](https://github.com/sapcc/limes).

## Installation

We provide pre-compiled binaries for the [latest release](https://github.com/sapcc/limesctl/releases/latest).

Alternatively, you can build with `make` or install with `make install`. The latter
understands the conventional environment variables for choosing install locations:
`DESTDIR` and `PREFIX`.

### Homebrew

In addition to macOS, the `brew` package will also work with Homebrew on Linux.

```
$ brew tap sapcc/limesctl https://github.com/sapcc/limesctl.git
$ brew install sapcc/limesctl/limesctl
```

### Go

```
$ go install github.com/sapcc/limesctl@latest
```

## Usage

For usage instructions:

```
$ limesctl --help
```

**Note**: `limesctl` requires the full set of OpenStack auth environment
variables. See [documentation for openstackclient](https://docs.openstack.org/python-openstackclient/latest/cli/man/openstack.html) for details.
