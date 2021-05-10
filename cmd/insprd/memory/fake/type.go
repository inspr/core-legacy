package fake

import (
	"fmt"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
)

// Types - mocks the implementation of the TypeMemory interface methods
type Types struct {
	*MemManager
	fail  error
	Types map[string]*meta.Type
}

// Get - simple mock
func (t *Types) Get(context string, ctName string) (*meta.Type, error) {
	if t.fail != nil {
		return nil, t.fail
	}
	query := fmt.Sprintf("%s.%s", context, ctName)
	ct, ok := t.Types[query]
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
func (t *Types) Create(context string, ct *meta.Type) error {
	if t.fail != nil {
		return t.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := t.Types[query]
	if ok {
		return ierrors.
			NewError().
			AlreadyExists().
			Message("type %s already exists", query).
			Build()
	}
	t.Types[query] = ct
	return nil
}

// Delete - simple mock
func (t *Types) Delete(context string, ctName string) error {
	if t.fail != nil {
		return t.fail
	}
	query := fmt.Sprintf("%s.%s", context, ctName)
	_, ok := t.Types[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("type %s not found", query).
			Build()
	}

	delete(t.Types, query)
	return nil
}

// Update - simple mock
func (t *Types) Update(context string, ct *meta.Type) error {
	if t.fail != nil {
		return t.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := t.Types[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("type %s not found", query).
			Build()
	}
	t.Types[query] = ct
	return nil
}
