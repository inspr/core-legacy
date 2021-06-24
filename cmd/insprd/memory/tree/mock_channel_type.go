package tree

import (
	"inspr.dev/inspr/pkg/meta"
)

// TypeMockManager mocks a Type interface for testing
type TypeMockManager struct {
	*MockManager
}

// Create mocks a Type method for testing
func (ctm *TypeMockManager) Create(scope string, ct *meta.Type) error {
	return nil
}

// Get mocks a Type method for testing
func (ctm *TypeMockManager) Get(scope, name string) (*meta.Type, error) {
	return nil, nil
}

// Delete mocks a Type method for testing
func (ctm *TypeMockManager) Delete(scope, name string) error {
	return nil
}

// Update mocks a Type method for testing
func (ctm *TypeMockManager) Update(query string, ct *meta.Type) error {
	return nil
}
