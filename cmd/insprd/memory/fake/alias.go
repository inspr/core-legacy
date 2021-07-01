package fake

import (
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

// Alias - mocks the implementation of the AppMemory interface methods
type Alias struct {
	*TreeMemoryMock
	fail  error
	alias map[string]*meta.Alias
}

// Get - simple mock
func (a *Alias) Get(scope, aliasKey string) (*meta.Alias, error) {
	if a.fail != nil {
		return nil, a.fail
	}

	alias, ok := a.alias[scope]
	if !ok {
		return nil, ierrors.
			NewError().
			NotFound().
			Message("alias %s not found", scope).
			Build()
	}

	return alias, nil
}

// Create - simple mock
func (a *Alias) Create(query, targetBoundary string, alias *meta.Alias) error {
	if a.fail != nil {
		return a.fail
	}

	_, ok := a.alias[query]
	if ok {
		return ierrors.
			NewError().
			AlreadyExists().
			Message("alias %s already exists", query).
			Build()
	}
	a.alias[query] = alias
	return nil
}

// Delete - simple mock
func (a *Alias) Delete(scope, aliasKey string) error {
	if a.fail != nil {
		return a.fail
	}

	_, ok := a.alias[scope]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("type %s not found", scope).
			Build()
	}

	delete(a.alias, scope)
	return nil
}

// Update - simple mock
func (a *Alias) Update(scope, aliasKey string, alias *meta.Alias) error {
	if a.fail != nil {
		return a.fail
	}
	_, ok := a.alias[scope]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("type %s not found", scope).
			Build()
	}
	a.alias[scope] = alias
	return nil
}
