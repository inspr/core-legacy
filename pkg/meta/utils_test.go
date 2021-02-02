package meta

import "testing"

func TestStructureNameIsValid(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid lowercase app name",
			args: args{
				name: "onetwotree",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid lowercase app name with numbers",
			args: args{
				name: "0one1two2tree3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid uppercase app name",
			args: args{
				name: "ONETWOTREE",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid uppercase app name with numbers",
			args: args{
				name: "0ONE1TWO2TREE3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid app name with '-'",
			args: args{
				name: "ONE1-two2-TREE3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid app name with '_'",
			args: args{
				name: "ONE1-two_2-TREE3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid app name with starting '-'",
			args: args{
				name: "-ONE1-two2-TREE3",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with ending '-'",
			args: args{
				name: "ONE1-two2-TREE3-",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with '%'",
			args: args{
				name: "ONE1-two%2-TREE3",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with '/'",
			args: args{
				name: "ONE1-two/2-TREE3",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with '&'",
			args: args{
				name: "ONE1-two&2-TREE3",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with '#'",
			args: args{
				name: "ONE1-two#2-TREE3",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with '@'",
			args: args{
				name: "ONE1-two@2-TREE3-",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid app name with length >= 64 characters",
			args: args{
				name: "qwertyuiopasdfghjkl12345678901234567890zxcvbnmasdfghjklqwert3456",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StructureNameIsValid(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructureNameIsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StructureNameIsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
