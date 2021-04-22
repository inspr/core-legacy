package utils

import (
	"io/ioutil"
)

// GetToken gets a token from a filepath, returns an empty token if file does not exist.
func GetToken(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte("")
	}
	return content
}
