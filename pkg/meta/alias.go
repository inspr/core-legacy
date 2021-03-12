package meta

import "gitlab.inspr.dev/inspr/core/pkg/utils"

// Alias defines an alias for a channel
//
// Target is the channel which is being referenced by the alias
type Alias struct {
	Target        string            `yaml:"target"  json:"target"`
	ConnectedApps utils.StringArray `yaml:"connectedapps,omitempty"  json:"connectedapps"`
}
