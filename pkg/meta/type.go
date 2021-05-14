package meta

// Type is a struct that can be a reference to an outsourced type or can be a local type, one of it's uses is the definition of the content type that is going to be sent in a channel.
//
// type will be defined via the workspace and instantiated as a string on the cluster
type Type struct {
	Meta              Metadata `yaml:"meta,omitempty" json:"meta"`
	Schema            string   `yaml:"schema,omitempty" json:"schema"`
	ConnectedChannels []string `yaml:"connectedchannels,omitempty"  json:"connectedchannels"`
}
