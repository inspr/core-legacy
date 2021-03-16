package fake

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// Alias - mocks the implementation of the AppMemory interface methods
type Alias struct {
	*MemManager
	fail  error
	alias map[string]*meta.Alias
}

// Get - simple mock
func (a *Alias) Get(context string, aliasKey string) (*meta.Alias, error) {
	return &meta.Alias{}, nil
}

// CreateAlias - simple mock
func (a *Alias) CreateAlias(query string, targetBoundary string, alias *meta.Alias) error {
	return nil
}

// DeleteAlias - simple mock
func (a *Alias) DeleteAlias(context string, aliasKey string) error {
	return nil
}

// UpdateAlias - simple mock
func (a *Alias) UpdateAlias(contexcontext string, aliasKey string, alias *meta.Alias) error {
	return nil
}
