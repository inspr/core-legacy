package tree

import "gitlab.inspr.dev/inspr/core/pkg/meta"

//MockAppManager Mock
type MockAppManager struct {
	root *meta.App
	err  error
}

//GetApp mock
func (mock *MockAppManager) GetApp(query string) (*meta.App, error) {
	if mock.err != nil {
		return nil, mock.err
	}
	return mock.root, nil
}

//CreateApp mock
func (mock *MockAppManager) CreateApp(app *meta.App, context string) error {
	return nil
}

//DeleteApp mock
func (mock *MockAppManager) DeleteApp(query string) error {
	return nil
}

//UpdateApp mock
func (mock *MockAppManager) UpdateApp(app *meta.App, query string) error {
	return nil
}
