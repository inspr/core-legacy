package meta

// Alias defines an alias for a channel
// Target is the channel which is being referenced by the alias
type Alias struct {
	Meta   Metadata `yaml:"meta" json:"meta"`
	Target string   `yaml:"target"  json:"target"`
}
