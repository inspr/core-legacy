package utils

import "testing"

func TestGetScope(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "It should return the scope",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetScope()
			if (err != nil) != tt.wantErr {
				t.Errorf("getScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getScope() = %v, want %v", got, tt.want)
			}
		})
	}
}
