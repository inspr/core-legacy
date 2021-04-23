package cli

import (
	"context"
	"fmt"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/cmd/utils"
)

var authCommand = cmd.NewCmd("auth").AddSubCommand(
	cmd.NewCmd("init").WithCommonFlags().NoArgs(
		func(c context.Context) error {
			token, err := utils.GetCliClient().Authorization().Init(c)
			if err != nil {
				return err
			}
			output := utils.GetCliOutput()

			fmt.Fprintln(output, "This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.")
			fmt.Fprintf(output, "%s\n", token)
			return nil
		},
	),
).WithCommonFlags().Super()
