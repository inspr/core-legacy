package meta

// Node represents an inspr component that is a node.
type Node struct {
	Meta Metadata `yaml:"meta"  json:"meta"`
	Spec NodeSpec `yaml:"spec"  json:"spec"`
}

// NodeSpec represents a configuration for a node. The image represents the Docker image for the main container of the Node.
// If the node has an specific Kubernetes configuration, the configuration can be injected via the Kubernetes field. When
// Kubernetes is set, the Image field gets igored.
type NodeSpec struct {
	Image string `yaml:"image"  json:"image"`
}

// App is an inspr component that represents an App. An App can contain other apps, channels and other components.
type App struct {
	Metadata `yaml:"metadata" json:"metadata"`
	Spec     AppSpec `yaml:"spec" json:"spec"`
}

// AppBoundary represents the connections this app can make to other apps. These are the fields that can be overriten
// by the ChannelAliases when instantiating the app.
type AppBoundary struct {
	Input  []string `yaml:"input" json:"input"`
	Output []string `yaml:"output" json:"output"`
}

// AppSpec represents the configuration of an App.
//
// The app contains a list of apps and a list of nodes. The apps and nodes can be dereferenced by it's metadata
// reference, at CLI time.
//
// The boundary represent the possible connections to other apps, and the fields that can be overriten when instantiating the app.
type AppSpec struct {
	Node     Node               `yaml:"node" json:"node"`
	Apps     map[string]App     `yaml:"apps" json:"apps"`
	Channels map[string]Channel `yaml:"channels" json:"channels"`
	Boundary AppBoundary        `yaml:"boundary" json:"boundary"`
}
