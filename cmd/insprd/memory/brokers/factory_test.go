package brokers

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/sidecars/models"
)

func TestAbstractBrokerFactory_Subscribe(t *testing.T) {
	type args struct {
		broker  string
		factory models.SidecarFactory
	}
	tests := []struct {
		name    string
		exec    func(abf *AbstractBrokerFactory)
		args    args
		wantErr bool
	}{
		{
			name: "valid subscription from nil singleton",
			exec: func(abf *AbstractBrokerFactory) {

			},
			args: args{
				broker:  "broker_name",
				factory: nil,
			},
			wantErr: false,
		},
		{
			name: "valid subscription from instanced singleton",
			exec: func(abf *AbstractBrokerFactory) {
				abf.Subscribe("broker-name", nil)
			},
			args: args{
				broker:  "broker_name",
				factory: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid subscription from singleton containing broker",
			exec: func(abf *AbstractBrokerFactory) {
				abf.Subscribe("broker_name", nil)
			},
			args: args{
				broker:  "broker_name",
				factory: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFactories()
			abf := &AbstractBrokerFactory{}
			tt.exec(abf)
			if err := abf.Subscribe(tt.args.broker, tt.args.factory); (err != nil) != tt.wantErr {
				t.Errorf("AbstractBrokerFactory.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func resetFactories() {
	factories = nil
}

func TestAbstractBrokerFactory_Get(t *testing.T) {
	type args struct {
		broker string
	}
	tests := []struct {
		name    string
		args    args
		exec    func(abf *AbstractBrokerFactory)
		want    models.SidecarFactory
		wantErr bool
	}{
		{
			name: "invalid get from nil singleton",
			args: args{
				broker: "broner_name",
			},
			exec: func(abf *AbstractBrokerFactory) {

			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid get from invalid broker on singleton",
			args: args{
				broker: "broker_name",
			},
			exec: func(abf *AbstractBrokerFactory) {
				abf.Subscribe("broker-name", nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid get from nil broker subscribed on singleton",
			args: args{
				broker: "broker_name",
			},
			exec: func(abf *AbstractBrokerFactory) {
				abf.Subscribe("broker_name", nil)
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFactories()
			abf := &AbstractBrokerFactory{}
			tt.exec(abf)
			got, err := abf.Get(tt.args.broker)
			if (err != nil) != tt.wantErr {
				t.Errorf("AbstractBrokerFactory.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&got, &tt.want) {
				t.Errorf("AbstractBrokerFactory.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
