package ierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	tests := []struct {
		name string
		arg  interface{}
		want *ierror
	}{
		{
			name: "ierror_from_string",
			arg:  "mock",
			want: &ierror{
				err:  errors.New("mock"),
				code: Unknown,
			},
		}, {
			name: "ierror_from_standard_error",
			arg:  errors.New("mock"),
			want: &ierror{
				err:  errors.New("mock"),
				code: Unknown,
			},
		}, {
			name: "ierror_from_ierror",
			arg: &ierror{
				err:  errors.New("mock"),
				code: Unknown,
			},
			want: &ierror{
				err:  errors.New("mock"),
				code: Unknown,
			},
		}, {
			name: "invalid_interface",
			arg:  10,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.arg)

			if got == nil && got != tt.want {
				t.Errorf(
					"ierrors.New() got '%v', expected '%v'",
					got, tt.want,
				)
			}

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf(
					"ierrors.New() got '%v', expected '%v'",
					got.Error(), tt.want.Error(),
				)
			}
		})
	}
}

func TestNewIerror(t *testing.T) {
	type fields struct {
		format string
		values string
	}
	tests := []struct {
		name   string
		fields fields
		want   ierror
	}{
		{
			name: "simple_creation",
			fields: fields{
				format: "something",
			},
			want: ierror{err: errors.New("something"), code: Unknown},
		}, {
			name: "creation_with_values",
			fields: fields{
				format: "%v---",
				values: "value",
			},
			want: ierror{err: errors.New("value---"), code: Unknown},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *ierror
			if tt.fields.values == "" {
				got = newIerror(tt.fields.format)
			} else {
				got = newIerror(tt.fields.format, tt.fields.values)
			}

			if !reflect.DeepEqual(got.err.Error(), tt.want.err.Error()) {
				t.Errorf("Expected '%v', got '%v'", tt.want.err, got.err)
			}
			if !reflect.DeepEqual(got.code, tt.want.code) {
				t.Errorf("Expected %v, got %v", tt.want.code, got.code)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name:   "test_from_errors",
			fields: fields{err: errors.New("mock_err")},
			want: &ierror{
				err:  errors.New("mock_err"),
				code: Unknown,
			},
		},
		{
			name:   "test_from_ierrors",
			fields: fields{err: New("mock_err").InternalServer()},
			want: &ierror{
				err:  errors.New("mock_err"),
				code: InternalServer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := from(tt.fields.err)
			if got.Error() != tt.want.Error() {
				t.Errorf("Expected %v, got %v", got, tt.want)
			}
		})
	}
}

func TestIerror_Error(t *testing.T) {
	type fields struct {
		err *ierror
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "testing the error message of inspr error",
			fields: fields{err: New("mock_message")},
			want:   New("mock_message").Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.fields.err.Error(); got != tt.want {
				t.Errorf("ierror.Error() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "converting_nil",
			fields: fields{
				err: Wrap(
					New("mock_err"),
					"wrapper_1",
					"wrapper_2",
					"wrapper_3",
				),
			},
			want: "error : wrapper_3\n\twrapper_2\n\twrapper_1\n\tmock_err\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatError(tt.fields.err)
			if got != tt.want {
				t.Errorf(
					"FormatError got '%v', want '%v'",
					got, tt.want,
				)
			}
		})
	}
}

func TestCode(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   ErrCode
	}{
		{
			name:   "test_from_errors",
			fields: fields{err: errors.New("mock_err")},
			want:   Unknown,
		},
		{
			name:   "test_from_ierrors",
			fields: fields{err: New("mock_err").InternalServer()},
			want:   InternalServer,
		},
	}
	for _, tt := range tests {
		got := Code(tt.fields.err)
		if got != tt.want {
			t.Errorf("Expected %v, got %v", got, tt.want)
		}
	}
}

func TestIerror_Wrap(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantNil bool
	}{
		{
			name: "wrap_no_error",
			fields: fields{
				err: nil,
			},
			want:    "",
			wantNil: true,
		},
		{
			name: "wrap_standard_error_no_message",
			fields: fields{
				err: errors.New("mock_err"),
			},
			want: New("mock_err").Error(),
		},
		{
			name: "wrap_standard_error_with_message",
			fields: fields{
				err: errors.New("mock_err"),
			},
			args: args{
				msg: "wrapper_context",
			},
			want: fmt.Errorf(
				"error : %v : %w", "wrapper_context", errors.New("mock_err"),
			).Error(),
		},
		{
			name: "wrap_ierror_no_message",
			fields: fields{
				err: New("mock_err"),
			},
			want: New("mock_err").Error(),
		},
		{
			name: "wrap_ierror_with_message",
			fields: fields{
				err: New("mock_err"),
			},
			args: args{
				msg: "wrapper_context",
			},
			want: fmt.Errorf(
				"error : %v : %w", "wrapper_context", errors.New("mock_err"),
			).Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.fields.err, tt.args.msg)

			if (got == nil) != tt.wantNil {
				t.Errorf("expected nil, got %v", got)
			}
			if got != nil && got.Error() != tt.want {
				t.Errorf("ierror.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "unwrap_empty_err_stack",
			args: args{err: nil},
			want: nil,
		},
		{
			name: "unwrap_err_without_previous_wrap",
			args: args{err: errors.New("mock")},
			want: nil,
		},
		{
			name: "unwrap_err_with_previous_wrap",
			args: args{err: fmt.Errorf("wrap : %w", errors.New("mock"))},
			want: newIerror("mock"),
		},
		{
			name: "unwrap_ierror_without_previous_wrap",
			args: args{err: New("mock")},
			want: nil,
		},
		{
			name: "unwrap_ierror_with_previous_wrap",
			args: args{err: Wrap(New("mock"), "simple_wrap")},
			want: newIerror("mock"),
		},
		{
			name: "unwrap_ierror_with_previous_formatted_wrap",
			args: args{err: Wrap(New("mock"), "simple_wrap")},
			want: New("mock"),
		},
		{
			name: "unwrap_ierror_with_multiple_previous_wraps",
			args: args{
				err: Wrap(
					New("mock"),
					"first_context",
					"second_context",
					"third context",
				),
			},
			want: Wrap(
				New("mock"),
				"first_context",
				"second_context",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unwrap(tt.args.err)

			if got == nil && got != tt.want {
				t.Errorf("Expected %v, received %v", tt.want, got)
			}

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Unwrap() error = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIerror_MarshalJSON(t *testing.T) {

	expetedIerrorBytes := func(err *ierror) []byte {
		ps := parseStruct{
			Stack: err.err.Error(),
			Code:  err.code,
		}
		data, _ := json.Marshal(ps)
		return data
	}

	type fields struct {
		Err *ierror
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "simple_marshall",
			fields: fields{
				Err: New("mock_err"),
			},
			want: expetedIerrorBytes(New("mock_err")),
		},
		{
			// testing for the inner error being nil
			name: "inner_error_nil",
			fields: fields{
				Err: &ierror{err: nil},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := json.Marshal(tt.fields.Err)

			// error on marshal
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ierror.MarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			// comparing byte slices produced by parseStruct
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"ierror.MarshalJSON() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIerror_UnmarshalJSON(t *testing.T) {

	generateIerrorBytes := func(err *ierror) []byte {
		data, _ := json.Marshal(err)
		return data
	}

	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "unmarshal_simple_ierror",
			args: args{data: generateIerrorBytes(New("mock_err"))},
			want: New("mock_err"),
		},
		{
			name: "unmarshal_wrapped_error",
			args: args{data: generateIerrorBytes(
				from(
					Wrap(
						New("mock_err"),
						"mock_context",
					),
				),
			)},
			want: Wrap(
				New("mock_err"),
				"mock_context",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierr := New("")
			err := json.Unmarshal(tt.args.data, &ierr)

			if (err != nil) && !errors.Is(ierr, tt.want) {
				t.Errorf(
					"json.Unmarshal(ierror) got = %v, wanted %v",
					err,
					tt.want,
				)
			}
		})
	}
}

func TestIerror_stackError(t *testing.T) {
	type fields struct {
		stack string
		code  ErrCode
	}
	tests := []struct {
		name   string
		fields fields
		wanted error
	}{
		{
			name: "basic_test",
			fields: fields{
				stack: Wrap(
					errors.New("mock_err"),
					"wrap_1",
				).Error(),
				code: Unknown,
			},

			wanted: Wrap(
				errors.New("mock_err"),
				"wrap_1",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stackToError(tt.fields.stack, Unknown)

			if got.Error() != tt.wanted.Error() {
				t.Errorf(
					"ierror.StackToError() error = '%v', wanted = '%v'",
					got,
					tt.wanted,
				)
			}
		})
	}
}

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
		{
			name: "It should add the code InvalidToken the new error",
			fields: fields{
				err: New(""),
			},
			exec: func(e *ierror) *ierror {
				return e.InvalidToken()
			},
			want: InvalidToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.exec(tt.fields.err)

			if !reflect.DeepEqual(got.code, tt.want) {
				t.Errorf("ErrCode = %v, want %v", got, tt.want)
			}
		})
	}
}
