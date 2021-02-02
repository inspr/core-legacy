package meta

// Channel is an Inspr component that represents a Channel.
type Channel struct {
	Meta          Metadata        `yaml:"meta"  json:"meta" diff:"channelmeta"`
	Spec          ChannelSpec     `yaml:"spec"  json:"spec" diff:"channelspec"`
	ConnectedApps map[string]*App `yaml:"connectedapps"  json:"connectedapps" diff:"chconnectedapps"`
}

// ChannelSpec is the specification of a channel.
// 'Type' string references a Channel Type structure name
type ChannelSpec struct {
	Type string `yaml:"type"  json:"type" diff:"type"`
}
