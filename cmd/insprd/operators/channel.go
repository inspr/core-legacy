package operators

import (
	"context"
	"fmt"
	"reflect"

	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	kafkaop "inspr.dev/inspr/cmd/insprd/operators/kafka"
	"inspr.dev/inspr/cmd/sidecars"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metabrokers "inspr.dev/inspr/pkg/meta/brokers"
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

func (g GenOp) getOperator(scope, name string, deleteCmd bool) (ChannelOperatorInterface, error) {
	var channel *meta.Channel
	var err error

	// delete should get it's channel information from the unaltered tree,
	// otherwise the channel wouldnt be found
	if deleteCmd {
		channel, err = g.memory.Tree().Channels().Get(scope, name)
		if err != nil {
			return nil, err
		}
	} else {
		channel, err = g.memory.Channels().Get(scope, name)
		if err != nil {
			return nil, err
		}
	}

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
		operator, err := kafkaop.NewOperator(g.memory, *kafkaConfig)
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
func (g GenOp) Get(ctx context.Context, scope, name string) (*meta.Channel, error) {
	op, err := g.getOperator(scope, name, false)
	if err != nil {
		return nil, err
	}
	return op.Get(ctx, scope, name)
}

//Create executes Create method of correct operator given the desired channel's broker
func (g GenOp) Create(ctx context.Context, scope string, channel *meta.Channel) error {
	op, err := g.getOperator(scope, channel.Meta.Name, false)
	if err != nil {
		return err
	}
	return op.Create(ctx, scope, channel)
}

//Update executes Update method of correct operator given the desired channel's broker
func (g GenOp) Update(ctx context.Context, scope string, channel *meta.Channel) error {
	op, err := g.getOperator(scope, channel.Meta.Name, false)
	if err != nil {
		return err
	}
	return op.Update(ctx, scope, channel)
}

//Delete executes Delete method of correct operator given the desired channel's broker
func (g GenOp) Delete(ctx context.Context, scope, name string) error {
	fmt.Println("ABLUBLEBLEBLABLUBLABLE")
	op, err := g.getOperator(scope, name, true)
	if err != nil {
		return err
	}
	fmt.Println("XABLABLEXABLWXABLU")
	return op.Delete(ctx, scope, name)
}
