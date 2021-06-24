package nodes

import (
	"os"

	"go.uber.org/zap"
)

type k8sVariables struct {
	AppsNamespace string
}

var variables *k8sVariables
var logger *zap.Logger

func getK8SVariables() *k8sVariables {
	if variables == nil {
		variables = &k8sVariables{
			AppsNamespace: getFromEnv("NODES_APPS_NAMESPACE"),
		}
	}
	return variables
}

func getFromEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found")
}

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "sidecar-channel-operator")))
}
