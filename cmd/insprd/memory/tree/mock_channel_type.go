package tree

import (
	"github.com/inspr/inspr/pkg/meta"
)

// TypeMockManager mocks a Type interface for testing
type TypeMockManager struct {
	*MockManager
}

// Create mocks a Type method for testing
func (ctm *TypeMockManager) Create(context string, ct *meta.Type) error {
	return nil
}

// Get mocks a Type method for testing
func (ctm *TypeMockManager) Get(context string, ctName string) (*meta.Type, error) {
	return nil, nil
}

// Delete mocks a Type method for testing
func (ctm *TypeMockManager) Delete(context string, ctName string) error {
	return nil
}

// Update mocks a Type method for testing
func (ctm *TypeMockManager) Update(query string, ct *meta.Type) error {
	return nil
}
