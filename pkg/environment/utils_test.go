package environment

import (
	"os"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/utils"
)

func TestIsInChannelBoundary(t *testing.T) {
	type fields struct {
		InputChannels  []brokers.ChannelBroker
		OutputChannels []brokers.ChannelBroker
		UnixSocketAddr string
	}
	type args struct {
		channel string
	}
	defaultFields := fields{
		InputChannels: []brokers.ChannelBroker{
			{ChName: "a"},
			{ChName: "b"},
			{ChName: "c"},
			{ChName: "d"},
			{ChName: "e"},
		},
		OutputChannels: []brokers.ChannelBroker{
			{ChName: "1"},
			{ChName: "2"},
			{ChName: "3"},
			{ChName: "4"},
			{ChName: "5"},
		},
		UnixSocketAddr: "socket",
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "channel_found",
			fields: defaultFields,
			args: args{
				channel: "b",
			},
			want: true,
		},
		{
			name:   "channel_not_found",
			fields: defaultFields,
			args: args{
				channel: "f",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInChannelBoundary(tt.args.channel, tt.fields.InputChannels); got != tt.want {
				t.Errorf("InsprEnvironment.IsInInputChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChannelBoundaryList(t *testing.T) {
	type fields struct {
		InputChannels    string
		OutputChannels   string
		UnixSocketAddr   string
		InsprAppContext  string
		InsprEnvironment string
	}
	tests := []struct {
		name   string
		fields fields
		want   utils.StringArray
	}{
		{
			name: "It should get all the channels in the InputChannels env",
			fields: fields{
				InputChannels:    "ch1_b;ch2_b;ch3_b;ch4_b",
				OutputChannels:   "",
				UnixSocketAddr:   "",
				InsprAppContext:  "",
				InsprEnvironment: "",
			},
			want: utils.StringArray{"ch1", "ch2", "ch3", "ch4"},
		},
		{
			name: "Returns empty string slice",
			fields: fields{
				InputChannels:    "",
				OutputChannels:   "",
				UnixSocketAddr:   "",
				InsprAppContext:  "",
				InsprEnvironment: "",
			},
			want: utils.StringArray{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChannelBoundaryList(getChannelData(tt.fields.InputChannels)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprEnvVars.GetInputChannelList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSchema(t *testing.T) {
	os.Setenv("ch1_SCHEMA", "this_is_a_schema")
	defer os.Unsetenv("ch1_SCHEMA")
	type fields struct {
		InputChannels    string
		OutputChannels   string
		UnixSocketAddr   string
		InsprAppContext  string
		InsprEnvironment string
	}
	tests := []struct {
		name    string
		fields  fields
		channel string
		wantErr bool
		want    string
	}{
		{
			name: "Get valid env var with schema",
			fields: fields{
				InputChannels:    "",
				OutputChannels:   "ch1;ch2;ch3;ch4",
				UnixSocketAddr:   "",
				InsprAppContext:  "",
				InsprEnvironment: "",
			},
			channel: "ch1",
			wantErr: false,
			want:    "this_is_a_schema",
		},
		{
			name: "Get invalid env var with schema",
			fields: fields{
				InputChannels:    "",
				OutputChannels:   "ch1;ch2;ch3;ch4",
				UnixSocketAddr:   "",
				InsprAppContext:  "",
				InsprEnvironment: "",
			},
			channel: "ch5",
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSchema(tt.channel)
			if tt.wantErr && err == nil {
				t.Errorf("InsprEnvVars.GetSchema() = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprEnvVars.GetSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResolvedBoundaryChannelList(t *testing.T) {
	os.Setenv("ch1_RESOLVED", "channel1_resolved")
	os.Setenv("ch2_RESOLVED", "channel2_resolved")
	defer os.Unsetenv("ch1_RESOLVED")
	defer os.Unsetenv("ch2_RESOLVED")

	type args struct {
		boundary string
	}
	tests := []struct {
		name string
		args args
		want utils.StringArray
	}{
		{
			name: "Returns resolved boundary",
			args: args{
				boundary: "ch1_b;ch2_b",
			},
			want: utils.StringArray{"channel1_resolved", "channel2_resolved"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetResolvedBoundaryChannelList(getChannelData(tt.args.boundary)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResolvedBoundaryChannelList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResolvedChannel(t *testing.T) {
	os.Setenv("ch1_RESOLVED", "channel1_resolved")
	defer os.Unsetenv("ch1_RESOLVED")
	type args struct {
		channel    string
		inputChan  string
		outputChan string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Returns resolved channel",
			args: args{
				channel:   "ch1",
				inputChan: "ch1_b;ch2_b",
			},
			wantErr: false,
			want:    "channel1_resolved",
		},
		{
			name: "Invalid channel",
			args: args{
				channel:   "ch3",
				inputChan: "ch1_b;ch2_b",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResolvedChannel(tt.args.channel, getChannelData(tt.args.inputChan), getChannelData(tt.args.outputChan))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResolvedChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetResolvedChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
