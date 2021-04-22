package utils

import (
	"io/ioutil"
	"os"
)

type Authenticator struct {
	tokenPath string
}

var bearer = []byte("Bearer ")

func (a Authenticator) GetToken() ([]byte, error) {
	token, err := ioutil.ReadFile(a.tokenPath)
	if err != nil {
		return nil, err
	}
	token = append(bearer, token...)
	return token, nil
}

func (a Authenticator) SetToken(token []byte) error {
	token = token[len(bearer):]
	err := ioutil.WriteFile(a.tokenPath, token, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
