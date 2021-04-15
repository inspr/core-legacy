package ierrors

import (
	"errors"
	"testing"
)

func TestHasCode(t *testing.T) {
	type args struct {
		target error
		code   InsprErrorCode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "non_ierror_test",
			args: args{
				target: errors.New("mock"),
				code:   0,
			},
			want: false,
		},
		{
			name: "ierror_test",
			args: args{
				target: NewError().AlreadyExists().Build(),
				code:   AlreadyExists,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasCode(tt.args.target, tt.args.code); got != tt.want {
				t.Errorf("HasCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIerror(t *testing.T) {
	type args struct {
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "non_ierror_test",
			args: args{target: errors.New("mock")},
			want: false,
		},
		{
			name: "ierror_test",
			args: args{target: NewError().AlreadyExists().Build()},
			want: true,
		},
		{
			name: "ierror_test_empty_errCode",
			args: args{target: &InsprError{Code: 0}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIerror(tt.args.target); got != tt.want {
				t.Errorf("IsIerror() = %v, want %v", got, tt.want)
			}
		})
	}
}
