package auth

import (
	"io/ioutil"
	"os"
	"strings"

	"inspr.dev/inspr/pkg/ierrors"
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
		return nil, ierrors.New(err).InvalidFile()
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
		return ierrors.New(err).InvalidFile()
	}
	return nil
}
