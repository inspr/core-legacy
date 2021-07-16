package meta

import "inspr.dev/inspr/pkg/utils"

// Node represents an inspr component that is a node.
type Node struct {
	Meta Metadata `yaml:"meta,omitempty"  json:"meta"`
	Spec NodeSpec `yaml:"spec,omitempty"  json:"spec"`
}

// NodePort represents a connection for a node
type NodePort struct {
	Port       int `yaml:"port" json:"port"`
	TargetPort int `yaml:"targetPort" json:"targetPort"`
}

// SidecarPort represents the port for communication between node and load balancer sidecar
type SidecarPort struct {
	LBRead  int `json:"lbRead"`
	LBWrite int `json:"lbWrite"`
}

// NodeSpec represents a configuration for a node. The image represents the Docker image for the main container of the Node.
type NodeSpec struct {
	Ports         []NodePort           `yaml:"ports,omitempty" json:"ports,omitempty"`
	Image         string               `yaml:"image,omitempty"  json:"image"`
	Replicas      int                  `yaml:"replicas,omitempty" json:"replicas"`
	RestartPolicy string               `yaml:"restartPolicy,omitempty" json:"restartPolicy"`
	Environment   utils.EnvironmentMap `yaml:"environment,omitempty" json:"environment"`
	SidecarPort   SidecarPort          `yaml:"sidecarPort,omitempty" json:"sidecarPort"`
}

// App is an inspr component that represents an dApp. An App can contain other apps, channels and other components.
type App struct {
	Meta Metadata `yaml:"meta,omitempty" json:"meta"`
	Spec AppSpec  `yaml:"spec,omitempty" json:"spec"`
}

// AppBoundary represents the connections this app can make to other apps. These are the fields that can be overriten
// by the ChannelAliases when instantiating the app.
type AppBoundary struct {
	Input  utils.StringArray `yaml:"input,omitempty" json:"input"`
	Output utils.StringArray `yaml:"output,omitempty" json:"output"`
}

// AppSpec represents the configuration of an App.
//
// The boundary represent the possible connections to other apps, and the fields that can be overriten when instantiating the app.
type AppSpec struct {
	Node     Node                `yaml:"node,omitempty"   json:"node"`
	Apps     map[string]*App     `yaml:"apps,omitempty"   json:"apps"`
	Channels map[string]*Channel `yaml:"channels,omitempty"   json:"channels"`
	Types    map[string]*Type    `yaml:"types,omitempty"   json:"types"`
	Aliases  map[string]*Alias   `yaml:"aliases"   json:"aliases"`
	Boundary AppBoundary         `yaml:"boundary,omitempty"   json:"boundary"`
	Auth     AppAuth             `yaml:"auth"  json:"auth"`
	LogLevel string              `yaml:"logLevel" json:"logLevel"`
}

// AppAuth represents the permissions that a dApp (and its children) contains
type AppAuth struct {
	Scope       string            `yaml:"scope"  json:"scope"`
	Permissions utils.StringArray `yaml:"permissions"  json:"permissions"`
}
