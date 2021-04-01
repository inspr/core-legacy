package auth

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name string
		want Builder
	}{
		{
			name: "NewUser test",
			want: &builder{
				usr: User{
					Scopes: make([]string, 0),
					Role:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_SetName(t *testing.T) {
	type fields struct {
		usr User
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Builder
	}{
		{
			name: "SetName test",
			fields: fields{
				usr: User{
					Scopes: make([]string, 0),
					Role:   0,
				},
			},
			args: args{
				name: "name",
			},
			want: &builder{
				usr: User{
					Name:   "name",
					Scopes: make([]string, 0),
					Role:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &builder{
				usr: tt.fields.usr,
			}
			if got := bd.SetName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.SetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_SetPassword(t *testing.T) {
	type fields struct {
		usr User
	}
	type args struct {
		pwd string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Builder
	}{
		{
			name: "SetPassword test",
			fields: fields{
				usr: User{
					Scopes: make([]string, 0),
					Role:   0,
				},
			},
			args: args{
				pwd: "pass",
			},
			want: &builder{
				usr: User{
					Password: "pass",
					Scopes:   make([]string, 0),
					Role:     0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &builder{
				usr: tt.fields.usr,
			}
			if got := bd.SetPassword(tt.args.pwd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.SetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_SetScope(t *testing.T) {
	type fields struct {
		usr User
	}
	type args struct {
		scope []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Builder
	}{
		{
			name: "SetScope test",
			fields: fields{
				usr: User{
					Scopes: make([]string, 0),
					Role:   0,
				},
			},
			args: args{
				scope: []string{"scope", "scopo"},
			},
			want: &builder{
				usr: User{
					Scopes: []string{"scope", "scopo"},
					Role:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &builder{
				usr: tt.fields.usr,
			}
			if got := bd.SetScope(tt.args.scope...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.SetScope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_AsAdmin(t *testing.T) {
	type fields struct {
		usr User
	}
	tests := []struct {
		name   string
		fields fields
		want   Builder
	}{

		{
			name: "AsAdmin test",
			fields: fields{
				usr: User{
					Scopes: make([]string, 0),
					Role:   0,
				},
			},

			want: &builder{
				usr: User{
					Role:   1,
					Scopes: make([]string, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &builder{
				usr: tt.fields.usr,
			}
			if got := bd.AsAdmin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.AsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}
