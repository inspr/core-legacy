package handler

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type appHandler struct{}

// AppHandler is an empty object to separate route functions
var AppHandler appHandler = struct{}{}

// GetAllApps returns all Apps
func (ch *appHandler) GetAllApps() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// CreateInfo informs the data needed to create a App
func (ch *appHandler) CreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// CreateApp todo doc
func (ch *appHandler) CreateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// GetAppByRef todo doc
func (ch *appHandler) GetAppByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// UpdateApp todo doc
func (ch *appHandler) UpdateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// GetAllApps todo doc
func (ch *appHandler) DeleteApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}
