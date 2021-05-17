package brokers

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
)

func TestBrokersMemoryManager_GetAll(t *testing.T) {
	tests := []struct {
		name string
		want brokers.BrokerStatusArray
	}{
		{
			name: "getall from empty brokerMM",
			want: brokers.BrokerStatusArray{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := GetBrokerMemory()
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
			bmm := GetBrokerMemory()
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
		want *brokers.Brokers
	}{
		{
			name: "getall from instanciated singleton",
			bmm: &BrokerMemoryManager{
				broker: &brokers.Brokers{
					Available: metautils.StrSet{
						"brk1": true,
						"brk2": true,
						"brk3": true,
					},
					Default: "brk1",
				},
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
			if got := tt.bmm.get(); !reflect.DeepEqual(got, tt.want) {
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
		exec    func(bmm Manager)
		wantErr bool
	}{
		{
			name: "valid create",
			args: args{
				broker: "brk1",
				config: nil,
			},
			exec: func(bmm Manager) {

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
			exec: func(bmm Manager) {
				bmm.Create("brk1", nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := GetBrokerMemory()
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
		exec    func(bmm Manager)
		wantErr bool
	}{
		{
			name: "invalid setdefault",
			args: args{
				broker: "brk1",
			},
			exec: func(bmm Manager) {
			},
			bmm:     &BrokerMemoryManager{},
			wantErr: true,
		},
		{
			name: "valid setdefault",
			args: args{
				broker: "brk1",
			},
			exec: func(bmm Manager) {
				bmm.Create("brk1", nil)
			},
			bmm:     &BrokerMemoryManager{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := GetBrokerMemory()
			tt.exec(bmm)
			if err := bmm.SetDefault(tt.args.broker); (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.SetDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func resetBrokers() {
	brokerMemory = nil
}
