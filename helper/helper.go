package helper

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

// reusing this all over the place
func StringFromFile(filename string) (string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(bytes.NewBuffer(b).String(), " \n"), nil
}

// Convenience Method for checking files
func IsFile(filename string) bool {
	if fi, err := os.Stat(filename); err == nil && fi.Mode().IsRegular() {
		return true
	}
	return false
}
