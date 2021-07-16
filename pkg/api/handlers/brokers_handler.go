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
	logger *zap.Logger
}

// NewBrokerHandler - returns the handle functions that regard brokers
func (handler *Handler) NewBrokerHandler() *BrokerHandler {
	return &BrokerHandler{
		Handler: handler,
		logger:  logger.With(zap.String("subSection", "brokers")),
	}
}

// HandleGet returns the get handler for brokers
func (bh *BrokerHandler) HandleGet() rest.Handler {
	l := bh.logger.With(zap.String("operation", "get"))
	l.Info("received brokers get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		brokers, err := bh.Memory.Brokers().Get()
		if err != nil {
			l.Error("unable to obtain currently available brokers on cluster", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l.Debug("current brokers:", zap.Strings("brokers", brokers.Available))

		rest.JSON(w, http.StatusOK, brokers)
	}
	return rest.Handler(handler)
}

// KafkaCreateHandler is the function that processes requests at the /brokers/kafka endpoint
func (bh *BrokerHandler) KafkaCreateHandler() rest.Handler {
	l := bh.logger.With(zap.String("operation", "create"), zap.String("broker", "kafka"))
	l.Info("received kafka broker create request")
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
			l.Error("unable to unmarshall config", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		if err = bh.Memory.Brokers().Create(
			&kafkaConfig,
		); err != nil {
			l.Error("error creating kafka broker on memory", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, http.StatusOK, nil)
	}
	return rest.Handler(handler)
}
