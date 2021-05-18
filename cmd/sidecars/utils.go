package sidecars

import (
	"strings"

	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

func generateReadPort() int32 { return int32(0) }

func generateWritePort() int32 { return int32(0) }

// toAppID - creates the kubernetes deployment name from the app
func toAppID(app *meta.App) string {
	var depNames utils.StringArray
	depNames = strings.Split(app.Meta.Parent, ".")
	if depNames[0] == "" {
		depNames = utils.StringArray{}
	}
	depNames = append(depNames, app.Meta.Name)
	return depNames.Join("-")
}

// insprAppIDConfig adds the dapp id that the sidecar is related to
func insprAppIDConfig(app *meta.App) k8s.ContainerOption {
	return k8s.ContainerWithEnv(corev1.EnvVar{
		Name:  "INSPR_APP_ID",
		Value: toAppID(app),
	})
}
