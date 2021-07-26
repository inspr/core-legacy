package ierrors

/*
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

/*

func Testierror_Is(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Code    ErrCode
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
				target: New("another message").BadRequest(),
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
				target: New("Another").InternalServer(),
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
			err := &ierror{
				err:  tt.fields.err,
				code: tt.fields.code,
			}
			if got := err.Is(tt.args.target); got != tt.want {
				t.Errorf("ierror.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Testierror_HasCode(t *testing.T) {
	type fields struct {
		Message string
		Err     error
		Code    ErrCode
	}
	type args struct {
		code ErrCode
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
			err := &ierror{
				Message: tt.fields.Message,
				Err:     tt.fields.Err,
				Code:    tt.fields.Code,
			}
			if got := err.HasCode(tt.args.code); got != tt.want {
				t.Errorf("ierror.HasCode() = %v, want %v", got, tt.want)
			}
		})
	}
}


*/
