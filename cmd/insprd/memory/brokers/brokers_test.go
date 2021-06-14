package brokers

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/utils"
)

var kafkaStructMock = sidecars.KafkaConfig{
	BootstrapServers: "",
	AutoOffsetReset:  "",
	KafkaInsprAddr:   "",
	SidecarImage:     "",
}

func TestBrokersMemoryManager_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		want    utils.StringArray
		wantErr bool
	}{
		{
			name:    "getall from empty brokerMM",
			want:    utils.StringArray{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := GetBrokerMemory()
			got, err := bmm.GetAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrokersMemoryManager.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_GetDefault(t *testing.T) {
	tests := []struct {
		name    string
		bmm     *BrokerMemoryManager
		want    string
		wantErr bool
	}{
		{
			name:    "getdefault from empty brokerMM",
			bmm:     &BrokerMemoryManager{},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			bmm := GetBrokerMemory()
			got, err := bmm.GetDefault()
			if (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.GetDefault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BrokersMemoryManager.GetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_get(t *testing.T) {
	tests := []struct {
		name    string
		bmm     *BrokerMemoryManager
		want    *brokers.Brokers
		wantErr bool
	}{
		{
			name: "get from instanciated singleton",
			bmm: &BrokerMemoryManager{
				broker: &brokers.Brokers{
					Available: brokers.BrokerStatusArray{
						"brk1": nil,
						"brk2": nil,
						"brk3": nil,
					},
					Default: "brk1",
				},
			},
			want: &brokers.Brokers{
				Available: brokers.BrokerStatusArray{
					"brk1": nil,
					"brk2": nil,
					"brk3": nil,
				},
				Default: "brk1",
			},
			wantErr: false,
		},
		{
			name: "get from nil singleton memory",
			bmm: &BrokerMemoryManager{
				broker: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBrokers()
			got, err := tt.bmm.get()
			if (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrokersMemoryManager.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockedConfigs struct {
	broker string
}

func (mb mockedConfigs) Broker() string {
	return mb.broker
}

func TestBrokersMemoryManager_Create_and_SetDefault(t *testing.T) {
	resetBrokers()

	tests := []struct {
		name    string
		bmm     *BrokerMemoryManager
		exec    func(bmm Manager) error
		wantErr bool
	}{
		{
			name: "invalid create - broker not supported",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm Manager) error {
				return bmm.Create(mockedConfigs{broker: "brk1"})
			},
			wantErr: true,
		},
		{
			name: "valid create",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm Manager) error {
				return bmm.Create(&kafkaStructMock)
			},
		},
		{
			name: "invalid create - broker already exists",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm Manager) error {
				return bmm.Create(&kafkaStructMock)
			},
			wantErr: true,
		},
		{
			name: "invalid setdefault",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm Manager) error {
				return bmm.SetDefault("brk1")
			},
			wantErr: true,
		},
		{
			name: "valid setdefault",
			bmm:  &BrokerMemoryManager{},
			exec: func(bmm Manager) error {
				return bmm.SetDefault(brokers.Kafka)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bmm := GetBrokerMemory()
			if tt.exec != nil {
				if err := tt.exec(bmm); (err != nil) != tt.wantErr {
					t.Errorf("BrokersMemoryManager method error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func resetBrokers() {
	brokerMemory = nil
}
