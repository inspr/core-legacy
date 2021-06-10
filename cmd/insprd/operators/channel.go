package operators

import (
	"context"
	"reflect"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/operators/kafka"
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/brokers"
)

type GenOp struct {
	mem     memory.Manager
	configs map[string]struct {
		config brokers.BrokerConfiguration
		op     ChannelOperatorInterface
	}
}

func NewGeneralOperator(memory memory.Manager) *GenOp {
	return &GenOp{
		mem: memory,
		configs: make(map[string]struct {
			config brokers.BrokerConfiguration
			op     ChannelOperatorInterface
		}),
	}
}

func (g GenOp) getOperator(scope string, name string) ChannelOperatorInterface {
	channel, _ := g.mem.Root().Channels().Get(scope, name)
	return g.configs[channel.Spec.SelectedBroker].op
}

// func (g GenOp) SetOperator(broker string) func(config BrokerConfig, mem memory.Manager) ChannelOperatorInterface {
// 	if c, b := g.configs[broker]; !reflect.DeepEqual(c, config) {

// 		switch broker {
// 		case "kafka":
// 			return kafka.NewKafkaOperator()

// 		}
// 	} else {
// 		return b
// 	}
// }

func (g GenOp) SetOperator(config brokers.BrokerConfiguration, mem memory.Manager) error {
	var err error
	if obj, ok := g.configs[config.Broker()]; !reflect.DeepEqual(obj.config, config) || !ok {
		switch config.Broker() {
		case "kafka":
			operator, err := kafka.NewOperator(mem, config.(sidecars.KafkaConfig))
			if err == nil {
				g.configs[config.Broker()] = struct {
					config brokers.BrokerConfiguration
					op     ChannelOperatorInterface
				}{
					config: config,
					op:     operator,
				}
			}
		default:
			err = ierrors.NewError().Message("").Build()
		}
	}
	return err
}

func (g GenOp) Get(ctx context.Context, scope string, name string) (*meta.Channel, error) {
	return g.getOperator(scope, name).Get(ctx, scope, name)
}

func (g GenOp) Create(ctx context.Context, scope string, channel *meta.Channel) error {
	return g.getOperator(scope, channel.Meta.Name).Create(ctx, scope, channel)
}

func (g GenOp) Update(ctx context.Context, scope string, channel *meta.Channel) error {
	return g.getOperator(scope, channel.Meta.Name).Update(ctx, scope, channel)
}

func (g GenOp) Delete(ctx context.Context, scope string, name string) error {
	return g.getOperator(scope, name).Delete(ctx, scope, name)
}
