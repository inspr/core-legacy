package ierrors

import (
	"errors"
	"io/fs"
	"testing"
)

func TestHasCode(t *testing.T) {
	type args struct {
		target error
		code   ErrCode
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
				target: New("").AlreadyExists(),
				code:   AlreadyExists,
			},
			want: true,
		},
		{
			name: "ierror_test_with_different_code",
			args: args{
				target: New(""),
				code:   AlreadyExists,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasCode(tt.args.target, tt.args.code)
			if got != tt.want {
				t.Errorf("HasCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIerror_Is(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "two_equal_errors_code",
			fields: fields{
				err: New("new message").BadRequest(),
			},
			args: args{
				target: New("another message").BadRequest(),
			},
			want: true,
		},
		{
			name: "two_different_errors_code",
			fields: fields{
				err: New("new message").InternalServer(),
			},
			args: args{
				target: New("another message").BadRequest(),
			},
			want: false,
		},
		{
			name: "external_error_not_in_stack",
			fields: fields{
				err: Wrap(
					New(fs.ErrClosed).ExternalErr(),
					"ctx_1",
					"ctx_2",
					"ctx_3",
				),
			},
			args: args{
				target: fs.SkipDir,
			},
			want: false,
		},
		{
			name: "external_error_in_stack",
			fields: fields{
				err: Wrap(
					Wrap(
						Wrap(
							fs.ErrClosed,
							"ctx_1",
						),
						"ctx_2",
					),
					"ctx_3",
				),
			},
			args: args{
				target: fs.ErrClosed,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Is(tt.fields.err, tt.args.target)
			if got != tt.want {
				t.Errorf("ierror.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
