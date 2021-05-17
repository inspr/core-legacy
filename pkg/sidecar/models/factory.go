package models

import (
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
)

// SidecarConnections object to store a sidecar's connectio ports
type SidecarConnections struct {
	InPort  int32
	OutPort int32
}

// SidecarFactory function type responsible for creating a sidecar for a broker
type SidecarFactory func(app *meta.App, conn SidecarConnections) k8s.DeploymentOption
