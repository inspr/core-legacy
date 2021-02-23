package cli

import (
	"context"
	"encoding/json"
	"io"

	"github.com/spf13/viper"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

/* INTEGRATION
Given the current state of the tasks some things will change when integrating everything,
this is a list of things to do  related to this package

	- remove the RunMethod and use the one defined in the method Factory
	- Establish a prerun in the apply cmd where the factory.Subscribe is done for this method
	- in the channels().create -> the false parameter is the cmd.InsprOptions.DryRunFlag variable
		- replace the mock_flags with the ones in the apply cmd
*/

// RunMethod defines the method that will run for the component
type RunMethod func(data []byte, out io.Writer) error

// ApplyChannel is of the type RunMethod, it calls the pkg/controller/client functions.
var ApplyChannel RunMethod = func(data []byte, out io.Writer) error {
	url := viper.GetString("reqUrl")

	rc := request.NewClient().
		BaseURL(url).
		Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).
		Build()

	// create controller client
	c := client.NewControllerClient(rc)

	// unmarshal into a channel
	channel, err := utils.YamlToChannel(data)
	if err != nil {
		return err
	}

	// INTEGRATION: REMOVE, flags from other story
	// todo use the flags from the apply cmd
	flagDryRun := false
	flagIsUpdate := false

	var clog diff.Changelog
	// creates or updates it
	if flagIsUpdate {
		clog, err = c.Channels().Update(context.Background(), channel.Meta.Parent, &channel, flagDryRun)
	} else {
		clog, err = c.Channels().Create(context.Background(), channel.Meta.Parent, &channel, flagDryRun)
	}

	if err != nil {
		return err
	}

	// prints differences
	clog.Print(out)

	return nil
}
