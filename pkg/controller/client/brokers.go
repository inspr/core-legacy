package client

import (
	"context"
	"net/http"

	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/rest/request"
)

// BrokersClient interacts with Brokers on the Insprd
type BrokersClient struct {
	reqClient *request.Client
}

// Get gets a brokers from the Insprd
func (bc *BrokersClient) Get(ctx context.Context) (*models.BrokersDI, error) {
	resp := &models.BrokersDI{}

	err := bc.reqClient.
		Header(rest.HeaderScopeKey, "").
		Send(ctx, "/brokers", http.MethodGet, nil, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Create creates a broker into the cluster via insprd
func (bc *BrokersClient) Create(ctx context.Context, brokerName string, config []byte) error {
	dataBody := models.BrokerConfigDI{
		BrokerName:   brokerName,
		FileContents: config,
	}
	err := bc.reqClient.
		Header(rest.HeaderScopeKey, "").
		Send(ctx, "/brokers/"+brokerName, http.MethodPost, dataBody, nil)

	return err
}
