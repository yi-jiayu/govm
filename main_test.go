package main

import (
	"testing"
	"fmt"
	"log"
	"github.com/urfave/cli"
)

func TestInstalledGoVersions(t *testing.T) {
	vs, err := InstalledGoVersions()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vs {
		fmt.Println(v)
	}
}

func TestCurrentGoVersion(t *testing.T) {
	cv, err := CurrentGoVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cv)
}

func TestList(t *testing.T) {
	List(&cli.Context{})
}
