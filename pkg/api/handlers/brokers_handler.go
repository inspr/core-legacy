package handler

import (
	"encoding/json"
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
			Default:   string(*def),
		}
		logger.Debug("current brokers:", zap.Any("brokers", brokers.Default))

		rest.JSON(w, http.StatusOK, brokers)
	}
	return rest.Handler(handler)
}

// HandlerCreate handles the creation of brokers in the insprd/cluster
func (bh *BrokerHandler) HandlerCreate() rest.Handler {
	logger.Info("handling Brokers' create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		var data models.BrokerDataDI

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("Unable to read body from request")
			rest.ERROR(w, err)
		}

		// parsing the bytes to specific handler config

		// err := bh.Brokers.Create(brokers.BrokerStatus(data.BrokerName), // TODO config parsed from the []bytes in data.FileContents)
		if err != nil {
			logger.Error("Unable to create broker with passed data")
			rest.ERROR(w, err)
		}
	}
	return rest.Handler(handler)
}
