package operator

import (
	meta "gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ClusterOperator defines an interface to operate on cluster
type ClusterOperator interface {
	// Nodes
	CreateNode(node *meta.Node) error
	DeleteNode(name string) error
	UpdateNode(name string, node *meta.Node) error
	GetNode(name string) (*meta.Node, error)
	ListNodes() ([]*meta.Node, error)

	// Apps
	CreateApp(app *meta.App) error
	DeleteApp(name string) error
	UpdateApp(name string, app *meta.App) error
	GetApp(name string) (*meta.App, error)
	ListApps() ([]*meta.App, error)

	// Channels
	CreateChannel(channel *meta.Channel) error
	DeleteChannel(name string) error
	UpdateChannel(name string, channel *meta.Channel) error
	GetChannel(name string) (*meta.Channel, error)
	ListChannels() ([]*meta.Channel, error)

	// Utils
	UpdateNodeStatus() map[string]string
}

// MessageOperator defines an interface to read / write messages
// on message brokers
type MessageOperator interface {
	// Message Broker basic operator
	//
	// This Interface should be parallel, easy to understand and
	// easy to read.
	//
	OpenConnection() error
	ReadNextMessageFrom(channel string) (string, error)
	WriteMessageAt(message string, channel string) error
	CloseConnection() error
}
