package lib

import (
	"io/ioutil"
	"log"
	"os"
)

var logger = log.New(ioutil.Discard, "govm/lib ", 0)

func SetVerbose(v bool) {
	if !v {
		logger.SetPrefix("")
		logger.SetFlags(0)
		logger.SetOutput(ioutil.Discard)
	} else {
		logger.SetPrefix("govm/lib ")
		logger.SetFlags(log.LstdFlags)
		logger.SetOutput(os.Stderr)
	}
}
