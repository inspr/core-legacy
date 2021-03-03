package utils

import (
	"encoding/json"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

//GetClient creates and returns a new client for requesting on inspr deamon.
//About to be deprecated.
func GetClient() *client.Client {
	url := GetConfiguredServerIP()

	rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()
	return client.NewControllerClient(rc)
}
