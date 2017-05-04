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
In order to limit the amount of magic done by this program, govm makes a few assumptions about your environment and will work best when these assumptions hold. govm tries to be conservative and not do anything if these assumptions are wrong, but certain assumptions such as assuming `basename $GOROOT` to be "Go" are hardcoded.

For best results,
1. Set your `GOROOT` to `C:/Go` or leave it unset.
2. Add `GOROOT/bin` to your `PATH`.

## Setup
govm searches in the `dirname $GOROOT` directory for sibling directories to `GOROOT` with names matching `Go*` such as `Go1.6.4` or `Go1.8`. For example, if `GOROOT` is `C:/Go`, govm expects `C:/Go1.6.4` to contain Go version 1.6.4 and `C:/Go1.8` to contain Go version 1.8.

## How it works
All govm does is create a symlink at `GOROOT` which points to a folder containing a specific Go version. For example, assuming `GOROOT` is `C:/Go`, to use Go version 1.8 govm will create a directory symlink at `C:/Go` pointed at `C:/Go1.8`, which should contain Go version 1.8. When switching between Go versions, govm will delete the current `GOROOT` symlink and create a new one to the target version, but only if `GOROOT` is indeed a symlink pointing at a Go installation. govm simply execs `cmd /c mklink /d` to create directory symlinks, and `cmd /c rmdir` to delete directory symlinks.

## [Change log](CHANGELOG.md)
