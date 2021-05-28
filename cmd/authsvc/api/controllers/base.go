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
	s.logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "auth-provider")))

	keyPem, ok := os.LookupEnv("JWT_PRIVATE_KEY")
	if !ok {
		panic("[ENV VAR] JWT_PRIVATE_KEY not found")
	}
	privKey, _ := pem.Decode([]byte(keyPem))

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
			s.logger.Error("unable to parse RSA private key",
				zap.Any("error", err))

			err = ierrors.NewError().InternalServer().Message("error parsing RSA private key: %v", err).Build()
			panic(err)
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		s.logger.Error("unable to parse RSA private key")
		err = ierrors.NewError().InternalServer().Message("unable to parse RSA private key").Build()
		panic(err)
	}

	err = privateKey.Validate()
	if err != nil {
		s.logger.Error("unable to validate private key", zap.Any("error", err))
		err = ierrors.NewError().InternalServer().Message("unable to validate private key").Build()
		panic(err)
	}

	s.privKey = privateKey
	s.Mux = http.NewServeMux()
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("authsvc rest api is up! Listening on port: %s\n", addr)
	s.logger.Fatal("authsvc crashed: ", zap.Any("error", http.ListenAndServe(addr, s.Mux)))
}
