// Package applymethods is the package with the functions
// to be inserted into the factory method
package applymethods

import (
	"encoding/json"

	"github.com/spf13/viper"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

/* INTEGRATION
Given the current state of the tasks some things will change when integrating everything,
this is a list of things to do  related to this package

- remove the RunMethod and use the one defined in the method Factory
- Establisha a prerun in the apply cmd where the factory.Subscribe is done for this method
- how to determine
*/

// RunMethod defines the method that will run for the component
type RunMethod func([]byte) error

// ApplyChannel is of the type RunMethod, it calls the pkg/controller/client functions.
var ApplyChannel RunMethod = func([]byte) error {
	url := viper.Get("reqUrl")
	rc := request.NewClient().
		BaseURL(url).
		Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).
		Build()

	// create controller client
	c := client.NewControllerClient(rc)

	// unmarshal into a channel
	// channel, err :=

	// creates or updates it
	// = c.Channels().Create(context.Background(),)
	return nil
}
