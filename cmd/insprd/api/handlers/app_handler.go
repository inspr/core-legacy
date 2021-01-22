package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
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
func (ah *AppHandler) HandleGetAllApps() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// not implemented yet
	}
	return rest.Handler(handler)
}

// HandleCreateInfo informs the data needed to create a App
func (ah *AppHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// not implemented yet
	}
	return rest.Handler(handler)
}

// HandleCreateApp todo doc
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

// HandleGetAppByRef todo doc
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

// HandleUpdateApp todo doc
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

// HandleDeleteApp todo doc
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
