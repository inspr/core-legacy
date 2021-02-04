package sidecarserv

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

func mockServer(err error) *Server {
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

func (mr *mockReader) ReadMessage(channel string) (models.Message, error) {
	if mr.err != nil {
		return models.Message{}, mr.err
	}
	return models.Message{}, nil
}

func (mr *mockReader) CommitMessage(channel string) error {
	if mr.err != nil {
		return mr.err
	}
	return nil
}

func (mw *mockWriter) WriteMessage(channel string, msg models.Message) error {
	if mw.err != nil {
		return mw.err
	}
	return nil
}
