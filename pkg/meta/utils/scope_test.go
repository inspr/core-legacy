package utils

import "testing"

func TestIsValidScope(t *testing.T) {
	type args struct {
		scope string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Invalid scope - it should return false",
			args: args{
				scope: "app1.app2.",
			},
			want: false,
		},
		{
			name: "Invalid scope - it should return false",
			args: args{
				scope: ".app1.app2",
			},
			want: false,
		},
		{
			name: "Invalid scope - it should return false",
			args: args{
				scope: "..app1.app2",
			},
			want: false,
		},
		{
			name: "Invalid scope - it should return false",
			args: args{
				scope: "app1..app2",
			},
			want: false,
		},
		{
			name: "Valid scope - empty scope",
			args: args{
				scope: "",
			},
			want: true,
		},
		{
			name: "Valid scope",
			args: args{
				scope: "app1.app2",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidScope(tt.args.scope); got != tt.want {
				t.Errorf("IsValidScope() = %v, want %v", got, tt.want)
			}
		})
	}
}
