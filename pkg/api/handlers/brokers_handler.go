package handler

import (
	"net/http"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
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

// HandleGet returns the get handler for brokers
func (bh *BrokerHandler) HandleGet() rest.Handler {
	logger.Info("handling Brokers get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		available, err := bh.Brokers.GetAll()
		if err != nil {
			logger.Error("unable to obtain currently available brokers on cluster",
				zap.Any("error", err))
			rest.ERROR(w, err)
		}
		def, err := bh.Brokers.GetDefault()
		if err != nil {
			logger.Error("unable to obtain currently default brokers on cluster",
				zap.Any("error", err))
			rest.ERROR(w, err)
		}
		brokers := &models.BrokersDI{
			Installed: available,
			Default:   def,
		}
		logger.Debug("current brokers:", zap.Any("brokers", brokers.Default))

		rest.JSON(w, http.StatusOK, brokers)
	}
	return rest.Handler(handler)
}
