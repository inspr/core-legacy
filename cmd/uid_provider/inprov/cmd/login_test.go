package cmd

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func Test_login(t *testing.T) {
	type args struct {
		ctx      context.Context
		login    string
		password string
		file     string
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
				password: "password",
				file:     "file",
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
			loginOptions = loginOptionsDT{
				output: tt.args.file,
			}
			defer os.Remove(tt.args.file)
			if err := login(tt.args.ctx, tt.args.login, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_loginAction(t *testing.T) {
	recoverStdout := os.Stdout
	var writer *os.File
	var reader *os.File
	defer func() {
		os.Stdout = recoverStdout
	}()
	type args struct {
		c context.Context
	}
	tests := []struct {
		name      string
		args      args
		before    func()
		after     func()
		getReader func()
		options   loginOptionsDT
		wantErr   bool
	}{
		{
			name: "stdout output",
			args: args{
				c: context.Background(),
			},
			before: func() {
				var err error
				reader, writer, err = os.Pipe()
				if err != nil {
					panic(err)
				}
				os.Stdout = writer
			},
			options: loginOptionsDT{
				stdout:   true,
				user:     "username",
				password: "password",
			},
			after: func() {
				os.Stdout = recoverStdout
			},
		},

		{
			name: "file output",
			args: args{
				c: context.Background(),
			},
			options: loginOptionsDT{
				output:   "afile",
				user:     "username2",
				password: "password2",
			},
			getReader: func() {
				reader, _ = os.Open("afile")
			},
			after: func() {
				os.Remove("afile")
			},
		},

		{
			name: "invalid file",
			args: args{
				c: context.Background(),
			},
			options: loginOptionsDT{
				output:   "/usr/bin/afile",
				user:     "username3",
				password: "password3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer = nil, nil
			if tt.before != nil {
				tt.before()
			}
			loginOptions = tt.options
			if tt.after != nil {
				defer tt.after()
			}

			cl = mockCl{
				login: func(c context.Context, s1, s2 string) (string, error) {
					if s1 != tt.options.user {
						t.Errorf("username is not set correctly\n%v\n!=\n%v", s1, tt.options.user)
					}
					if s2 != tt.options.password {
						t.Errorf("password is not set correctly\n%v\n!=\n%v", s2, tt.options.password)
					}
					return "this is a test", nil
				},
			}

			if err := loginAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("loginAction() error = %v, wantErr %v", err, tt.wantErr)
			}
			writer.Close()

			if tt.getReader != nil {
				tt.getReader()
			}
			bytes, err := ioutil.ReadAll(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("loginAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.options.stdout {
				if (string(bytes) != "Successfully logged in!\nthis is a test") != tt.wantErr {
					t.Errorf("loginAction() = %v, want %v", string(bytes), "Successfully logged in!\nthis is a test")
				}

			} else {
				if (string(bytes) != "this is a test") != tt.wantErr {
					t.Errorf("loginAction() = %v, want %v", string(bytes), "this is a test")
				}
			}
		})
	}
}
