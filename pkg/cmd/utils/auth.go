package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

// Authenticator is responsible for implementing the interface methods
// defined in the rest/request pkg.
type Authenticator struct {
	TokenPath string
}

var bearer = []byte("Bearer ")

// GetToken read the token from the file specified in the struct
// TokenPath and returns it's bytes
func (a Authenticator) GetToken() ([]byte, error) {
	token, err := ioutil.ReadFile(a.TokenPath)
	if err != nil {
		return nil, err
	}

	token = []byte(strings.TrimSpace(string(token)))
	token = append(bearer, token...)
	return token, nil
}

// SetToken receives a new token as a parameter and then writes it
// in the file specified in the TokenPath
func (a Authenticator) SetToken(token []byte) error {
	token = token[len(bearer):]
	err := ioutil.WriteFile(a.TokenPath, token, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
