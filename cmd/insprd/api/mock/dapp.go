package repos

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Apps todo doc
type Apps struct {
	memory.AppMemory
}

// GetApp todo doc
func (apps *Apps) GetApp(ref string) (*meta.AppSpec, error) {
	return &meta.AppSpec{}, nil
}

// CreateApp todo doc
func (apps *Apps) CreateApp(appSpec *meta.AppSpec) error {
	return nil
}

// DeleteApp todo doc
func (apps *Apps) DeleteApp(ref string) error {
	return nil
}

// UpdateApp todo doc
func (apps *Apps) UpdateApp(appSpec *meta.AppSpec, ref string) error {
	return nil
}
