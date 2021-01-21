package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			App meta.App `json:"app"`
			Ctx string   `json:"ctx"`
		}{}
		json.Unmarshal(body, &data)

		/// testing
		fmt.Println(data.App)
		fmt.Println(data.Ctx)

		err = ah.CreateApp(&data.App, data.Ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			Query string `json:"query"`
		}{}
		json.Unmarshal(body, &data)
		app, err := ah.GetApp(data.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}

		// respond with json
		fmt.Fprintf(w, "%v", app)
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleUpdateApp todo doc
func (ah *AppHandler) HandleUpdateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			App meta.App `json:"app"`
			Ctx string   `json:"ctx"`
		}{}
		json.Unmarshal(body, &data)

		err = ah.UpdateApp(&data.App, data.Ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			Query string `json:"query"`
		}{}
		json.Unmarshal(body, &data)
		err = ah.DeleteApp(data.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
