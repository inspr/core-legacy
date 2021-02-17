package meta

// Node represents an inspr component that is a node.
type Node struct {
	Meta Metadata `yaml:"meta"  json:"meta"`
	Spec NodeSpec `yaml:"spec"  json:"spec"`
}

// NodeSpec represents a configuration for a node. The image represents the Docker image for the main container of the Node.
type NodeSpec struct {
	Image       string         `yaml:"image"  json:"image" diff:"image"`
	Replicas    int            `yaml:"replicas" json:"replicas" diff:"replicas"`
	Environment EnvironmentMap `yaml:"envioronment" json:"envioronment" diff:"envioronment"`
}

// App is an inspr component that represents an App. An App can contain other apps, channels and other components.
type App struct {
	Meta Metadata `yaml:"meta" json:"meta"`
	Spec AppSpec  `yaml:"spec" json:"spec"`
}

// AppBoundary represents the connections this app can make to other apps. These are the fields that can be overriten
// by the ChannelAliases when instantiating the app.
type AppBoundary struct {
	Input  StringArray `yaml:"input" json:"input"`
	Output StringArray `yaml:"output" json:"output"`
}

// AppSpec represents the configuration of an App.
//
// The app contains a list of apps and a list of nodes. The apps and nodes can be dereferenced by it's metadata
// reference, at CLI time.
//
// The boundary represent the possible connections to other apps, and the fields that can be overriten when instantiating the app.
type AppSpec struct {
	Node         Node                    `yaml:"node" json:"node"`
	Apps         map[string]*App         `yaml:"apps" json:"apps"`
	Channels     map[string]*Channel     `yaml:"channels" json:"channels"`
	ChannelTypes map[string]*ChannelType `yaml:"channeltypes" json:"channeltypes"`
	Boundary     AppBoundary             `yaml:"boundary" json:"boundary"`
}
