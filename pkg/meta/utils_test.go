package meta

import (
	"testing"

	test "gitlab.inspr.dev/inspr/core/pkg/testutils"
	"gotest.tools/v3/assert/cmp"
	kubeCore "k8s.io/api/core/v1"
)

func TestStructureNameIsValid(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid lowercase app name",
			args: args{
				name: "onetwotree",
			},
			wantErr: false,
		},
		{
			name: "Valid lowercase app name with numbers",
			args: args{
				name: "0one1two2tree3",
			},
			wantErr: false,
		},
		{
			name: "Valid uppercase app name",
			args: args{
				name: "ONETWOTREE",
			},
			wantErr: false,
		},
		{
			name: "Valid uppercase app name with numbers",
			args: args{
				name: "0ONE1TWO2TREE3",
			},
			wantErr: false,
		},
		{
			name: "Valid app name with '-'",
			args: args{
				name: "ONE1-two2-TREE3",
			},
			wantErr: false,
		},
		{
			name: "Valid app name with '_'",
			args: args{
				name: "ONE1-two_2-TREE3",
			},
			wantErr: false,
		},
		{
			name: "Invalid app name with starting '-'",
			args: args{
				name: "-ONE1-two2-TREE3",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with ending '-'",
			args: args{
				name: "ONE1-two2-TREE3-",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with '%'",
			args: args{
				name: "ONE1-two%2-TREE3",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with '/'",
			args: args{
				name: "ONE1-two2-TREE3/",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with '&'",
			args: args{
				name: "ONE1-two&2-TREE3",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with '#'",
			args: args{
				name: "ONE1-two#2-TREE3",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with '@'",
			args: args{
				name: "ONE1-two@2-TREE3-",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with length >= 64 characters",
			args: args{
				name: "qwertyuiopasdfghjkl12345678901234567890zxcvbnmasdfghjklqwert3456",
			},
			wantErr: true,
		},
		{
			name: "Invalid app name with 0 length",
			args: args{
				name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StructureNameIsValid(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructureNameIsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func Test_parseToK8sArrEnv(t *testing.T) {
	type args struct {
		arrappEnv EnvironmentMap
	}
	tests := []struct {
		name string
		args args
		want []kubeCore.EnvVar
	}{
		{
			name: "successful_test",
			args: args{
				arrappEnv: map[string]string{
					"key_1": "value_1",
					"key_2": "value_2",
					"key_3": "value_3",
				},
			},
			want: []kubeCore.EnvVar{
				{
					Name:  "key_1",
					Value: "value_1",
				},
				{
					Name:  "key_2",
					Value: "value_2",
				},
				{
					Name:  "key_3",
					Value: "value_3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.arrappEnv.ParseToK8sArrEnv()

			var comp cmp.Comparison = cmp.DeepEqual(got, tt.want, test.GetMapCompareOptions())
			if !comp().Success() {
				t.Errorf("parseToK8sArrEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
