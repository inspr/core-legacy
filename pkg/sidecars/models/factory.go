package models

import (
	"github.com/inspr/inspr/pkg/meta"
	corev1 "k8s.io/api/core/v1"
)

// SidecarConnections object to store a sidecar's connectio ports
type SidecarConnections struct {
	InPort  int32
	OutPort int32
}

// SidecarFactory function type responsible for creating a sidecar for a broker
type SidecarFactory func(app *meta.App, conn *SidecarConnections) corev1.Container
