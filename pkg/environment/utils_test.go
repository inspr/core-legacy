package environment

import "testing"

func TestInsprEnvironment_IsInInputChannel(t *testing.T) {
	type fields struct {
		InputChannels  string
		OutputChannels string
		UnixSocketAddr string
	}
	type args struct {
		word      string
		separator string
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
				word:      "b",
				separator: ";",
			},
			want: true,
		},
		{
			name:   "channel_not_found",
			fields: defaultFields,
			args: args{
				word:      "f",
				separator: ";",
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
			if got := insprEnv.IsInInputChannel(tt.args.word, tt.args.separator); got != tt.want {
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
		word      string
		separator string
	}

	defaultFields := fields{
		InputChannels:  "a;b;c;d;e",
		OutputChannels: "1-2-3-4-5",
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
				word:      "1",
				separator: "-",
			},
			want: true,
		},
		{
			name:   "channel_not_found",
			fields: defaultFields,
			args: args{
				word:      "f",
				separator: "-",
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
			if got := insprEnv.IsInOutputChannel(tt.args.word, tt.args.separator); got != tt.want {
				t.Errorf("InsprEnvironment.IsInOutputChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
