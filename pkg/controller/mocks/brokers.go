package mocks

import (
	"context"

	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/controller"
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

// Create mocks the Brokers controller interface method
func (bm *BrokersMock) Create(ctx context.Context, brokerName string, config []byte) error {
	return bm.err
}
