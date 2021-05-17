package brokers

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
)

func TestMemoryManager_Brokers(t *testing.T) {

	tests := []struct {
		name string
		want BrokerInterface
	}{
		{
			name: "standard Brokers() behaviour",
			want: &BrokerMemoryManager{
				factory: &AbstractBrokerFactory{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bm := &BrokerManager{}
			if got := bm.Brokers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryManager.Brokers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_GetAll(t *testing.T) {
	tests := []struct {
		name string
		bmm  *BrokerMemoryManager
		want brokers.BrokerStatusArray
	}{
		{
			name: "getall from empty brokerMM",
			bmm:  &BrokerMemoryManager{},
			want: brokers.BrokerStatusArray{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := &BrokerMemoryManager{}
			if got := bmm.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrokersMemoryManager.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_GetDefault(t *testing.T) {
	tests := []struct {
		name string
		bmm  *BrokerMemoryManager
		want brokers.BrokerStatus
	}{
		{
			name: "getall from empty brokerMM",
			bmm:  &BrokerMemoryManager{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := &BrokerMemoryManager{}
			if got := bmm.GetDefault(); got != tt.want {
				t.Errorf("BrokersMemoryManager.GetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_get(t *testing.T) {
	tests := []struct {
		name string
		bmm  *BrokerMemoryManager
		exec func(bmm *BrokerMemoryManager)
		want *brokers.Brokers
	}{
		{
			name: "getall from empty brokerMM",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm *BrokerMemoryManager) {

			},
			want: &brokers.Brokers{
				Available: make(metautils.StrSet),
			},
		},
		{
			name: "getall from instanciated singleton",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm *BrokerMemoryManager) {
				bmm.Create("brk1", nil)
				bmm.Create("brk2", nil)
				bmm.Create("brk3", nil)
				bmm.SetDefault("brk1")
			},
			want: &brokers.Brokers{
				Available: metautils.StrSet{
					"brk1": true,
					"brk2": true,
					"brk3": true,
				},
				Default: "brk1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := &BrokerMemoryManager{}
			tt.exec(bmm)
			if got := bmm.get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrokersMemoryManager.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_Create(t *testing.T) {
	type args struct {
		broker brokers.BrokerStatus
		config brokers.BrokerConfiguration
	}
	tests := []struct {
		name    string
		bmm     *BrokerMemoryManager
		args    args
		exec    func(bmm *BrokerMemoryManager)
		wantErr bool
	}{
		{
			name: "valid create",
			args: args{
				broker: "brk1",
				config: nil,
			},
			exec: func(bmm *BrokerMemoryManager) {

			},
			bmm:     &BrokerMemoryManager{},
			wantErr: false,
		},
		{
			name: "invalid create",
			args: args{
				broker: "brk1",
				config: nil,
			},
			bmm: &BrokerMemoryManager{},
			exec: func(bmm *BrokerMemoryManager) {
				bmm.Create("brk1", nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := &BrokerMemoryManager{}
			tt.exec(bmm)
			if err := bmm.Create(tt.args.broker, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBrokersMemoryManager_SetDefault(t *testing.T) {
	type args struct {
		broker brokers.BrokerStatus
	}
	tests := []struct {
		name    string
		bmm     *BrokerMemoryManager
		args    args
		exec    func(bmm *BrokerMemoryManager)
		wantErr bool
	}{
		{
			name: "invalid setdefault",
			args: args{
				broker: "brk1",
			},
			exec: func(bmm *BrokerMemoryManager) {
			},
			bmm:     &BrokerMemoryManager{},
			wantErr: true,
		},
		{
			name: "valid setdefault",
			args: args{
				broker: "brk1",
			},
			exec: func(bmm *BrokerMemoryManager) {
				bmm.Create("brk1", nil)
			},
			bmm:     &BrokerMemoryManager{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := &BrokerMemoryManager{}
			tt.exec(bmm)
			if err := bmm.SetDefault(tt.args.broker); (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.SetDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func resetBrokers() {
	broker = nil
}
