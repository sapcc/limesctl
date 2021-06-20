# limesctl

[![GitHub Release](https://img.shields.io/github/v/release/sapcc/limesctl)](https://github.com/sapcc/limesctl/releases/latest)
[![CI](https://github.com/sapcc/limesctl/actions/workflows/ci.yaml/badge.svg)](https://github.com/sapcc/limesctl/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/limesctl)](https://goreportcard.com/report/github.com/sapcc/limesctl)

`limesctl` is the command-line interface for [Limes](https://github.com/sapcc/limes).

## Usage

You can download pre-compiled binaries for the [latest release](https://github.com/sapcc/limesctl/releases/latest).

Alternatively, you can build with `make`, install with `make install`, or `go get`.

For usage instructions:

```
$ limesctl --help
```

**Note**: `limesctl` requires the full set of OpenStack auth environment
variables. See [documentation for openstackclient](https://docs.openstack.org/python-openstackclient/latest/cli/man/openstack.html) for details.
