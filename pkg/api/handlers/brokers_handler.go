package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/cmd/sidecars"
	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/rest"
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
			return
		}
		def, err := bh.Brokers.GetDefault()
		if err != nil {
			logger.Error("unable to obtain currently default brokers on cluster",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}
		brokers := &models.BrokersDI{
			Installed: available,
			Default:   def,
		}
		logger.Debug("current brokers:", zap.Any("brokers", brokers.Installed))

		rest.JSON(w, http.StatusOK, brokers)
	}
	return rest.Handler(handler)
}

// KafkaCreateHandler is the function that processes requests at the /brokers/kafka endpoint
func (bh *BrokerHandler) KafkaCreateHandler() rest.Handler {
	logger.Info("handling the brokers' kafka route request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		// decode into the bytes of yaml file
		var content models.BrokerConfigDI
		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		var kafkaConfig sidecars.KafkaConfig
		// parsing the bytes into a Kafka config structure
		err = yaml.Unmarshal(content.FileContents, &kafkaConfig)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if err = bh.Brokers.Create(
			&kafkaConfig,
		); err != nil {
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, http.StatusOK, nil)
	}
	return rest.Handler(handler)
}
