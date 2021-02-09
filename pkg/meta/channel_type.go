package meta

// ChannelType is the type of the channel. It can be a reference to an outsourced type or can be a local type. This local
// type will be defined via the workspace and instantiated as a []byte on the cluster
type ChannelType struct {
	Meta   Metadata `yaml:"meta" json:"meta" diff:"ctypemeta"`
	Schema []byte   `yaml:"schema" json:"schema" diff:"schema"`
	ConnectedChannels []string `yaml:"connectedchannels"  json:"connectedchannels" diff:"connectedchannels"`
}
