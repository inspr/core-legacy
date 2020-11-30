package errors

import (
	"reflect"
	"testing"
)

func TestInsprError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *InsprError
		want string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("InsprError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}

	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundError(tt.args.name, tt.args.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.args.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEncodingError(t *testing.T) {
	type args struct {
		msg        string
		innerError error
	}

	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEncodingError(tt.args.msg, tt.args.innerError); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncodingError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEncoding(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEncoding(tt.args.err); got != tt.want {
				t.Errorf("IsEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAlreadyExistsError(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}

	tests := []struct {
		name string
		args args
		want *InsprError
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlreadyExistsError(tt.args.name, tt.args.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlreadyExistsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAlreadyExists(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlreadyExists(tt.args.err); got != tt.want {
				t.Errorf("IsAlreadyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		err     *InsprError
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("InsprError.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsprError_Is(t *testing.T) {
	type args struct {
		target error
	}

	tests := []struct {
		name string
		err  *InsprError
		args args
		want bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Is(tt.args.target); got != tt.want {
				t.Errorf("InsprError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsprError_ToNative(t *testing.T) {
	tests := []struct {
		name string
		err  *InsprError
		want interface{}
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.ToNative(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsprError.ToNative() = %v, want %v", got, tt.want)
			}
		})
	}
}
