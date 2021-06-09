package mocks

import (
	"context"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/controller"
)

// BrokersMock mock structure for the operations of the controller.Brokerss()
type BrokersMock struct {
	err error
}

// NewBrokersMock exports a mock of the Brokers.interface
func NewBrokersMock(err error) controller.BrokersInterface {
	return &BrokersMock{
		err: err,
	}
}

// Get is the BrokersMock Get
func (bm *BrokersMock) Get(ctx context.Context) (*models.BrokersDI, error) {
	if bm.err != nil {
		return &models.BrokersDI{}, bm.err
	}
	return &models.BrokersDI{}, nil
}
