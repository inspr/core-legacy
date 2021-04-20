package utils

import (
	"io/ioutil"
)

func GetToken(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte("")
	}
	return content
}
