package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAuthenticator_GetToken(t *testing.T) {
	folder := t.TempDir()

	type fields struct {
		TokenPath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
		before  func(auth Authenticator)
	}{
		{
			name: "valid token location",
			fields: fields{
				TokenPath: filepath.Join(folder, "token"),
			},
			want:    []byte("Bearer this is a token"),
			wantErr: false,
			before: func(auth Authenticator) {
				ioutil.WriteFile(
					filepath.Join(folder, "token"),
					[]byte("this is a token"),
					os.ModePerm,
				)
			},
		},
		{
			name: "invalid token location",
			fields: fields{
				TokenPath: filepath.Join(folder, "token2"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Authenticator{
				TokenPath: tt.fields.TokenPath,
			}
			if tt.before != nil {
				tt.before(a)
			}
			got, err := a.GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Authenticator.GetToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticator.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticator_SetToken(t *testing.T) {
	folder := t.TempDir()

	type fields struct {
		TokenPath string
	}
	type args struct {
		token []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantString string
		wantErr    bool
	}{
		{
			name: "base_setToken",
			fields: fields{
				TokenPath: filepath.Join(folder, "token"),
			},
			args: args{
				token: []byte("Bearer mock_token"),
			},
			wantErr:    false,
			wantString: "mock_token",
		},
		{
			name: "error_setToken",
			fields: fields{
				TokenPath: filepath.Join(folder, "token/err_file"),
			},
			args: args{
				token: []byte("Bearer mock_token"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Authenticator{
				TokenPath: tt.fields.TokenPath,
			}
			err := a.SetToken(tt.args.token)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Authenticator.SetToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			content, _ := ioutil.ReadFile(tt.fields.TokenPath)

			got := string(content)
			if got != tt.wantString {
				t.Errorf(
					"Authenticator.SetToken(), STRING error = %v, wantErr %v",
					got,
					tt.wantString,
				)
			}
		})
	}
}
