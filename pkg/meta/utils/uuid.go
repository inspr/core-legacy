package utils

import (
	"github.com/google/uuid"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func InjectUUID(m meta.Metadata) meta.Metadata {
	m.UUID = uuid.New().String()
	return m
}
