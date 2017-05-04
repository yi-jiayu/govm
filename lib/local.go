package lib

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	GOROOT = "C:/Go"
)

var (
	goroot = Goroot()

	PermissionDenied = errors.New("permission denied")
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

func CreateSymbolicLink(path, target string) error {
	var stderr bytes.Buffer

	cmd := exec.Command("cmd", fmt.Sprintf("/c mklink /d %s %s", path, target))
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if strings.Contains(stderr.String(), "You do not have sufficient privilege to perform this operation.") {
			return PermissionDenied
		} else {
			return err
		}
	}

	return nil
}

func CreateSymbolicLinkWithElevation(path, target string) error {
	powershellCmd := fmt.Sprintf(`Start-Process cmd -Wait -Verb RunAs -ArgumentList "cmd /c mklink /d %s %s"`, path, target)

	name := "powershell"
	args := []string{"-NoProfile", "-Command", powershellCmd}

	c := exec.Command(name, args...)

	err := c.Run()
	if err != nil {
		return err
	}

	return nil
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

	output, err := cmd2.Output()
	if err != nil {
		return gv, err
	}

	return string(output), nil
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

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	// create new go symlink
	target := filepath.Join(filepath.Dir(abs), "Go"+tv)

	err = CreateSymbolicLink(abs, target)
	if err != nil {
		if err == PermissionDenied {
			err := CreateSymbolicLinkWithElevation(abs, target)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func ExtractDownloadedGoVersion(dl string) (string, error) {
	var target string

	target, err := ioutil.TempDir("", "govm-install-")
	if err != nil {
		return target, err
	}

	reader, err := zip.OpenReader(dl)
	if err != nil {
		return target, err
	}
	defer reader.Close()

	// closure so that we don't run out of file descriptors
	extract := func(file *zip.File) error {
		var name string
		if !strings.HasPrefix(file.Name, "go") {
			return errors.New("unexpected file in archive")
		} else {
			name = strings.TrimPrefix(file.Name, "go")
		}

		path := filepath.Join(target, name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			// read file in zip
			fr, err := file.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			// create file to be written to
			os.MkdirAll(filepath.Dir(path), file.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			//logger.Printf("extracting %s to %s\n", file.Name, path)
			_, err = io.Copy(f, fr)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, file := range reader.File {
		err := extract(file)
		if err != nil {
			return target, err
		}
	}

	logger.Printf("extracted %s to %s", dl, target)
	return target, nil
}

func InstallGoVersion(version, source string) error {
	if !ValidateSemver(version) {
		return errors.New("invalid version string")
	}

	dest := filepath.Join(filepath.Dir(goroot), "Go"+version)
	_, err := os.Lstat(dest)
	if err == nil || !os.IsNotExist(err) {
		return errors.New("destination already exists")
	}

	logger.Printf("Installing to %s...\n", dest)
	err = os.Rename(source, dest)
	if err != nil {
		return err
	}

	logger.Printf("renamed %s to %s", source, dest)
	return nil
}

func UninstallGoVersion(version string) error {
	if !ValidateSemver(version) {
		return errors.New("invalid version string")
	}

	path := filepath.Join(filepath.Dir(goroot), "Go"+version)

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
