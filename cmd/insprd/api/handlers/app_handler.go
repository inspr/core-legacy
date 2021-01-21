package handler

import (
	"fmt"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// AppHandler todo doc
type AppHandler struct {
	memory.AppMemory
}

// NewAppHandler todo doc
func NewAppHandler(memManager memory.Manager) *AppHandler {
	return &AppHandler{
		AppMemory: memManager.Apps(),
	}
}

// HandleGetAllApps returns all Apps
func (ch *AppHandler) HandleGetAllApps() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateInfo informs the data needed to create a App
func (ch *AppHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateApp todo doc
func (ch *AppHandler) HandleCreateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// process json and other things
		newApp := &meta.App{}
		ctx := ""
		err := ch.CreateApp(newApp, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
	}
	return rest.Handler(handler)
}

// HandleGetAppByRef todo doc
func (ch *AppHandler) HandleGetAppByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// process json and other things
		newApp := &meta.App{}
		ctx := ""
		err := ch.CreateApp(newApp, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
	}
	return rest.Handler(handler)
}

// HandleUpdateApp todo doc
func (ch *AppHandler) HandleUpdateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// process json and other things
		newApp := &meta.App{}
		ctx := ""
		err := ch.CreateApp(newApp, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
	}
	return rest.Handler(handler)
}

// HandleDeleteApp todo doc
func (ch *AppHandler) HandleDeleteApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// process json and other things
		newApp := &meta.App{}
		ctx := ""
		err := ch.CreateApp(newApp, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
	}
	return rest.Handler(handler)
}
