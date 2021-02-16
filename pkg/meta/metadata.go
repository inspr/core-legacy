package meta

// Metadata represents an arbitrary inspr component. It represents this components' ID inside the cluster and the hub.
// Annotations can be used as query methods for getting components inside Inspr, just like kubernetes' annotations.
//
// The parent field represents the parent app of the object.
// The full reference for this component will be {it's parent's reference}.Name
type Metadata struct {
	Name        string            `yaml:"name" json:"name"`
	Reference   string            `yaml:"reference" json:"reference"`
	Annotations map[string]string `yaml:"annotations" json:"annotations"`
	Parent      string            `yaml:"parent" json:"parent"`
	SHA256      string
}
