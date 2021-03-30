package tree

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// MockAppManager to Mock App Manager
type MockAppManager struct {
	*MockManager
	err error
}

// Get Mock
func (mock *MockAppManager) Get(query string) (*meta.App, error) {
	if mock.err != nil {
		return nil, mock.err
	}
	return mock.root, nil
}

// Create Mock
func (mock *MockAppManager) Create(context string, app *meta.App) error {
	return nil
}

// Delete Mock
func (mock *MockAppManager) Delete(query string) error {
	return nil
}

// Update Mock
func (mock *MockAppManager) Update(query string, app *meta.App) error {
	return nil
}

// ResolveBoundary Mock
func (mock *MockAppManager) ResolveBoundary(app *meta.App) (map[string]string, error) {
	return nil, nil
}
