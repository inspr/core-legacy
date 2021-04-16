package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func Test_doStuff(t *testing.T) {
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
		err        error
	}{
		{
			name: "correct functionality",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			wantOutput: "this is a token",
			wantErr:    false,
		},
		{
			name: "server returns an error",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			wantOutput: "",
			wantErr:    true,
			err:        errors.New("this is an error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl = mockCl{
				login: func(c context.Context, s1, s2 string) (string, error) {
					if s1 != tt.args.login {
						t.Errorf("login is not passed correctly \n%v\n!=\n%v", s1, tt.args.login)
					}
					if s2 != tt.args.password {
						t.Errorf("password is not passed correctly \n%v\n!=\n%v", s2, tt.args.password)
					}
					if tt.err != nil {
						return "", tt.err
					}

					return "this is a token", nil
				},
			}
			output := &bytes.Buffer{}
			if err := login(tt.args.ctx, tt.args.login, tt.args.password, output); (err != nil) != tt.wantErr {
				t.Errorf("doStuff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("doStuff() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
