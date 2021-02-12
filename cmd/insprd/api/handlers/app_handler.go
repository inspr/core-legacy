package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
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
		tree.GetTreeMemory().InitTransaction()
		if !data.DryRun {
			defer tree.GetTreeMemory().Commit()
		} else {
			defer tree.GetTreeMemory().Cancel()
		}
		err = ah.CreateApp(&data.App, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		diff, err := tree.GetTreeMemory().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		fmt.Println(json.Marshal(diff))
		rest.JSON(w, http.StatusOK, diff)
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
		app, err := ah.GetApp(data.Ctx)
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
		tree.GetTreeMemory().InitTransaction()
		if !data.DryRun {
			defer tree.GetTreeMemory().Commit()
		} else {
			defer tree.GetTreeMemory().Cancel()
		}
		err = ah.UpdateApp(&data.App, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
		}
		diff, err := tree.GetTreeMemory().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
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
		tree.GetTreeMemory().InitTransaction()
		if !data.DryRun {
			defer tree.GetTreeMemory().Commit()
		} else {
			defer tree.GetTreeMemory().Cancel()
		}
		err = ah.DeleteApp(data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		diff, err := tree.GetTreeMemory().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
