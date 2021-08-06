package brokers

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/sidecars"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/sidecars/models"
)

// Get returns the brokers configured data
func (bmm *brokerMemoryManager) Get() (*apimodels.BrokersDI, error) {
	l := logger.With(zap.String("operation", "get"))
	l.Debug("received broker get request")
	bmm.available.Lock()
	l.Debug("available mutex locked", zap.String("type", "mutex"))

	defer l.Debug("available mutex unlocked", zap.String("type", "mutex"))
	defer bmm.available.Unlock()

	mem, err := bmm.get()
	if err != nil {
		l.Debug("unable to get memory manager")
		return nil, err
	}
	l.Debug("retrieved brokers, returning")
	return &apimodels.BrokersDI{
		Available: mem.Available.Brokers(),
		Default:   mem.Default,
	}, nil
}

func (bmm *brokerMemoryManager) get() (*brokers.Brokers, error) {
	if bmm.broker == nil {
		return nil, ierrors.NewError().Message("broker status memory is empty").Build()
	}
	return bmm.broker, nil
}

// Create configures a new broker on insprd
func (bmm *brokerMemoryManager) Create(config brokers.BrokerConfiguration) error {
	l := logger.With(
		zap.String("operation", "create"),
		zap.Any("configs", config),
	)
	bmm.available.Lock()
	l.Debug("available mutex locked", zap.String("type", "mutex"))

	defer l.Debug("available mutex unlocked", zap.String("type", "mutex"))
	defer bmm.available.Unlock()

	l.Info("creating new broker")
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	broker := config.Broker()

	if _, ok := mem.Available[broker]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}

	var factory models.SidecarFactory
	switch broker {
	case brokers.Kafka:
		l.Debug("configuring broker as kafka brojer")
		obj, _ := config.(*sidecars.KafkaConfig)
		factory = sidecars.KafkaToDeployment(*obj)
	default:
		l.Debug("found unsupported broker config, rejecting request")
		return ierrors.NewError().Message("broker %s is not supported", broker).Build()
	}

	l.Debug("subscribing broker to sidecar factory")
	err = bmm.Factory().Subscribe(broker, factory)
	if err != nil {
		l.Error("unable to subscribe broker")
		return err
	}

	mem.Available[broker] = config
	if mem.Default == "" {
		l.Debug("no default broker found - setting broker as default")
		bmm.SetDefault(broker)
	}
	return nil
}

// SetDefault sets a previously configured broker as insprd's default broker
func (bmm *brokerMemoryManager) SetDefault(broker string) error {
	l := logger.With(zap.String("operation", "set-default"), zap.String("broker", broker))
	bmm.def.Lock()
	l.Debug("def mutex locked", zap.String("type", "mutex"))

	defer l.Debug("def mutex unlocked", zap.String("type", "mutex"))
	defer bmm.def.Unlock()

	l.Debug("received default broker change request")
	mem, err := bmm.get()
	if err != nil {
		l.Debug("unable to get broker memory")
		return err
	}

	if _, ok := mem.Available[broker]; !ok {
		l.Debug("broker not configured")
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	mem.Default = broker
	l.Debug("default broker set")
	return nil
}

// Factory provides the struct implementation for Sidecarfactory
func (bmm *brokerMemoryManager) Factory() SidecarManager {
	return bmm.factory
}

//Configs returns the configurations for a given broker
func (bmm *brokerMemoryManager) Configs(broker string) (brokers.BrokerConfiguration, error) {
	l := logger.With(zap.String("operation", "get-configs"), zap.String("broker", broker))
	bmm.available.Lock()
	l.Debug("available mutex locked", zap.String("type", "mutex"))
	defer l.Debug("available mutex unlocked", zap.String("type", "mutex"))
	defer bmm.available.Unlock()

	l.Info("getting config for broker sidecar")

	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}

	config, ok := mem.Available[broker]
	if !ok {
		return nil, ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	return config, nil
}
