package utils

import (
	"testing"
)

func TestGetFreePorts(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name         string
		args         args
		wantedLength int
		wantErr      bool
	}{
		{
			name: "return X valid ports",
			args: args{
				count: 4,
			},
			wantedLength: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFreePorts(tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFreePorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantedLength {
				t.Errorf("GetFreePorts() length = %v, want %v", got, tt.wantedLength)
			}
		})
	}
}
