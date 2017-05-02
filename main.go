package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	GOROOT = "C:/Go"
)

func init() {
	if os.Getenv("GOROOT") != "" {
		GOROOT = os.Getenv("GOROOT")
	}
}

func IsSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0
}

func InstalledGoVersions() ([]string, error) {
	installedVersions := []string{}

	// look for folders in dirname GOROOT starting with "Go"
	dirname := filepath.Dir(GOROOT)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return installedVersions, err
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "Go") {
			if len(file.Name()) > 2 {
				installedVersions = append(installedVersions, file.Name()[2:])
			}
		}
	}

	return installedVersions, nil
}

func CurrentGoVersion() (string, error) {
	// CurrentGoVersion is only valid if GOROOT is a directory symlink pointing to a sibling Go$(version) directory
	// 1. check if GOROOT is a symlink
	// 2. check if GOROOT's target is a directory
	// 3. check if GOROOT's target is a sibling directory
	// 4. check if GOROOT's target is in InstalledGoVersions()
	var cv string

	fi, err := os.Lstat(GOROOT)
	if err != nil {
		// if GOROOT does not exist, we consider it as no current Go version and safe to create a new symlink
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

	path, err := filepath.EvalSymlinks(GOROOT)
	if err != nil {
		return cv, err
	}

	target, err := filepath.Abs(path)
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

	dirname := filepath.Dir(GOROOT)
	if filepath.Dir(target) == dirname {
		cv = filepath.Base(target)[2:]

		vs, err := InstalledGoVersions()
		if err != nil {
			return cv, err
		}

		for _, v := range vs {
			if cv == v {
				return cv, nil
			}
		}
	}

	return cv, errors.New("GOROOT does not point to an installed Go version")
}

func List(c *cli.Context) error {
	vs, err := InstalledGoVersions()
	if err != nil {
		return err
	}

	cv, err := CurrentGoVersion()
	if err != nil {
		return err
	}

	for _, v := range vs {
		if v == cv {
			fmt.Printf("  * %s\n", v)
		} else {
			fmt.Printf("    %s\n", v)
		}
	}

	return nil
}

func Use(c *cli.Context) error {
	// 1. check if CurrentGoVersion is valid in our context
	// 2. delete GOROOT if CurrentGoVersion is not empty
	// 3. symlink GOROOT

	cv, err := CurrentGoVersion()
	if err != nil {
		return err
	}

	if c.NArg() != 1 {
		return cli.ShowCommandHelp(c, "use")
	}

	tv := c.Args().Get(0)

	fmt.Printf("Current Go version is: %s\n", cv)
	fmt.Printf("You are trying to use Go version: %s\n", tv)

	vs, err := InstalledGoVersions()
	if err != nil {
		return err
	}

	found := false
	for _, v := range vs {
		if v == tv {
			found = true
			break
		}
	}

	if found {
		if cv == tv {
			fmt.Printf("Go version %s is already the currently active version.\n", cv)

			return nil
		} else {
			fmt.Printf("Changing to Go version %s...\n", tv)
			fmt.Fprintln(os.Stderr, "not implemented")

			return nil
		}
	} else {
		fmt.Printf("Go version %s is not currently installed.\n", tv)

		return nil
	}
}

func Root(c *cli.Context) error {
	if c.NArg() == 0 {
		fp, err := filepath.Abs(GOROOT)
		if err != nil {
			return err
		}

		fmt.Printf("GOROOT is currently set to: %s\n", fp)

		return nil
	} else if c.NArg() == 1 {
		fmt.Printf("You are trying to set GOROOT to: %s\n", c.Args().Get(0))
		fmt.Fprintln(os.Stderr, "not implemented")

		return nil
	} else {
		return cli.ShowCommandHelp(c, "root")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "gvm"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install a new Go version",
			Action: func(c *cli.Context) error {
				fmt.Fprintln(os.Stderr, "not implemented")
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List installed Go versions",
			Action:  List,
		},
		{
			Name:   "root",
			Usage:  "Display the current GOROOT",
			Action: Root,
		},
		{
			Name:      "use",
			Usage:     "Switch to the specified Go version",
			ArgsUsage: "[version]",
			Action:    Use,
		},
		{
			Name:  "version",
			Usage: "Display the current gvm version",
			Action: func(c *cli.Context) error {
				fmt.Fprintln(os.Stderr, "not implemented")
				return nil
			},
		},
	}

	app.Run(os.Args)
}
