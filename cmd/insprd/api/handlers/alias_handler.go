package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
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
		data := models.AliasDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		ah.Memory.InitTransaction()

		err = ah.Memory.Alias().CreateAlias(data.Ctx, data.Target, &data.Alias)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleGet - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleGet() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		ah.Memory.InitTransaction()

		app, err := ah.Memory.Root().Alias().Get(data.Ctx, data.Key)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		defer ah.Memory.Cancel()

		rest.JSON(w, http.StatusOK, app)
	}
	return rest.Handler(handler)
}

// HandleUpdateAlias - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleUpdateAlias() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.Memory.InitTransaction()

		err = ah.Memory.Alias().UpdateAlias(data.Ctx, data.Target, &data.Alias)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDeleteAlias - handler that generates the rest.Handle
// func to manage the http request
func (ah *AliasHandler) HandleDeleteAlias() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasQueryDI{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.Memory.InitTransaction()

		err = ah.Memory.Alias().DeleteAlias(data.Ctx, data.Key)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.Channels().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
