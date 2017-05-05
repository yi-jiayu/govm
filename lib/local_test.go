package lib

import (
	"fmt"
	"log"
	"testing"
)

func TestInstalledGoVersions(t *testing.T) {
	installed, err := InstalledGoVersions("C:/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", installed)
}

func TestCurrentGoVersion(t *testing.T) {
	current, err := CurrentGoVersion(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(current)
}
