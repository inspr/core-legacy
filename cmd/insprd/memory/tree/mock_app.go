package tree

import (
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta"
)

// MockAppManager to Mock App Manager
type MockAppManager struct {
	*MockManager
	err error
}

// Get Mock
func (mock *MockAppManager) Get(scope string) (*meta.App, error) {
	if mock.err != nil {
		return nil, mock.err
	}
	return mock.root, nil
}

// Create Mock
func (mock *MockAppManager) Create(scope string, app *meta.App, brokers *apimodels.BrokersDI) error {
	return nil
}

// Delete Mock
func (mock *MockAppManager) Delete(scope string) error {
	return nil
}

// Update Mock
func (mock *MockAppManager) Update(scope string, app *meta.App, brokers *apimodels.BrokersDI) error {
	return nil
}

// ResolveBoundary Mock
func (mock *MockAppManager) ResolveBoundary(app *meta.App) (map[string]string, error) {
	return nil, nil
}
