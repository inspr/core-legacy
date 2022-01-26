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
func (a *Alias) Get(scope, name string) (*meta.Alias, error) {
	if a.fail != nil {
		return nil, a.fail
	}

	alias, ok := a.alias[scope]
	if !ok {
		return nil, ierrors.New("alias %s not found", scope).NotFound()
	}

	return alias, nil
}

// Create - simple mock
func (a *Alias) Create(scope string, alias *meta.Alias) error {
	if a.fail != nil {
		return a.fail
	}

	_, ok := a.alias[scope]
	if ok {
		return ierrors.New("alias %s already exists", scope).AlreadyExists()
	}
	a.alias[scope] = alias
	return nil
}

// Delete - simple mock
func (a *Alias) Delete(scope, name string) error {
	if a.fail != nil {
		return a.fail
	}

	_, ok := a.alias[scope]
	if !ok {
		return ierrors.New("type %s not found", scope).NotFound()
	}

	delete(a.alias, scope)
	return nil
}

// Update - simple mock
func (a *Alias) Update(scope string, alias *meta.Alias) error {
	if a.fail != nil {
		return a.fail
	}
	_, ok := a.alias[scope]
	if !ok {
		return ierrors.New("type %s not found", scope).NotFound()
	}
	a.alias[scope] = alias
	return nil
}
