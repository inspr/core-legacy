package client

import (
	"context"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest/request"
)

// BrokersClient interacts with Brokers on the Insprd
type BrokersClient struct {
	reqClient *request.Client
}

// Get gets a brokers from the Insprd
func (cc *BrokersClient) Get(ctx context.Context) (*models.BrokersDI, error) {
	resp := &models.BrokersDI{}

	err := cc.reqClient.Send(ctx, "/brokers", "GET", nil, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
