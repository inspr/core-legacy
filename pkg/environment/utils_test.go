package environment

import (
	"reflect"
	"testing"
)

func TestInsprEnvironment_IsInInputChannel(t *testing.T) {
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
			insprEnv := &InsprEnvironmentVariables{
				InputChannels:  tt.fields.InputChannels,
				OutputChannels: tt.fields.OutputChannels,
				UnixSocketAddr: tt.fields.UnixSocketAddr,
			}
			if got := insprEnv.IsInInputChannel(tt.args.channel); got != tt.want {
				t.Errorf("InsprEnvironment.IsInInputChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprEnvironment_IsInOutputChannel(t *testing.T) {
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
				channel: "1",
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
			insprEnv := &InsprEnvironmentVariables{
				InputChannels:  tt.fields.InputChannels,
				OutputChannels: tt.fields.OutputChannels,
				UnixSocketAddr: tt.fields.UnixSocketAddr,
			}
			if got := insprEnv.IsInOutputChannel(tt.args.channel); got != tt.want {
				t.Errorf("InsprEnvironment.IsInOutputChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprEnvironmentVariables_GetInputChannelList(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insprEnv := &InsprEnvironmentVariables{
				InputChannels:    tt.fields.InputChannels,
				OutputChannels:   tt.fields.OutputChannels,
				UnixSocketAddr:   tt.fields.UnixSocketAddr,
				InsprAppContext:  tt.fields.InsprAppContext,
				InsprEnvironment: tt.fields.InsprEnvironment,
			}
			if got := insprEnv.GetInputChannelList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprEnvironmentVariables.GetInputChannelList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprEnvironmentVariables_GetOutputChannelList(t *testing.T) {
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
				InputChannels:    "",
				OutputChannels:   "ch1;ch2;ch3;ch4",
				UnixSocketAddr:   "",
				InsprAppContext:  "",
				InsprEnvironment: "",
			},
			want: []string{"ch1", "ch2", "ch3", "ch4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insprEnv := &InsprEnvironmentVariables{
				InputChannels:    tt.fields.InputChannels,
				OutputChannels:   tt.fields.OutputChannels,
				UnixSocketAddr:   tt.fields.UnixSocketAddr,
				InsprAppContext:  tt.fields.InsprAppContext,
				InsprEnvironment: tt.fields.InsprEnvironment,
			}
			if got := insprEnv.GetOutputChannelList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprEnvironmentVariables.GetOutputChannelList() = %v, want %v", got, tt.want)
			}
		})
	}
}
