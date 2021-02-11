package nodes

import "os"

type K8SVariables struct {
	AppsNamespace string
}

var variables *K8SVariables

func GetK8SVariables() *K8SVariables {
	if variables == nil {
		variables = &K8SVariables{
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
