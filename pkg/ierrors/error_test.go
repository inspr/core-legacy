package ierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"testing"
)

func TestInsprError_Error(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "testing the error message of inspr error",
			fields: fields{
				Message: "mock_message",
				Err:     nil,
				Stack:   "no_stack",
				Code:    0,
			},
			want: "mock_message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("InsprError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_Is(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Code    InsprErrorCode
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
			name: "It should return true since the two Errors Codes are equal",
			fields: fields{
				Code:    BadRequest,
				Err:     nil,
				Message: "A new message",
			},
			args: args{
				target: &InsprError{
					Code:    BadRequest,
					Message: "Another message",
					Err:     nil,
				},
			},
			want: true,
		},
		{
			name: "It should return false since the two Errors Codes are different",
			fields: fields{
				Code:    BadRequest,
				Err:     nil,
				Message: "A new message",
			},
			args: args{
				target: &InsprError{
					Code:    InternalServer,
					Message: "Another message",
					Err:     nil,
				},
			},
			want: false,
		},
		{
			name: "The error given is in the error stack",
			fields: fields{
				Err: fmt.Errorf(
					"layer2: %w",
					fmt.Errorf(
						"layer1: %w",
						fs.ErrClosed,
					),
				),
			},
			args: args{
				target: fs.ErrClosed,
			},
			want: true,
		},
		{
			name: "The error given is NOT in the error stack",
			fields: fields{
				Err: fmt.Errorf(
					"layer2: %w",
					fmt.Errorf(
						"layer1: %w",
						fs.ErrClosed,
					),
				),
			},
			args: args{
				target: fs.SkipDir,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Code:    tt.fields.Code,
			}
			if got := err.Is(tt.args.target); got != tt.want {
				t.Errorf("InsprError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_HasCode(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Code    InsprErrorCode
	}
	type args struct {
		code InsprErrorCode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "It should return true since the Error code and the code are equal",
			fields: fields{
				Code:    BadRequest,
				Err:     nil,
				Message: "A new message",
			},
			args: args{
				code: BadRequest,
			},
			want: true,
		},
		{
			name: "It should return false since the Error code and the code are different",
			fields: fields{
				Code:    BadRequest,
				Err:     nil,
				Message: "A new message",
			},
			args: args{
				code: AlreadyExists,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Code:    tt.fields.Code,
			}
			if got := err.HasCode(tt.args.code); got != tt.want {
				t.Errorf("InsprError.HasCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_Wrap(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "wrap_simple_test",
			fields: fields{
				Message: "",
				Err:     nil,
				Stack:   "",
				Code:    0,
			},
			args: args{
				message: "mock_message",
			},
			want: "mock_message",
		},
		{
			name: "wrap_composed_test",
			fields: fields{
				Message: "",
				Err:     errors.New("first"),
				Stack:   "",
				Code:    0,
			},
			args: args{
				message: "second",
			},
			want: "second: first",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			ierror.Wrap(tt.args.message)

			got := ierror.Stack
			if got != tt.want {
				t.Errorf("InsprError.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_Wrapf(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}
	type args struct {
		format string
		values []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "wrap_simple_test",
			fields: fields{
				Message: "",
				Err:     nil,
				Stack:   "",
				Code:    0,
			},
			args: args{
				format: "%s",
				values: []interface{}{"mock_message"},
			},
			want: "mock_message",
		},
		{
			name: "wrap_composed_test",
			fields: fields{
				Message: "",
				Err:     errors.New("first"),
				Stack:   "",
				Code:    0,
			},
			args: args{
				format: "%s-%s",
				values: []interface{}{"third", "second"},
			},
			want: "third-second: first",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			ierror.Wrapf(tt.args.format, tt.args.values...)

			got := ierror.Stack
			if got != tt.want {
				t.Errorf("InsprError.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_MarshalJSON(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}

	// mocking the insprErr and getting it's bytes representation
	ie := InsprError{
		Message: "mock",
		Err:     nil,
		Stack:   "mock",
		Code:    0,
	}
	bytes, _ := json.Marshal(ie)

	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:   "simple_marshall",
			fields: fields(ie),
			// json.Marshal result of the above structure inside an IError
			want:    bytes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			got, err := ierror.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"InsprError.MarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"InsprError.MarshalJSON() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestInsprError_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}
	type args struct {
		data []byte
	}

	wanted := InsprError{
		Message: "mock_error",
		Err:     errors.New("mock_error"),
		Stack:   "mock_error",
		Code:    0,
	}
	wantedBytes, _ := json.Marshal(wanted)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "basic_test",
			fields: fields{
				Message: wanted.Message,
				Err:     wanted.Err,
				Stack:   wanted.Stack,
				Code:    0,
			},
			args:    args{data: wantedBytes},
			wantErr: false,
		},
		{
			name:    "error_test",
			fields:  fields{},
			args:    args{data: []byte{0}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}

			// ierror.Unmarshal could be used but
			err := json.Unmarshal(tt.args.data, &ierror)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"InsprError.UnmarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestInsprError_StackToError(t *testing.T) {
	type fields struct {
		Err   error
		Stack string
	}
	tests := []struct {
		name   string
		fields fields
		wanted string // used for comparison with err.Error()
	}{
		{
			name: "basic_test",
			fields: fields{
				Err:   nil,
				Stack: "hello: stack: test",
			},
			wanted: "hello: stack: test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Err:   tt.fields.Err,
				Stack: tt.fields.Stack,
			}
			ierror.StackToError()

			got := ierror.Err.Error()
			if got != tt.wanted {
				t.Errorf(
					"InsprError.StackToError() error = %v, wanted = %v",
					got,
					tt.wanted,
				)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	mutipleErrors := NewError().InnerError(errors.New("mock")).Build()
	mutipleErrors.Wrap("new_error_message")

	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "unwrap_empty_err_stack",
			args:    args{err: nil},
			wantErr: false,
		},
		{
			name:    "unwrap_simple_err_stack",
			args:    args{err: errors.New("mock")},
			wantErr: false,
		},
		{
			name:    "unwrap_simple_inspr_err_stack",
			args:    args{err: NewError().InnerError(errors.New("mock")).Build()},
			wantErr: false,
		},
		{
			name:    "unwrap_complex_inspr_err_stack",
			args:    args{err: mutipleErrors},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
