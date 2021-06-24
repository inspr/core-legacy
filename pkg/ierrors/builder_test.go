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
			got := b.Message(tt.args.msg)
			if !reflect.DeepEqual(got.Build().Message, tt.want.Build().Message) {
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

			got := b.InnerError(tt.args.err).Build().Err.Error()
			want := tt.want.Build().Err.Error()
			if got != want {
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

func TestErrBuilder_InvalidType(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		{
			name: "It should add the code Invalid Type to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: InvalidType,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.InvalidType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.InvalidType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_Unauthorized(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		{
			name: "It should add the code Unauthorized to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: Unauthorized,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.Unauthorized(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.Unauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBuilder_Forbidden(t *testing.T) {
	type fields struct {
		err *InsprError
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrBuilder
	}{
		{
			name: "It should add the code Forbidden to the new error",
			fields: fields{
				err: &InsprError{},
			},
			want: &ErrBuilder{
				err: &InsprError{
					Code: Forbidden,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrBuilder{
				err: tt.fields.err,
			}
			if got := b.Forbidden(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBuilder.Forbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}
