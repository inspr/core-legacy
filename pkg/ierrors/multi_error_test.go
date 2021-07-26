package ierrors

/*
func TestMultiError_Error(t *testing.T) {
	type fields struct {
		Errors []error
		Code   InsprErrorCode
	}
	tests := []struct {
		name    string
		fields  fields
		wantRet string
	}{
		{
			name: "Empty error",
			fields: fields{
				Errors: []error{},
				Code:   128,
			},
			wantRet: "",
		},
		{
			name: "Multiple errors",
			fields: fields{
				Errors: []error{
					NewError().Message("err1").Build(),
					NewError().Message("err2").Build(),
				},
				Code: 128,
			},
			wantRet: "err1\nerr2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MultiError{
				Errors: tt.fields.Errors,
				Code:   tt.fields.Code,
			}
			if gotRet := e.Error(); gotRet != tt.wantRet {
				t.Errorf("MultiError.Error() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestMultiError_Add(t *testing.T) {
	type fields struct {
		Errors []error
		Code   InsprErrorCode
	}
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRet string
	}{
		{
			name: "Adds error",
			fields: fields{
				Errors: []error{NewError().Message("err1").Build()},
				Code:   128,
			},
			args: args{
				NewError().Message("err2").Build(),
			},
			wantRet: "err1\nerr2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MultiError{
				Errors: tt.fields.Errors,
				Code:   tt.fields.Code,
			}
			e.Add(tt.args.err)
			if e.Error() != tt.wantRet {
				t.Errorf("MultiError.Empty() = %v, want %v", e.Error(), tt.wantRet)
			}
		})
	}
}

func TestMultiError_Empty(t *testing.T) {
	type fields struct {
		Errors []error
		Code   InsprErrorCode
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Empty error",
			fields: fields{
				Errors: []error{},
				Code:   128,
			},
			want: true,
		},
		{
			name: "Non empty error",
			fields: fields{
				Errors: []error{
					NewError().Message("err2").Build(),
				},
				Code: 128,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MultiError{
				Errors: tt.fields.Errors,
				Code:   tt.fields.Code,
			}
			if got := e.Empty(); got != tt.want {
				t.Errorf("MultiError.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
