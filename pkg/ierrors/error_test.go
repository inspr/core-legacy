package ierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	msg := "mockError"
	want := ierror{
		err:  errors.New(msg),
		code: Unknown,
	}

	got := New(msg)

	if !reflect.DeepEqual(got.err.Error(), want.err.Error()) {
		t.Errorf("Expected '%v', got '%v'", want.err, got.err)
	}
	if !reflect.DeepEqual(got.code, want.code) {
		t.Errorf("Expected %v, got %v", want.code, got.code)
	}
}

// TODO add Error case with different code
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
			want:   fmt.Sprintf("Code %d : mock_message", Unknown),
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
			want: "Code 1 : mock_err",
		},

		{
			name: "wrap_standard_error_with_message",
			fields: fields{
				err: errors.New("mock_err"),
			},
			args: args{
				msg: "wrapper_context",
			},
			want: "Code 1 : wrapper_context : mock_err",
		},
		{
			name: "wrap_standard_error_with_formatted_message",
			fields: fields{
				err: errors.New("mock_err"),
			},
			args: args{
				msg: "%w wrapper_context",
			},
			want: "Code 1 : mock_err wrapper_context",
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
			want: errors.New("mock"),
		},
		{
			name: "unwrap_ierror_without_previous_wrap",
			args: args{err: New("mock")},
			want: nil,
		},
		{
			name: "unwrap_ierror_with_previous_wrap",
			args: args{err: Wrap(New("mock"), "simple_wrap")},
			want: New("mock"),
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
					Wrap(
						Wrap(
							New("mock"),
							"first_context",
						),
						"second_context",
					),
					"third context",
				),
			},
			want: Wrap(
				Wrap(
					New("mock"),
					"first_context",
				),
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
		name    string
		args    args
		want    *ierror
		wantErr bool
	}{
		// {
		// 	name:    "empty_bytes",
		// 	args:    args{data: []byte{}},
		// 	wantErr: true,
		// },
		{
			name: "unmarshal_simple_ierror",
			args: args{data: generateIerrorBytes(New("mock_err"))},
			want: New("mock_err"),
		},
		// {
		// 	name:    "unmarshal_wrapped_error",
		// 	args:    args{data: []byte{0}},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ierr := New("")
			fmt.Println(tt.args.data, ierr)
			err := json.Unmarshal(tt.args.data, &ierr)

			if (err == nil) != tt.wantErr {
				t.Errorf(
					"json.Unmarshal(ierror) got = %v, wanted %v",
					err,
					tt.wantErr,
				)
			}

			if ierr.Error() != tt.want.Error() {
				t.Errorf(
					"json.Unmarshal(ierror) got = %v, wanted %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

/*

func Testierror_StackToError(t *testing.T) {
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
			ierror := &ierror{
				Err:   tt.fields.Err,
				Stack: tt.fields.Stack,
			}
			ierror.StackToError()

			got := ierror.Err.Error()
			if got != tt.wanted {
				t.Errorf(
					"ierror.StackToError() error = %v, wanted = %v",
					got,
					tt.wanted,
				)
			}
		})
	}
}
*/
