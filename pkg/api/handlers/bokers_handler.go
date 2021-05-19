package handler

import (
	"net/http"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest"
)

// BrokerHandler - contains handlers that uses the BrokerManager interface methods
type BrokerHandler struct {
	*Handler
}

// NewBrokerHandler - returns the handle functions that regard brokers
func (handler *Handler) NewBrokerHandler() *BrokerHandler {
	return &BrokerHandler{
		handler,
	}
}

func (bh *BrokerHandler) HandleGet() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		brokers := &models.BrokersDi{
			Installed: bh.Brokers.GetAll(),
			Default:   string(bh.Brokers.GetDefault()),
		}

		rest.JSON(w, http.StatusOK, brokers)
	}
	return rest.Handler(handler)
}
