package cmd

import (
	"context"
	"errors"
	"testing"
)

func Test_deleteAction(t *testing.T) {
	type args struct {
		c context.Context
		s []string
	}
	tests := []struct {
		name    string
		args    args
		options deleteUsrOptionsDT
		wantErr bool
	}{
		{
			name:    "correct functionality testing",
			options: deleteUsrOptionsDT{username: "user to be deleted"},
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
		},
		{
			name:    "error testing",
			options: deleteUsrOptionsDT{username: "user to be deleted"},
			args: args{
				c: context.Background(),

				s: []string{"user", "password"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteUsrOptions = tt.options
			cl = mockCl{
				deleteUser: func(c context.Context, s1, s2, s3 string) error {
					if s1 != tt.args.s[0] {
						t.Errorf("username not set correctly \n%v\n!=\n%v", s1, tt.args.s[0])
					}
					if s2 != tt.args.s[1] {
						t.Errorf("password not set correctly \n%v\n!=\n%v", s2, tt.args.s[1])
					}
					if s3 != tt.options.username {
						t.Errorf("user to be deleted not set correctly \n%v\n!=\n%v", s3, tt.options.username)
					}

					if tt.wantErr {
						return errors.New("this is an error")
					}
					return nil
				},
			}
			if err := deleteAction(tt.args.c, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("deleteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
