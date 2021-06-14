package fake

import (
	"fmt"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
)

// Types - mocks the implementation of the TypeMemory interface methods
type Types struct {
	*MemManager
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
		return nil, ierrors.
			NewError().
			NotFound().
			Message("type %s not found", query).
			Build()
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
		return ierrors.
			NewError().
			AlreadyExists().
			Message("type %s already exists", query).
			Build()
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
		return ierrors.
			NewError().
			NotFound().
			Message("type %s not found", query).
			Build()
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
		return ierrors.
			NewError().
			NotFound().
			Message("type %s not found", query).
			Build()
	}
	t.insprTypes[query] = tp
	return nil
}
