package meta

// ChannelType is the type of the channel. It can be a reference to an outsourced type or can be a local type. This local
// type will be defined via the workspace and instantiated as a string on the cluster
type ChannelType struct {
	Meta              Metadata `yaml:"meta,omitempty" json:"meta"`
	Schema            []byte   `yaml:"schema,omitempty" json:"schema"`
	ConnectedChannels []string `yaml:"connectedchannels,omitempty"  json:"connectedchannels"`
}
