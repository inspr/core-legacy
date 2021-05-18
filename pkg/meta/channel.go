package meta

import "github.com/inspr/inspr/pkg/utils"

// Channel is an Inspr component that represents a Channel.
type Channel struct {
	Meta             Metadata          `yaml:"meta,omitempty"  json:"meta"`
	Spec             ChannelSpec       `yaml:"spec,omitempty"  json:"spec"`
	ConnectedApps    utils.StringArray `yaml:"connectedapps,omitempty"  json:"connectedapps"`
	ConnectedAliases utils.StringArray `yaml:"connectedaliases,omitempty"  json:"connectedaliasses"`
}

// ChannelSpec is the specification of a channel.
// 'Type' string references a Type structure name
type ChannelSpec struct {
	Type   string `yaml:"type,omitempty"  json:"type" `
	Broker string `yaml:"broker,omitempty" json:"broker"`
}
