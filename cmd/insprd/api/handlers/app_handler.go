package handler

import (
	"encoding/json"
	"io/ioutil"
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

// HandleCreateInfo informs the data needed to create a App
func (ah *AppHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// not implemented yet
	}
	return rest.Handler(handler)
}

// HandleCreateApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleCreateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.AppDI{}
		json.Unmarshal(body, &data)

		err = ah.CreateApp(&data.App, data.Ctx)
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.AppQueryDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		app, err := ah.GetApp(data.Query)
		if err != nil {
			rest.ERROR(w, http.StatusConflict, err)
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.AppDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ah.UpdateApp(&data.App, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleDeleteApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleDeleteApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.AppQueryDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
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
