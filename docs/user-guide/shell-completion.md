# Shell Completion

Canvas CLI supports tab completion for Bash, Zsh, Fish, and PowerShell.

## Installation

### Bash

```bash
# Add to ~/.bashrc
source <(canvas completion bash)

# Or save to a file
canvas completion bash > /etc/bash_completion.d/canvas
```

### Zsh

```bash
# Add to ~/.zshrc
source <(canvas completion zsh)

# Or add to fpath
canvas completion zsh > "${fpath[1]}/_canvas"
```

!!! tip "Oh My Zsh"
    If you use Oh My Zsh, save the completion script to:
    ```bash
    canvas completion zsh > ~/.oh-my-zsh/completions/_canvas
    ```

### Fish

```bash
canvas completion fish > ~/.config/fish/completions/canvas.fish
```

### PowerShell

```powershell
# Add to your PowerShell profile
canvas completion powershell | Out-String | Invoke-Expression

# Or save to a file
canvas completion powershell > canvas.ps1
```

## Usage

Once installed, press `Tab` to:

- Complete command names
- Complete flag names
- Complete flag values (where supported)

## Examples

```bash
# Complete commands
canvas cour<Tab>
# → canvas courses

# Complete subcommands
canvas courses <Tab>
# → list  get  create  update  delete

# Complete flags
canvas courses list --<Tab>
# → --output  --no-cache  --instance  --help
```

## Troubleshooting

### Completions Not Working

1. Ensure the completion script is sourced in your shell config
2. Restart your shell or source the config file
3. Check that Canvas CLI is in your PATH

### Zsh: Command Not Found

If you see "command not found: compdef", add this before sourcing:

```bash
autoload -Uz compinit && compinit
```
