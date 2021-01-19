package meta

// ChannelType is the type of the channel. It can be a reference to an outsourced type or can be a local type. This local
// type will be defined via the workspace and instantiated as a []byte on the cluster
type ChannelType struct {
	Metadata `yaml:"metadata" json:"metadata"`
	Schema   []byte `yaml:"schema" json:"schema,omitempty"  json:"schema"`
}
