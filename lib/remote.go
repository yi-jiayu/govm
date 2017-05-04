package lib

import (
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	ENDPOINT = "https://storage.googleapis.com/golang/go$VERSION.windows-amd64.zip"
)

var (
	semverRegex = regexp.MustCompile(`[0-9].[0-9]+.?[0-9]*`)

	client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
)

func ValidateSemver(version string) bool {
	return semverRegex.MatchString(version)
}

func CheckRemoteVersionExists(version string) (bool, error) {
	var exists bool

	if !ValidateSemver(version) {
		return exists, errors.New("invalid version string")
	}

	resp, err := client.Head(strings.Replace(ENDPOINT, "$VERSION", version, 1))
	if err != nil {
		return exists, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		exists = true
	} else {
		exists = false
	}

	return exists, nil
}

func DownloadRemoteVersion(version string) (string, error) {
	var dl string

	if !ValidateSemver(version) {
		return dl, errors.New("invalid version string")
	}

	temp, err := ioutil.TempFile("", "govm-dl-")
	if err != nil {
		return dl, err
	}
	defer temp.Close()

	uri := strings.Replace(ENDPOINT, "$VERSION", version, 1)
	resp, err := client.Get(uri)
	if err != nil {
		return dl, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return dl, errors.New("status code error")
	}

	logger.Printf("downloaded %d bytes from %s\n", resp.ContentLength, uri)

	n, err := io.Copy(temp, resp.Body)
	if err != nil {
		return dl, err
	}
	dl = temp.Name()

	logger.Printf("copied %d bytes to %s\n", n, dl)
	return dl, nil
}
