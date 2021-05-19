package sidecars

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	corev1 "k8s.io/api/core/v1"
)

func Test_generateReadPort(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{
			name: "needs to be implemented",
			want: int32(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateReadPort(); got != tt.want {
				t.Errorf("generateReadPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateWritePort(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{
			name: "needs to be implemented",
			want: int32(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateWritePort(); got != tt.want {
				t.Errorf("generateWritePort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toAppID(t *testing.T) {
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "dapp with parent",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{Name: "dapp3", Parent: "dapp1.dapp2"},
				},
			},
			want: "dapp1-dapp2-dapp3",
		},
		{
			name: "dapp without parent",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{Name: "dapp1"},
				},
			},
			want: "dapp1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toAppID(tt.args.app); got != tt.want {
				t.Errorf("toAppID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insprAppIDConfig(t *testing.T) {
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name string
		args args
		want k8s.ContainerOption
	}{
		{
			name: "dappID_testing",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{Name: "dapp3", Parent: "dapp1.dapp2"},
				},
			},
			want: k8s.ContainerWithEnv(corev1.EnvVar{
				Name:  "INSPR_APP_ID",
				Value: "dapp1-dapp2-dapp3",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InsprAppIDConfig(tt.args.app)

			gotContainer := k8s.NewContainer("", "", got)
			wantContainer := k8s.NewContainer("", "", tt.want)

			if !reflect.DeepEqual(gotContainer, wantContainer) {

				t.Errorf("insprAppIDConfig() = %v, want %v",
					gotContainer,
					wantContainer)
			}
		})
	}
}
