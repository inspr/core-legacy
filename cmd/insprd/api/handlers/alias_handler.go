package handler

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// AliasHandler - contains handlers that uses the AliasMemory interface methods
type AliasHandler struct {
	*Handler
}

// NewAliasHandler - generates a new AliasHandler through the memoryManager interface
func (handler *Handler) NewAliasHandler() *AliasHandler {
	return &AliasHandler{
		handler,
	}
}

// HandleCreateAlias - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleCreateAlias() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

	}
	return rest.Handler(handler)
}

// HandleGet - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleGet() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

	}
	return rest.Handler(handler)
}

// HandleUpdateAlias - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleUpdateAlias() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

	}
	return rest.Handler(handler)
}

// HandleDeleteAlias - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleDeleteAlias() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

	}
	return rest.Handler(handler)
}
