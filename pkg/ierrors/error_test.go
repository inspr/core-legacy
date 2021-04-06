package ierrors

import (
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

// TODO test for errors.Is
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
	}{
		// TODO: Add test cases.
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
	}{
		// TODO: Add test cases.
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
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
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
				t.Errorf("InsprError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprError.MarshalJSON() = %v, want %v", got, tt.want)
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			if err := ierror.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("InsprError.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsprError_StackToError(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			ierror.StackToError()
		})
	}
}

func TestInsprError_FormatedError(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
		Code    InsprErrorCode
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierror := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			ierror.FormatedError()
		})
	}
}
