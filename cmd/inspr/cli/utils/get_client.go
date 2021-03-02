package utils

import (
	"encoding/json"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

func GetClient() *client.Client {
	url := GetConfiguredServerIp()

	rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()
	return client.NewControllerClient(rc)
}
