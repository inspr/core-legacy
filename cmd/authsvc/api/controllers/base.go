package controllers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux     *http.ServeMux
	logger  *zap.Logger
	privKey *rsa.PrivateKey
	alevel  *zap.AtomicLevel
}

// Init - configures the server
func (s *Server) Init() {
	var err error
	s.logger, _ = logs.Logger(zap.Fields(zap.String("section", "server")))

	keyPem, ok := os.LookupEnv("JWT_PRIVATE_KEY")
	if !ok {
		panic("[ENV VAR] JWT_PRIVATE_KEY not found")
	}
	privKey, _ := pem.Decode([]byte(keyPem))

	var privPemBytes []byte
	if privKey.Type != "RSA PRIVATE KEY" {
		s.logger.Error("RSA private key is of the wrong type")
		err = ierrors.New("RSA private key is of the wrong type").
			InternalServer()
		panic(err)
	}

	privPemBytes = privKey.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			s.logger.Error("unable to parse RSA private key",
				zap.Any("error", err))

			err = ierrors.Wrap(
				ierrors.New(err).InternalServer(),
				"error parsing RSA private key",
			)
			panic(err)
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		s.logger.Error("unable to parse RSA private key")
		err = ierrors.New("unable to parse RSA private key").InternalServer()
		panic(err)
	}

	err = privateKey.Validate()
	if err != nil {
		s.logger.Error("unable to validate private key", zap.Any("error", err))
		err = ierrors.New("unable to validate private key").InternalServer()
		panic(err)
	}

	s.privKey = privateKey
	s.Mux = http.NewServeMux()
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	s.logger.Info("authsvc rest api is up!", zap.String("Port", addr))
	s.logger.Fatal(
		"authsvc crashed: ",
		zap.Any("error", http.ListenAndServe(addr, s.Mux)),
	)
}
