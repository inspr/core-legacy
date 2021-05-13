package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(inspr completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ inspr completion bash > /etc/bash_completion.d/inspr
  # macOS:
  $ inspr completion bash > /usr/local/etc/bash_completion.d/inspr

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ inspr completion zsh > "${fpath[1]}/_inspr"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ inspr completion fish | source

  # To load completions for each session, execute once:
  $ inspr completion fish > ~/.config/fish/completions/inspr.fish

PowerShell:

  PS> inspr completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> inspr completion powershell > inspr.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
			fmt.Print("compdef _inspr inspr")
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func completeDapps(cm *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return generateCompletion(func(app *meta.App) []string {
		scopes := []string{}
		for name := range app.Spec.Apps {
			scopes = append(scopes, name)
		}
		return scopes
	})(cm, args, toComplete)
}
func completeChannels(cm *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return generateCompletion(func(app *meta.App) []string {
		scopes := []string{}
		for name := range app.Spec.Apps {
			scopes = append(scopes, name+".")
		}
		for name := range app.Spec.Channels {
			scopes = append(scopes, name)
		}
		return scopes
	})(cm, args, toComplete)
}
func generateCompletion(tg func(*meta.App) []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cm *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		toComplete = strings.TrimSuffix(toComplete, ".")
		client := cliutils.GetCliClient()
		scope, err := cliutils.GetScope()

		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		newScope, err := utils.JoinScopes(scope, toComplete)
		if _, err := client.Apps().Get(context.Background(), newScope); err != nil {
			newScope, _, _ = utils.RemoveLastPartInScope(newScope)
		}

		app, err := client.Apps().Get(context.Background(), newScope)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		newOnes := tg(app)
		for i := range newOnes {
			if newScope != "" {
				newOnes[i] = newScope + "." + newOnes[i]
			}
		}
		ret := []string{}
		for _, s := range newOnes {
			if strings.HasPrefix(s, toComplete) {
				ret = append(ret, s)
			}
		}
		return ret, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

}
func completeTypes(cm *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return generateCompletion(func(app *meta.App) []string {
		scopes := []string{}
		for name := range app.Spec.Apps {
			scopes = append(scopes, name+".")
		}
		for name := range app.Spec.Types {
			scopes = append(scopes, name)
		}
		return scopes
	})(cm, args, toComplete)
}
func completeAliases(cm *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return generateCompletion(func(app *meta.App) []string {
		scopes := []string{}
		for name := range app.Spec.Apps {
			scopes = append(scopes, name+".")
		}
		for name := range app.Spec.Aliases {
			scopes = append(scopes, name)
		}
		return scopes
	})(cm, args, toComplete)
}
