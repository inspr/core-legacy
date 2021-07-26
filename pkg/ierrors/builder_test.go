package ierrors

import (
	"reflect"
	"testing"
)

func TestErrCode(t *testing.T) {
	type fields struct {
		err *ierror
	}
	tests := []struct {
		name   string
		fields fields
		exec   func(e *ierror) *ierror
		want   ErrCode
	}{
		{
			name: "It should receive the error Unknown",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e
			},
			want: Unknown,
		},
		{
			name: "It should add the code NotFound to the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.NotFound()
			},
			want: NotFound,
		},
		{
			name: "It should add the code AlreadyExists the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.AlreadyExists()
			},
			want: AlreadyExists,
		},
		{
			name: "It should add the code BadRequest the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.BadRequest()
			},
			want: BadRequest,
		},
		{
			name: "It should add the code InternalServer the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InternalServer()
			},
			want: InternalServer,
		},
		{
			name: "It should add the code InvalidName the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidName()
			},
			want: InvalidName,
		},
		{
			name: "It should add the code InvalidApp the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidApp()
			},
			want: InvalidApp,
		},
		{
			name: "It should add the code InvalidChannel the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidChannel()
			},
			want: InvalidChannel,
		},
		{
			name: "It should add the code InvalidType the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidType()
			},
			want: InvalidType,
		},
		{
			name: "It should add the code InvalidFile the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidFile()
			},
			want: InvalidFile,
		},
		{
			name: "It should add the code InvalidArgs the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidArgs()
			},
			want: InvalidArgs,
		},
		{
			name: "It should add the code Forbidden the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.Forbidden()
			},
			want: Forbidden,
		},
		{
			name: "It should add the code Unauthorized the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.Unauthorized()
			},
			want: Unauthorized,
		},
		{
			name: "It should add the code ExternalErr the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.ExternalErr()
			},
			want: ExternalPkg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.exec(tt.fields.err)
			if !reflect.DeepEqual(got.code, tt.want) {
				t.Errorf("ErrCode = %v, want %v", got, tt.want)
			}
		})
	}
}
