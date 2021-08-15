package fake

import (
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
)

// Apps - mocks the implementation of the AppMemory interface methods
type Apps struct {
	*TreeMemoryMock
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
		return nil, ierrors.New("dapp %s not found", query).NotFound()
	}
	return ct, nil
}

// Create - simple mock
func (a *Apps) Create(
	scope string,
	app *meta.App,
	brokers *apimodels.BrokersDI,
) error {
	if a.fail != nil {
		return a.fail
	}
	query, _ := utils.JoinScopes(scope, app.Meta.Name)

	_, ok := a.apps[query]
	if ok {
		return ierrors.New("dapp %s already exists", query).AlreadyExists()
	}
	a.apps[query] = app
	return nil
}

// Delete - simple mock
func (a *Apps) Delete(query string) error {
	if a.fail != nil {
		return a.fail
	}
	_, ok := a.apps[query]
	if !ok {
		return ierrors.New("dapp %s not found", query).NotFound()
	}

	delete(a.apps, query)
	return nil
}

// Update - simple mock
func (a *Apps) Update(
	scope string,
	app *meta.App,
	brokers *apimodels.BrokersDI,
) error {
	if a.fail != nil {
		return a.fail
	}
	query, _ := utils.JoinScopes(scope, app.Meta.Name)
	_, ok := a.apps[query]
	if !ok {
		return ierrors.New("dapp %s not found", query).NotFound()
	}
	a.apps[query] = app
	return nil
}

// ResolveBoundary mock
func (*Apps) ResolveBoundary(
	app *meta.App,
	usePermTree bool,
) (map[string]string, error) {
	ret := map[string]string{}
	for _, ch := range app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output) {
		ret[ch] = ch
	}
	return ret, nil
}
