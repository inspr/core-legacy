package cli

import (
	"context"
	"fmt"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/cmd/utils"
)

var clusterCommand = cmd.NewCmd("cluster").AddSubCommand(
	cmd.NewCmd("init").
		WithCommonFlags().
		ExactArgs(1,
			func(c context.Context, args []string) error {
				token, err := utils.GetCliClient().Authorization().Init(c, args[0])
				if err != nil {
					return err
				}
				output := utils.GetCliOutput()

				fmt.Fprintln(output, "This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.")
				fmt.Fprintf(output, "%s\n", token)
				return nil
			},
		),
	getBorkers,
).WithCommonFlags().Super()
