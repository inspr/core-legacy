package nodes

import "os"

type k8sVariables struct {
	AppsNamespace string
}

var variables *k8sVariables

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
