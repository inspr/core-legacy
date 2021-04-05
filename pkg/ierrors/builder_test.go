package ierrors

import (
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name string
		want *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_Message(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	type args struct {
		format string
		values []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.Message(tt.args.format, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_InnerError(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InnerError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InnerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_Build(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *InsprError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_Is(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Stack   string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
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
		Stack   string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Stack:   tt.fields.Stack,
				Code:    tt.fields.Code,
			}
			if got := err.HasCode(tt.args.code); got != tt.want {
				t.Errorf("InsprError.HasCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_NotFound(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.NotFound(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.NotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_AlreadyExists(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.AlreadyExists(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.AlreadyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_BadRequest(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.BadRequest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.BadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_InternalServer(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InternalServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InternalServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_InvalidName(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InvalidName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InvalidName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_InvalidApp(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InvalidApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InvalidApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_InvalidChannel(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InvalidChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InvalidChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_InvalidChannelType(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InvalidChannelType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InvalidChannelType() = %v, want %v", got, tt.want)
			}
		})
	}
}
