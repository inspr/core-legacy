package fake

import (
	"fmt"

	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

// Types - mocks the implementation of the TypeMemory interface methods
type Types struct {
	*TreeMemoryMock
	fail       error
	insprTypes map[string]*meta.Type
}

// Get - simple mock
func (t *Types) Get(scope, name string) (*meta.Type, error) {
	if t.fail != nil {
		return nil, t.fail
	}
	query := fmt.Sprintf("%s.%s", scope, name)
	ct, ok := t.insprTypes[query]
	if !ok {
		return nil, ierrors.New("type %s not found", query).NotFound()
	}
	return ct, nil
}

// Create - simple mock
func (t *Types) Create(scope string, tp *meta.Type) error {
	if t.fail != nil {
		return t.fail
	}
	query := fmt.Sprintf("%s.%s", scope, tp.Meta.Name)
	_, ok := t.insprTypes[query]
	if ok {
		return ierrors.New("type %s already exists", query).AlreadyExists()
	}
	t.insprTypes[query] = tp
	return nil
}

// Delete - simple mock
func (t *Types) Delete(scope, name string) error {
	if t.fail != nil {
		return t.fail
	}
	query := fmt.Sprintf("%s.%s", scope, name)
	_, ok := t.insprTypes[query]
	if !ok {
		return ierrors.New("type %s not found", query).NotFound()
	}

	delete(t.insprTypes, query)
	return nil
}

// Update - simple mock
func (t *Types) Update(scope string, tp *meta.Type) error {
	if t.fail != nil {
		return t.fail
	}
	query := fmt.Sprintf("%s.%s", scope, tp.Meta.Name)
	_, ok := t.insprTypes[query]
	if !ok {
		return ierrors.New("type %s not found", query).NotFound()
	}
	t.insprTypes[query] = tp
	return nil
}
