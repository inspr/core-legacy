package tree

import "gitlab.inspr.dev/inspr/core/pkg/meta"

type MockAppManager struct {
	root *meta.App
	err  error
}

func (mock *MockAppManager) GetApp(query string) (*meta.App, error) {
	if mock.err != nil {
		return nil, mock.err
	}
	return mock.root, nil
}
func (mock *MockAppManager) CreateApp(app *meta.App, context string) error {
	return nil
}
func (mock *MockAppManager) DeleteApp(query string) error {
	return nil
}
func (mock *MockAppManager) UpdateApp(app *meta.App, query string) error {
	return nil
}
