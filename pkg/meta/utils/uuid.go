package utils

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

func InjectUUID(m meta.Metadata) meta.Metadata {
	m.UUID = utils.NewUUID()
	return m
}
