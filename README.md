# limesctl

[![GitHub release](https://img.shields.io/github/release/sapcc/limesctl.svg)](https://github.com/sapcc/limesctl/releases/latest)
[![Build Status](https://travis-ci.org/sapcc/limesctl.svg?branch=master)](https://travis-ci.org/sapcc/limesctl) 
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/limesctl)](https://goreportcard.com/report/github.com/sapcc/limesctl)
[![GoDoc](https://godoc.org/github.com/sapcc/limesctl?status.svg)](https://godoc.org/github.com/sapcc/limesctl)

limesctl is the CLI client for [Limes](limes).

## Installation

You can download the latest pre-compiled binary from the [releases page](https://github.com/sapcc/limesctl/releases/latest).

Alternatively, you can also build from source:

The only required build dependency is [Go](https://golang.org/). 

```
$ go get github.com/sapcc/limesctl
$ cd $GOPATH/src/github.com/sapcc/limesctl
$ make install
```

this will put the binary in `/usr/bin/limesctl` or `/usr/local/bin/limesctl` for macOS.


## Usage

To get an overview of all the operations that are available:

```
$ limesctl help
```

limesctl requires a valid Keystone Token for all its operation. Your mileage may vary depending on your authorization scope.

[limes]: https://github.com/sapcc/limes
