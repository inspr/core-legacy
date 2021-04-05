package utils

import (
	"github.com/google/uuid"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// InjectUUID injects a new UUID on a metadata
func InjectUUID(m meta.Metadata) meta.Metadata {
	m.UUID = uuid.New().String()
	return m
}
