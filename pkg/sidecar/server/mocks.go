package sidecarserv

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// MockServer returns a mocked server to do tests
func MockServer(err error) *Server {
	return &Server{
		Mux:    http.NewServeMux(),
		Reader: &mockReader{err},
		Writer: &mockWriter{err},
	}
}

type mockReader struct {
	err error
}
type mockWriter struct {
	err error
}

func (mr *mockReader) ReadMessage() (models.BrokerData, error) {
	if mr.err != nil {
		return models.BrokerData{Message: models.Message{Data: "mock_data"}, Channel: ""}, mr.err
	}
	return models.BrokerData{Message: models.Message{Data: "mock_data"}, Channel: ""}, nil
}

func (mr *mockReader) CommitMessage() error {
	if mr.err != nil {
		return mr.err
	}
	return nil
}

func (mw *mockWriter) WriteMessage(channel string, msg interface{}) error {
	if mw.err != nil {
		return mw.err
	}
	return nil
}