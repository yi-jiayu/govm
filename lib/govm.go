package lib

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
	"os/exec"
)

const (
	GOROOT = "C:/Go"
)

var (
	goroot = Goroot()
)

func Goroot() string {
	goroot := GOROOT

	if gr := os.Getenv("GOROOT"); gr != "" {
		gr, err := filepath.Abs(gr)
		if err == nil {
			goroot = gr
		}
	}

	return goroot
}

func IsSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0
}

func InstalledGoVersions() ([]string, error) {
	ivs := []string{}

	abs, err := filepath.Abs(goroot)
	if err != nil {
		return ivs, err
	}

	// look for folders in (dirname goroot) starting with "Go"
	dirname := filepath.Dir(abs)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return ivs, err
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "Go") {
			if len(file.Name()) > 2 {
				ivs = append(ivs, file.Name()[2:])
			}
		}
	}

	return ivs, nil
}

func CurrentGoVersion() (string, error) {
	// CurrentGoVersion is only valid if goroot is a directory symlink pointing to a sibling Go$VERSION directory
	// 1. check if goroot is a symlink
	// 2. check if goroot's target is a directory
	// 3. check if goroot's target is a sibling directory
	// 4. check if goroot's target is in InstalledGoVersions()
	var cv string

	fi, err := os.Lstat(goroot)
	if err != nil {
		// if goroot does not exist, we consider it as no current Go version and safe to create a new symlink
		if os.IsNotExist(err) {
			cv = ""
			return cv, nil
		} else {
			return cv, err
		}
	}

	if !IsSymlink(fi) {
		return cv, errors.New("GOROOT is not a symlink")
	}

	p, err := filepath.EvalSymlinks(goroot)
	if err != nil {
		return cv, err
	}

	target, err := filepath.Abs(p)
	if err != nil {
		return cv, err
	}

	fi, err = os.Stat(target)
	if err != nil {
		return cv, err
	}

	if !fi.IsDir() {
		return cv, errors.New("GOROOT is not a directory symlink")
	}

	dirname := filepath.Dir(goroot)
	if filepath.Dir(target) == dirname {
		cv = filepath.Base(target)[2:]

		ivs, err := InstalledGoVersions()
		if err != nil {
			return cv, err
		}

		for _, iv := range ivs {
			if cv == iv {
				return cv, nil
			}
		}
	}

	return cv, errors.New("GOROOT does not point to an installed Go version")
}

func GoVersionOutput() (string, error) {
	var gv string

	cmd2 := exec.Command("go", "version")
	cmd2.Stderr = os.Stderr

	outp, err := cmd2.Output()
	if err != nil {
		return gv, err
	}

	return string(outp), nil
}

func SwitchGoVersion(tv string) error {
	abs, err := filepath.Abs(goroot)
	if err != nil {
		return err
	}

	cv, err := CurrentGoVersion()
	if err != nil {
		return err
	}

	if cv != "" {
		// delete current Go symlink
		cmd := exec.Command("cmd", fmt.Sprintf("/c rmdir %s", abs))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdin

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	// create new go symlink
	target := filepath.Join(filepath.Dir(abs), "Go"+tv)

	cmd := exec.Command("cmd", fmt.Sprintf("/c mklink /d %s %s", abs, target))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdin

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
