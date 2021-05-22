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
func (cm *BrokersMock) Get(ctx context.Context) (*models.BrokersDI, error) {
	if cm.err != nil {
		return &models.BrokersDI{}, cm.err
	}
	return &models.BrokersDI{}, nil
}

// Create moscks the Brokers controller interface method
func (cm *BrokersMock) Create(ctx context.Context, brokerName string, config []byte) error {
	return cm.err
}
