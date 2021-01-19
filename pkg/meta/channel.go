package meta

// Channel is an Inspr component that represents a Channel.
type Channel struct {
	Meta Metadata    `yaml:"meta"  json:"meta"`
	Spec ChannelSpec `yaml:"spec"  json:"spec"`
}

// ChannelSpec is the specification of a channel. (the external variable is just an idea)
type ChannelSpec struct {
	Type ChannelType `yaml:"type"  json:"type"`
}
