package cli

import (
	"context"
	"fmt"

	"github.com/inspr/inspr/pkg/cmd"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
)

var getBorkers = cmd.NewCmd("brokers").
	WithDescription("Retreives brokers curently installed on cluster").
	NoArgs(
		func(ctx context.Context) error {
			client := cliutils.GetCliClient()
			out := cliutils.GetCliOutput()
			resp, err := client.Brokers().Get(context.Background())
			if err != nil {
				cliutils.RequestErrorMessage(err, out)
				return err
			}

			fmt.Fprintf(out, "DEFAULT:\n%s\n", resp.Default)
			fmt.Fprintln(out, "AVAILABLE:")
			lines := []string(resp.Installed)
			printTab(&lines)
			return nil
		})
