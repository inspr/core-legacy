package brokers

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
)

func TestGetBrokerMemory(t *testing.T) {
	tests := []struct {
		name string
		want Manager
		exec func()
	}{
		{
			name: "GetManager from nil pointer",
			want: &BrokerMemoryManager{
				broker: &brokers.Brokers{
					Available: make(metautils.StrSet),
				},
				factory: &AbstractBrokerFactory{},
			},
			exec: func() {
				brokerMemory = nil
			},
		},
		{
			name: "GetManager from intanced pointer",
			want: &BrokerMemoryManager{
				broker: &brokers.Brokers{
					Available: make(metautils.StrSet),
				},
				factory: &AbstractBrokerFactory{},
			},
			exec: func() {
				GetBrokerMemory()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.exec()
			if got := GetBrokerMemory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBrokerMemory() = %v, want %v", got, tt.want)
			}
		})
	}
}
