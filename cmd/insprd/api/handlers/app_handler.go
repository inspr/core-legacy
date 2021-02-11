package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// AppHandler - contains handlers that uses the AppMemory interface methods
type AppHandler struct {
	memory.AppMemory
}

// NewAppHandler - generates a new AppHandler through the memoryManager interface
func NewAppHandler(memManager memory.Manager) *AppHandler {
	return &AppHandler{
		AppMemory: memManager.Apps(),
	}
}

// HandleCreateApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleCreateApp() rest.Handler {

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ah.CreateApp(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleGetAppByRef - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleGetAppByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		app, err := ah.GetApp(data.Query)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, app)
	}
	return rest.Handler(handler)
}

// HandleUpdateApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleUpdateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ah.UpdateApp(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
	return rest.Handler(handler)
}

// HandleDeleteApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleDeleteApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ah.DeleteApp(data.Query)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
