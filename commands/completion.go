package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate completion scripts for various shells.

To enable shell completion, follow these instructions:

Bash:
  $ canvas completion bash > /etc/bash_completion.d/canvas

  # Or for current user only:
  $ canvas completion bash > ~/.local/share/bash-completion/completions/canvas
  $ source ~/.bashrc

Zsh:
  $ canvas completion zsh > "${fpath[1]}/_canvas"

  # Or add to your .zshrc:
  $ canvas completion zsh > ~/.canvas-completion.zsh
  $ echo "source ~/.canvas-completion.zsh" >> ~/.zshrc
  $ source ~/.zshrc

Fish:
  $ canvas completion fish > ~/.config/fish/completions/canvas.fish
  $ source ~/.config/fish/config.fish

PowerShell:
  PS> canvas completion powershell > canvas.ps1
  PS> . .\canvas.ps1

  # Or add to your PowerShell profile:
  PS> canvas completion powershell >> $PROFILE`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(ExactArgsWithUsage(1, "shell"), cobra.OnlyValidArgs),
	RunE:      runCompletion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) error {
	shell := args[0]

	switch shell {
	case "bash":
		return rootCmd.GenBashCompletionV2(os.Stdout, true)
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	case "powershell":
		return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return fmt.Errorf("invalid shell: %s. Valid shells are: bash, zsh, fish, powershell", shell)
	}
}
