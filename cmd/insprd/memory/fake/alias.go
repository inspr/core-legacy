package fake

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Alias - mocks the implementation of the AppMemory interface methods
type Alias struct {
	*MemManager
	fail  error
	alias map[string]*meta.Alias
}

// Get - simple mock
func (a *Alias) Get(context string, aliasKey string) (*meta.Alias, error) {
	if a.fail != nil {
		return nil, a.fail
	}
	query := fmt.Sprintf("%s", context)
	alias, ok := a.alias[query]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message(fmt.Sprintf("alias %s not found", query)).Build()
	}

	return alias, nil
}

// CreateAlias - simple mock
func (a *Alias) CreateAlias(query string, targetBoundary string, alias *meta.Alias) error {
	if a.fail != nil {
		return a.fail
	}

	_, ok := a.alias[query]
	if ok {
		return ierrors.NewError().AlreadyExists().Message(fmt.Sprintf("alias %s already exists", query)).Build()
	}
	a.alias[query] = alias
	return nil
}

// DeleteAlias - simple mock
func (a *Alias) DeleteAlias(context string, aliasKey string) error {
	if a.fail != nil {
		return a.fail
	}
	query := fmt.Sprintf("%s", context)
	_, ok := a.alias[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}

	delete(a.alias, query)
	return nil
}

// UpdateAlias - simple mock
func (a *Alias) UpdateAlias(context string, aliasKey string, alias *meta.Alias) error {
	if a.fail != nil {
		return a.fail
	}
	query := fmt.Sprintf("%s", context)
	_, ok := a.alias[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}
	a.alias[query] = alias
	return nil
}
