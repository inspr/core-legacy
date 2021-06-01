package fake

import "github.com/inspr/inspr/pkg/sidecars/models"

// Factory mock of AbsstractBrokersFactory
type Factory struct {
	fail     error
	abstract map[string]models.SidecarFactory
}

// Subscribe mock of factory subscription method
func (f *Factory) Subscribe(broker string, factory models.SidecarFactory) error {
	return f.fail
}

// Get mock of factory get method
func (f *Factory) Get(broker string) (models.SidecarFactory, error) {
	if f.fail != nil {
		return nil, f.fail
	}
	return nil, nil
}
