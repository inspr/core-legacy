package fake

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Apps - mocks the implementation of the AppMemory interface methods
type Apps struct {
	*MemManager
	fail error
	apps map[string]*meta.App
}

// Get - simple mock
func (ch *Apps) Get(query string) (*meta.App, error) {
	if ch.fail != nil {
		return nil, ch.fail
	}
	ct, ok := ch.apps[query]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message(fmt.Sprintf("dapp %s not found", query)).Build()
	}
	return ct, nil
}

// CreateApp - simple mock
func (ch *Apps) CreateApp(context string, ct *meta.App) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := ch.apps[query]
	if ok {
		return ierrors.NewError().AlreadyExists().Message(fmt.Sprintf("dapp %s already exists", query)).Build()
	}
	ch.apps[query] = ct
	return nil
}

// DeleteApp - simple mock
func (ch *Apps) DeleteApp(query string) error {
	if ch.fail != nil {
		return ch.fail
	}
	_, ok := ch.apps[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("dapp %s not found", query)).Build()
	}

	delete(ch.apps, query)
	return nil
}

// UpdateApp - simple mock
func (ch *Apps) UpdateApp(context string, ct *meta.App) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := ch.apps[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("dapp %s not found", query)).Build()
	}
	ch.apps[query] = ct
	return nil
}
