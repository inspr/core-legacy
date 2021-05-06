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
func (chType *Types) Get(context string, ctName string) (*meta.Type, error) {
	if chType.fail != nil {
		return nil, chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ctName)
	ct, ok := chType.Types[query]
	if !ok {
		return nil, ierrors.
			NewError().
			NotFound().
			Message("Type %s not found", query).
			Build()
	}
	return ct, nil
}

// Create - simple mock
func (chType *Types) Create(context string, ct *meta.Type) error {
	if chType.fail != nil {
		return chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := chType.Types[query]
	if ok {
		return ierrors.
			NewError().
			AlreadyExists().
			Message("Type %s already exists", query).
			Build()
	}
	chType.Types[query] = ct
	return nil
}

// Delete - simple mock
func (chType *Types) Delete(context string, ctName string) error {
	if chType.fail != nil {
		return chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ctName)
	_, ok := chType.Types[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("Type %s not found", query).
			Build()
	}

	delete(chType.Types, query)
	return nil
}

// Update - simple mock
func (chType *Types) Update(context string, ct *meta.Type) error {
	if chType.fail != nil {
		return chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := chType.Types[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("Type %s not found", query).
			Build()
	}
	chType.Types[query] = ct
	return nil
}
