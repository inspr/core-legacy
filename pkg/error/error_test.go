package errors

import (
	"reflect"
	"testing"
)

func TestInsprError_Error(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Code    InsprErrorCode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "It should return the InsprError message",
			fields: fields{
				Code:    BadRequest,
				Message: "A test message",
				Err:     nil,
			},
			want: "A test message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InsprError{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Code:    tt.fields.Code,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("InsprError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name string
		want *ErrBuilder
	}{
		{
			name: "It should return a empty Inspr Err Build",
			want: &ErrBuilder{
				err: &InsprError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
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
		{
			name: "It should add the code Not Found to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: NotFound,
				},
			},
		},
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
		{
			name: "It should add the code Already Exists to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: AlreadyExists,
				},
			},
		},
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
		{
			name: "It should add the code Bad Request to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: BadRequest,
				},
			},
		},
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
		{
			name: "It should add the code Internal Server to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: InternalServer,
				},
			},
		},
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
		{
			name: "It should add the code Invalid Name to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: InvalidName,
				},
			},
		},
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
		{
			name: "It should add the code Invalid App to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: InvalidApp,
				},
			},
		},
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
		{
			name: "It should add the code Invalid Channel to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: InvalidChannel,
				},
			},
		},
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
		{
			name: "It should add the code Invalid Channel Type to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: InvalidChannelType,
				},
			},
		},
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

func TestErrBuilder_Message(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrBuilder
	}{
		{
			name: "It should add a message to the new Inspr Error",
			fields: fields{
				err: &InsprError{},
			},
			args: args{
				msg: "A brand new error message",
			},
			want: &ErrBuilder{
				err: &InsprError{
					Message: "A brand new error message",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.Message(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
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
		{
			name: "It should add a inner error to the new Inspr Error",
			fields: fields{
				err: &InsprError{},
			},
			args: args{
				err: NewError().AlreadyExists().Message("Hello").Build(),
			},
			want: &ErrBuilder{
				err: &InsprError{
					Err: &InsprError{
						Code:    AlreadyExists,
						Message: "Hello",
					},
				},
			},
		},
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
		{
			name: "It should return the created error",
			fields: fields{
				err: &InsprError{
					Code:    NotFound,
					Message: "A brand new error message",
				},
			},
			want: &InsprError{
				Code:    NotFound,
				Message: "A brand new error message",
			},
		},
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
