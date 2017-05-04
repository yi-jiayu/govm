package lib

import (
	"fmt"
	"os"
	"testing"
)

type TestCase struct {
	Input    interface{}
	Expected interface{}
}

func TestValidateSemver(t *testing.T) {
	testCases := []struct {
		Version  string
		Expected bool
	}{
		{
			"1",
			false,
		},
		{
			"1.8",
			true,
		},
		{
			"1.6.4",
			true,
		},
	}

	fail := false
	for _, tc := range testCases {
		if ValidateSemver(tc.Version) != tc.Expected {
			fmt.Fprintf(os.Stderr, "Expected ValidateSemver(\"%s\") to be %t, got %t.\n", tc.Version, tc.Expected, ValidateSemver(tc.Version))
			fail = true
		}
	}

	if fail {
		t.Fail()
	}
}

func TestCheckRemoteVersionExists(t *testing.T) {
	var version string
	var exists bool
	var err error

	version = "1"
	exists, err = CheckRemoteVersionExists(version)
	if err == nil || err.Error() != "invalid version string" {
		t.Fail()
	}

	version = "1.8"
	exists, err = CheckRemoteVersionExists(version)
	if err != nil || exists != true {
		t.Fail()
	}

	version = "1.9"
	exists, err = CheckRemoteVersionExists(version)
	if err != nil || exists != false {
		t.Fail()
	}

	version = "1.6.4"
	exists, err = CheckRemoteVersionExists(version)
	if err != nil || exists != true {
		t.Fail()
	}

	version = "1.6.9"
	exists, err = CheckRemoteVersionExists(version)
	if err != nil || exists != false {
		t.Fail()
	}
}

//func TestDownloadRemoteVersion(t *testing.T) {
//	DownloadRemoteVersion("1.8.1")
//}
