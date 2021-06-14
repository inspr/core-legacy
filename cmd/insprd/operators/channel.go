package operators

import (
	"context"
	"reflect"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/cmd/insprd/operators/kafka"
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	metabrokers "github.com/inspr/inspr/pkg/meta/brokers"
)

//GenOp is a general channel operator for dinamically selecting required operator
type GenOp struct {
	brokers brokers.Manager
	memory  memory.Manager
	configs map[string]struct {
		config metabrokers.BrokerConfiguration
		op     ChannelOperatorInterface
	}
}

//NewGeneralOperator creates an instance of GenOp for a given broker and memory manager
func NewGeneralOperator(brokers brokers.Manager, memory memory.Manager) *GenOp {
	return &GenOp{
		brokers: brokers,
		memory:  memory,
		configs: make(map[string]struct {
			config metabrokers.BrokerConfiguration
			op     ChannelOperatorInterface
		}),
	}
}

func (g GenOp) getOperator(scope string, name string) (ChannelOperatorInterface, error) {
	channel, _ := g.memory.Channels().Get(scope, name)
	broker := channel.Spec.SelectedBroker

	config, err := g.brokers.Configs(broker)
	if err != nil {
		return nil, err
	}

	if obj, ok := g.configs[broker]; !reflect.DeepEqual(obj.config, config) || !ok {
		err = g.setOperator(config)
		if err != nil {
			return nil, err
		}
	}
	return g.configs[channel.Spec.SelectedBroker].op, nil
}

func (g GenOp) setOperator(config metabrokers.BrokerConfiguration) error {
	var err error
	switch config.Broker() {
	case "kafka":
		kafkaConfig := config.(*sidecars.KafkaConfig)
		operator, err := kafka.NewOperator(g.memory, *kafkaConfig)
		if err == nil {
			g.configs[config.Broker()] = struct {
				config metabrokers.BrokerConfiguration
				op     ChannelOperatorInterface
			}{
				config: config,
				op:     operator,
			}
		}
	default:
		err = ierrors.NewError().Message("").Build()
	}
	return err
}

//Get executes Get method of correct operator given the desired channel's broker
func (g GenOp) Get(ctx context.Context, scope string, name string) (*meta.Channel, error) {
	op, err := g.getOperator(scope, name)
	if err != nil {
		return nil, err
	}
	return op.Get(ctx, scope, name)
}

//Create executes Create method of correct operator given the desired channel's broker
func (g GenOp) Create(ctx context.Context, scope string, channel *meta.Channel) error {
	op, err := g.getOperator(scope, channel.Meta.Name)
	if err != nil {
		return err
	}
	return op.Create(ctx, scope, channel)
}

//Update executes Update method of correct operator given the desired channel's broker
func (g GenOp) Update(ctx context.Context, scope string, channel *meta.Channel) error {
	op, err := g.getOperator(scope, channel.Meta.Name)
	if err != nil {
		return err
	}
	return op.Update(ctx, scope, channel)
}

//Delete executes Delete method of correct operator given the desired channel's broker
func (g GenOp) Delete(ctx context.Context, scope string, name string) error {
	op, err := g.getOperator(scope, name)
	if err != nil {
		return err
	}
	return op.Delete(ctx, scope, name)
}
