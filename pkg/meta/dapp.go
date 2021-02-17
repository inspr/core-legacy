package meta

// Node represents an inspr component that is a node.
type Node struct {
	Meta Metadata `yaml:"meta,omitempty"  json:"meta"`
	Spec NodeSpec `yaml:"spec,omitempty"  json:"spec"`
}

// NodeSpec represents a configuration for a node. The image represents the Docker image for the main container of the Node.
type NodeSpec struct {
	Image string `yaml:"image,omitempty"  json:"image"`
}

// App is an inspr component that represents an App. An App can contain other apps, channels and other components.
type App struct {
	Meta Metadata `yaml:"meta,omitempty" json:"meta"`
	Spec AppSpec  `yaml:"spec,omitempty" json:"spec"`
}

// AppBoundary represents the connections this app can make to other apps. These are the fields that can be overriten
// by the ChannelAliases when instantiating the app.
type AppBoundary struct {
	Input  []string `yaml:"input,omitempty" json:"input"`
	Output []string `yaml:"output,omitempty" json:"output"`
}

// AppSpec represents the configuration of an App.
//
// The app contains a list of apps and a list of nodes. The apps and nodes can be dereferenced by it's metadata
// reference, at CLI time.
//
// The boundary represent the possible connections to other apps, and the fields that can be overriten when instantiating the app.
type AppSpec struct {
	Node         Node                    `yaml:"node,omitempty" json:"node"`
	Apps         map[string]*App         `yaml:"apps,omitempty" json:"apps"`
	Channels     map[string]*Channel     `yaml:"channels,omitempty" json:"channels"`
	ChannelTypes map[string]*ChannelType `yaml:"channeltypes,omitempty" json:"channeltypes"`
	Boundary     AppBoundary             `yaml:"boundary,omitempty" json:"boundary"`
}
