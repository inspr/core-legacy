package client

import (
	"context"
	"net/http"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest/request"
)

// BrokersClient interacts with Brokers on the Insprd
type BrokersClient struct {
	reqClient *request.Client
}

// Get gets a brokers from the Insprd
func (bc *BrokersClient) Get(ctx context.Context) (*models.BrokersDI, error) {
	resp := &models.BrokersDI{}

	err := bc.reqClient.Send(ctx, "/brokers", http.MethodGet, nil, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Create creates a broker into the cluster via insprd
func (bc *BrokersClient) Create(ctx context.Context, brokerName string, config []byte) error {
	dataBody := models.BrokerDataDI{
		BrokerName:   brokerName,
		FileContents: config,
	}
	// TODO: how should i receive messages from this request? what composes the response body
	err := bc.reqClient.Send(
		ctx,
		"/brokers/"+brokerName,
		http.MethodPost,
		dataBody,
		nil)
	return err
}
