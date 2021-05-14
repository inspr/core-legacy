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
			want: &BrokersMemoryManager{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSingleton()
			bm := &BrokersManager{}
			if got := bm.Brokers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryManager.Brokers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_GetAll(t *testing.T) {
	tests := []struct {
		name string
		bmm  *BrokersMemoryManager
		want brokers.BrokerStatusArray
	}{
		{
			name: "getall from empty brokerMM",
			bmm:  &BrokersMemoryManager{},
			want: brokers.BrokerStatusArray{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSingleton()
			bmm := &BrokersMemoryManager{}
			if got := bmm.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrokersMemoryManager.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_GetDefault(t *testing.T) {
	tests := []struct {
		name string
		bmm  *BrokersMemoryManager
		want brokers.BrokerStatus
	}{
		{
			name: "getall from empty brokerMM",
			bmm:  &BrokersMemoryManager{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSingleton()
			bmm := &BrokersMemoryManager{}
			if got := bmm.GetDefault(); got != tt.want {
				t.Errorf("BrokersMemoryManager.GetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersMemoryManager_get(t *testing.T) {
	tests := []struct {
		name string
		bmm  *BrokersMemoryManager
		exec func(bmm *BrokersMemoryManager)
		want *brokers.Brokers
	}{
		{
			name: "getall from empty brokerMM",
			bmm:  &BrokersMemoryManager{},
			exec: func(bmm *BrokersMemoryManager) {

			},
			want: &brokers.Brokers{
				Available: make(metautils.StrSet),
			},
		},
		{
			name: "getall from instanciated singleton",
			bmm:  &BrokersMemoryManager{},
			exec: func(bmm *BrokersMemoryManager) {
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
			resetSingleton()
			bmm := &BrokersMemoryManager{}
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
		bmm     *BrokersMemoryManager
		args    args
		exec    func(bmm *BrokersMemoryManager)
		wantErr bool
	}{
		{
			name: "valid create",
			args: args{
				broker: "brk1",
				config: nil,
			},
			exec: func(bmm *BrokersMemoryManager) {

			},
			bmm:     &BrokersMemoryManager{},
			wantErr: false,
		},
		{
			name: "invalid create",
			args: args{
				broker: "brk1",
				config: nil,
			},
			bmm: &BrokersMemoryManager{},
			exec: func(bmm *BrokersMemoryManager) {
				bmm.Create("brk1", nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSingleton()
			bmm := &BrokersMemoryManager{}
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
		bmm     *BrokersMemoryManager
		args    args
		exec    func(bmm *BrokersMemoryManager)
		wantErr bool
	}{
		{
			name: "invalid setdefault",
			args: args{
				broker: "brk1",
			},
			exec: func(bmm *BrokersMemoryManager) {
			},
			bmm:     &BrokersMemoryManager{},
			wantErr: true,
		},
		{
			name: "valid setdefault",
			args: args{
				broker: "brk1",
			},
			exec: func(bmm *BrokersMemoryManager) {
				bmm.Create("brk1", nil)
			},
			bmm:     &BrokersMemoryManager{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSingleton()
			bmm := &BrokersMemoryManager{}
			tt.exec(bmm)
			if err := bmm.SetDefault(tt.args.broker); (err != nil) != tt.wantErr {
				t.Errorf("BrokersMemoryManager.SetDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func resetSingleton() {
	broker = &brokers.Brokers{
		Available: make(metautils.StrSet),
	}
}
