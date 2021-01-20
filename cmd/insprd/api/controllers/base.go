package controller

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	mux           *http.ServeMux
	memoryManager memory.Manager
}
