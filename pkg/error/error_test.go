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

func TestNewCustomError(t *testing.T) {
	type args struct {
		errCode InsprErrorCode
		errMsg  string
	}
	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		{
			name: "It should create a new Custom Inspr Error",
			args: args{
				errCode: BadRequest,
				errMsg:  "A brand new bad request error message",
			},
			want: &InsprError{
				Code:    BadRequest,
				Message: "A brand new bad request error message",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCustomError(tt.args.errCode, tt.args.errMsg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	type args struct {
		name string
		err  error
	}
	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		{
			name: "It should create a new Not Found Error",
			args: args{
				name: "Example",
				err:  nil,
			},
			want: &InsprError{
				Code:    NotFound,
				Message: "Component Example not found.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundError(tt.args.name, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAlreadyExistsError(t *testing.T) {
	type args struct {
		name string
		err  error
	}
	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		{
			name: "It should create a new Already Exists Error",
			args: args{
				name: "Example",
				err:  nil,
			},
			want: &InsprError{
				Code:    AlreadyExists,
				Message: "Component Example already exists.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlreadyExistsError(tt.args.name, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlreadyExistsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInternalServerError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		{
			name: "It should create a new Internal Server Error",
			args: args{
				err: nil,
			},
			want: &InsprError{
				Code:    InternalServer,
				Message: "There was a internal server error.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInternalServerError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInternalServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInvalidNameError(t *testing.T) {
	type args struct {
		name string
		err  error
	}
	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		{
			name: "It should create a new Invalid Name Error",
			args: args{
				name: "A invalid name example",
				err:  nil,
			},
			want: &InsprError{
				Code:    InvalidName,
				Message: "The name 'A invalid name example' is invalid.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInvalidNameError(tt.args.name, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvalidNameError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInvalidChannelError(t *testing.T) {
	tests := []struct {
		name string
		want *InsprError
	}{
		{
			name: "It should create a new Invalid Channel Error",
			want: &InsprError{
				Code:    InvalidChannel,
				Message: "The channel is invalid.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInvalidChannelError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvalidChannelError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInvalidAppError(t *testing.T) {
	tests := []struct {
		name string
		want *InsprError
	}{
		{
			name: "It should create a new Invalid app Error",
			want: &InsprError{
				Code:    InvalidApp,
				Message: "The app is invalid.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInvalidAppError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvalidAppError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInvalidChannelTypeError(t *testing.T) {
	tests := []struct {
		name string
		want *InsprError
	}{
		{
			name: "It should create a new Invalid Channel Type Error",
			want: &InsprError{
				Code:    InvalidChannelType,
				Message: "The ChannelType is invalid.",
				Err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInvalidChannelTypeError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvalidChannelTypeError() = %v, want %v", got, tt.want)
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
