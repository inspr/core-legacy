package environment

import (
	"os"
	"reflect"
	"testing"
)

func TestInsprEnvironment_IsInChannelBoundary(t *testing.T) {
	type fields struct {
		InputChannels  string
		OutputChannels string
		UnixSocketAddr string
	}
	type args struct {
		channel string
	}
	defaultFields := fields{
		InputChannels:  "a;b;c;d;e",
		OutputChannels: "1;2;3;4;5",
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

func TestInsprEnvVars_GetChannelBoundaryList(t *testing.T) {
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
		want   []string
	}{
		{
			name: "It should get all the channels in the InputChannels env",
			fields: fields{
				InputChannels:    "ch1;ch2;ch3;ch4",
				OutputChannels:   "",
				UnixSocketAddr:   "",
				InsprAppContext:  "",
				InsprEnvironment: "",
			},
			want: []string{"ch1", "ch2", "ch3", "ch4"},
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
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChannelBoundaryList(tt.fields.InputChannels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprEnvVars.GetInputChannelList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprEnvVars_GetSchema(t *testing.T) {
	os.Setenv("ch1_SCHEMA", "this_is_a_schema")
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
			got, err := GetSchema(tt.channel, tt.fields.InputChannels, tt.fields.OutputChannels)
			if tt.wantErr && err == nil {
				t.Errorf("InsprEnvVars.GetSchema() = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprEnvVars.GetSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
