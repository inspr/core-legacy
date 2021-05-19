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
func (cm *BrokersMock) Get(ctx context.Context) (*models.BrokersDi, error) {
	if cm.err != nil {
		return &models.BrokersDi{}, cm.err
	}
	return &models.BrokersDi{}, nil
}
