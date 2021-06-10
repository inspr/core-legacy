package meta

// Metadata represents an arbitrary Inspr component. It represents this components' ID inside the cluster and the hub.
// Annotations can be used as query methods for getting components inside Inspr, just like kubernetes' annotations.
// The parent field represents the parent app of the object.
type Metadata struct {
	Name        string            `yaml:"name" json:"name"`
	Reference   string            `yaml:"reference,omitempty" json:"reference"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations"`
	Parent      string            `yaml:"parent,omitempty" json:"parent"`
	UUID        string            `yaml:"uuid,omitempty"`
}
