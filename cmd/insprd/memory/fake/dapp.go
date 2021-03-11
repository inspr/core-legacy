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
func (a *Apps) Get(query string) (*meta.App, error) {
	if a.fail != nil {
		return nil, a.fail
	}
	ct, ok := a.apps[query]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message(fmt.Sprintf("dapp %s not found", query)).Build()
	}
	return ct, nil
}

// CreateApp - simple mock
func (a *Apps) CreateApp(context string, ct *meta.App) error {
	if a.fail != nil {
		return a.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)

	_, ok := a.apps[query]
	if ok {
		return ierrors.NewError().AlreadyExists().Message(fmt.Sprintf("dapp %s already exists", query)).Build()
	}
	a.apps[query] = ct
	return nil
}

// DeleteApp - simple mock
func (a *Apps) DeleteApp(query string) error {
	if a.fail != nil {
		return a.fail
	}
	_, ok := a.apps[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("dapp %s not found", query)).Build()
	}

	delete(a.apps, query)
	return nil
}

// UpdateApp - simple mock
func (a *Apps) UpdateApp(context string, ct *meta.App) error {
	if a.fail != nil {
		return a.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := a.apps[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("dapp %s not found", query)).Build()
	}
	a.apps[query] = ct
	return nil
}

// ResolveBoundary mock
func (ch *Apps) ResolveBoundary(app *meta.App) (map[string]string, error) {
	ret := map[string]string{}
	for _, ch := range app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output) {
		ret[ch] = ch + "_resolved"
	}
	return ret, nil
}
