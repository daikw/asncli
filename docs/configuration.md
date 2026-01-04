# Configuration

## Config File

Location:

- **macOS**: `~/Library/Application Support/asncli/config.json`
- **Linux**: `~/.config/asncli/config.json`
- **Windows**: `%APPDATA%\asncli\config.json`

Format:

```json
{
  "default_workspace": "1234567890",
  "default_workspace_name": "My Workspace"
}
```

The file is created automatically by `asn config set-workspace`.

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `ASNCLI_TOKEN` | Personal access token |
| `ASNCLI_DEFAULT_WORKSPACE` | Default workspace GID |

## Resolution Order

### Token
1. `ASNCLI_TOKEN` environment variable
2. Keyring/keychain stored token (via `asn auth login`)

### Workspace
1. `--workspace` flag
2. `ASNCLI_DEFAULT_WORKSPACE` environment variable
3. Config file `default_workspace`
4. Error if none found

## Commands

```bash
# Set default workspace (interactive)
asn config set-workspace

# Show current default workspace
asn config get-workspace

# Show all configuration
asn config show
```

## Examples

Interactive use:

```bash
asn auth login
asn config set-workspace

# Now these work without --workspace
asn tasks search --text "bug"
asn projects list
```

CI/automation:

```bash
export ASNCLI_TOKEN="your-token"
export ASNCLI_DEFAULT_WORKSPACE="123456789"

asn tasks search --text "bug" --json
```

Override default:

```bash
asn tasks search --workspace 987654321 --text "feature"
```
