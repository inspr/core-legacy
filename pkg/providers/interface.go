package providers

import (
	"io"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/pkg/ierrors"
)

// Provider is the interface for a provider for insprd operators,
type Provider interface {
	New(name string, op func(m memory.Manager, r io.Reader) operators.OperatorInterface) error
	Get(name string) (func(m memory.Manager, r io.Reader) operators.OperatorInterface, error)
}

var p providerFactory

// GetProvider gets the provider for insprd operators
func GetProvider() Provider {
	if p == nil {
		p = providerFactory{}
	}
	return p
}

type providerFactory map[string]func(memory.Manager, io.Reader) operators.OperatorInterface

func (p providerFactory) New(name string, op func(memory.Manager, io.Reader) operators.OperatorInterface) error {
	if _, ok := p[name]; ok {
		return ierrors.NewError().Message("provider already exists").AlreadyExists().Build()
	}
	p[name] = op
	return nil
}

func (p providerFactory) Get(name string) (op func(memory.Manager, io.Reader) operators.OperatorInterface, err error) {
	if op, ok := p[name]; ok {
		return op, nil

	}
	return nil, ierrors.NewError().NotFound().Message("provider %s does note exist", name).Build()
}
