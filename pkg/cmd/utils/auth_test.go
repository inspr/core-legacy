package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAuthenticator_GetToken(t *testing.T) {
	os.Mkdir("./temp", os.ModePerm)
	folder := "./temp"
	defer os.RemoveAll("./temp")
	type fields struct {
		tokenPath string
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
				tokenPath: filepath.Join(folder, "token"),
			},
			want:    []byte("Bearer this is a token"),
			wantErr: false,
			before: func(auth Authenticator) {
				ioutil.WriteFile(filepath.Join(folder, "token"), []byte("this is a token"), os.ModePerm)
			},
		},
		{
			name: "invalid token location",
			fields: fields{
				tokenPath: filepath.Join(folder, "token2"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Authenticator{
				tokenPath: tt.fields.tokenPath,
			}
			if tt.before != nil {
				tt.before(a)
			}
			got, err := a.GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticator.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticator.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
