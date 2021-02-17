package meta

// Channel is an Inspr component that represents a Channel.
type Channel struct {
	Meta          Metadata    `yaml:"meta,omitempty"  json:"meta"`
	Spec          ChannelSpec `yaml:"spec,omitempty"  json:"spec"`
	ConnectedApps []string    `yaml:"connectedapps,omitempty"  json:"connectedapps"`
}

// ChannelSpec is the specification of a channel.
// 'Type' string references a Channel Type structure name
type ChannelSpec struct {
	Type string `yaml:"type,omitempty"  json:"type" `
}
