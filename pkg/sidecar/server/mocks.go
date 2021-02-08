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

func (mr *mockReader) ReadMessage(channel string) (models.Message, error) {
	if mr.err != nil {
		return models.Message{
			Commit:  true,
			Channel: channel,
			Data:    "mock_data",
			Error:   mr.err,
		}, mr.err
	}
	return models.Message{
		Commit:  true,
		Channel: channel,
		Data:    "mock_data",
		Error:   nil,
	}, nil
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
