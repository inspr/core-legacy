package meta

// Component defines a general component so
// we can determine how to apply a given file to the cluster
type Component struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
}
