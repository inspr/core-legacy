package controllers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/inspr/inspr/pkg/ierrors"
	"go.uber.org/zap"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux     *http.ServeMux
	logger  *zap.Logger
	privKey *rsa.PrivateKey
}

// Init - configures the server
func (s *Server) Init() {
	var err error

	keyPem := []byte(os.Getenv("JWT_PRIVATE_KEY"))
	privKey, _ := pem.Decode(keyPem)

	var privPemBytes []byte
	if privKey.Type != "RSA PRIVATE KEY" {
		s.logger.Error("RSA private key is of the wrong type")
		err = ierrors.NewError().InternalServer().Message("RSA private key is of the wrong type").Build()
		panic(err)
	}

	privPemBytes = privKey.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			s.logger.Error("Unable to parse RSA private key")
			err = ierrors.NewError().InternalServer().Message("Unable to parse RSA private key").Build()
			panic(err)
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		s.logger.Error("Unable to parse RSA private key")
		err = ierrors.NewError().InternalServer().Message("Unable to parse RSA private key").Build()
		panic(err)
	}

	err = privateKey.Validate()
	if err != nil {
		s.logger.Error("Unable to validate private key", zap.Any("error", err))
		err = ierrors.NewError().InternalServer().Message("Unable to validate private key").Build()
		panic(err)
	}

	s.privKey = privateKey
	s.Mux = http.NewServeMux()
	s.logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "Auth-provider")))
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("insprd rest api is up! Listening on port: %s\n", addr)
	s.logger.Fatal("Authsvc crashed: ", zap.Any("error", http.ListenAndServe(addr, s.Mux)))
}
