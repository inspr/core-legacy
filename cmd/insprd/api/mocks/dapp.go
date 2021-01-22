package mocks

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Apps todo doc
type Apps struct {
	memory.AppMemory
}

// GetApp todo doc
func (apps *Apps) GetApp(query string) (*meta.App, error) {
	return &meta.App{}, nil
}

// CreateApp todo doc
func (apps *Apps) CreateApp(app *meta.App, context string) error {
	return nil
}

// DeleteApp todo doc
func (apps *Apps) DeleteApp(query string) error {
	return nil
}

// UpdateApp todo doc
func (apps *Apps) UpdateApp(app *meta.App, query string) error {
	return nil
}
