[![Go Report Card](https://goreportcard.com/badge/github.com/yi-jiayu/govm)](https://goreportcard.com/report/github.com/yi-jiayu/govm)

# govm
Manage multiple Go installations on Windows

## Installation
`go install -v github.com/yi-jiayu/govm` to clone this repository under `GOPATH/src/github.com/yi-jiayu/govm` and place a govm executable into `GOPATH/bin`.

Alternatively,
1. `go get -v github.com/yi-jiayu/govm` to clone the source into your GOPATH.
2. `cd` into the source directory: `cd $env:GOPATH/github.com/yi-jiayu/govm` (PowerShell) or `cd $GOPATH/github.com/yi-jiayu/govm` (Bash)
3. `go build` to compile the govm binary locally.

Precompiled binaries can also be found under [Releases](https://github.com/yi-jiayu/govm/releases).

## Usage
govm currently provides the following commands:
```
  help        Help about any command
  install     Install a new Go version
  list        List installed Go versions
  root        Print the current GOROOT
  uninstall   Uninstall an installed Go version
  use         Switch to a different Go version
  version     Show version information
```

## Background
Yes, I realise that there's not much need to manage multiple Go versions, but for a point in time I was struggling with various errors while building a Go on Google App Engine Standard Environment project and I thought that maybe using Go 1.6, which Go on Google App Engine runs, might help (it didn't).

## Assumptions
In order to limit the amount of magic done by this program, govm will not set any environment variables. You need to make sure that you set (or do not set) GOROOT and add GOROOT/bin to your PATH.

## How it works
govm downloads and saves Go binary installations to the govm install dir folder (default: `C:/`) as `GOVM_INSTALL_DIR/Go$version` and creates symlinks from GOROOT to the current version in use. If administrator permission in required to create a symlink, govm will trigger a UAC prompt.

For example, assuming `GOROOT` is `C:/Go` and the govm install dir is `C:/`, to use Go version 1.8 govm will create a directory symlink at `C:/Go` pointed at `C:/Go1.8`, which should contain Go version 1.8. When switching between Go versions, govm will delete the current `GOROOT` symlink and create a new one to the target version.

govm uses `cmd /c mklink /d` to create symlinks to Go installations.

## [Change log](CHANGELOG.md)
