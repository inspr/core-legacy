package mocks

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Apps - mocks the implementation of the AppMemory interface methods
type Apps struct {
	fail error
	*MemManager
}

// Get - simple mock
func (apps *Apps) Get(query string) (*meta.App, error) {
	if apps.fail != nil {
		return &meta.App{}, apps.fail
	}
	return &meta.App{}, nil
}

// CreateApp - simple mock
func (apps *Apps) CreateApp(app *meta.App, context string) error {
	if apps.fail != nil {
		return apps.fail
	}
	return nil
}

// DeleteApp - simple mock
func (apps *Apps) DeleteApp(query string) error {
	if apps.fail != nil {
		return apps.fail
	}
	return nil
}

// UpdateApp - simple mock
func (apps *Apps) UpdateApp(app *meta.App, query string) error {
	if apps.fail != nil {
		return apps.fail
	}
	return nil
}
