package meta

/*
Alias defines an alias for a channel or a route. It allows resources defined in different contexts
to be accessed in a higher or lower context.

(Deprecated)Target is the channel which is being referenced by the alias

Meta.Name is the new name given for the resource.

Resource is the name of the resource (channel or route) that is being referenced by the alias.

Source is the dapp that contains the resource that is being referenced.  When filled in, it means
that the resource is available in the specified child of the dapp. When ommited, it means that the
resource is available in the parent of the current dapp.

Destination is the dapp in which the alias is assigned. When filled in, it means that the resource is
intended for a specific child of the dapp. When ommited, it means that the resource is available
outside the context of the dapp.

*/
type Alias struct {
	Meta        Metadata `yaml:"meta" json:"meta"`
	Target      string   `yaml:"target" json:"target"`
	Resource    string   `yaml:"resource" json:"resource"`
	Source      string   `yaml:"source" json:"source"`
	Destination string   `yaml:"destination" json:"destination"`
}
