# asncli

Command-line interface for [Asana](https://asana.com/). Human-readable output by default, `--json` for automation and AI coding agents.

## Installation

Quick install

```bash
go install github.com/michalvavra/asncli/cmd/asn@latest
```

Or from source

```bash
git clone https://github.com/michalvavra/asncli.git
cd asncli
go install ./cmd/asn
```

Make sure `$GOPATH/bin` (typically `$HOME/go/bin`) is in your `PATH`.

## Setup

1. Create a Personal Access Token in Asana:
   - Go to Settings -> Apps -> Developer console
   - Create new token

2. Store your token:

```bash
asn auth login
```

3. (Optional) Set a default workspace:

```bash
asn config set-workspace
```

## Usage

```bash
# List and search tasks
asn tasks list --assignee=me
asn tasks list --project <project-gid>
asn tasks search --text "bug"

# Manage tasks
asn tasks get <task-gid>
asn tasks create --name "Fix login bug" --project <project-gid>
asn tasks update <task-gid> --completed=true

# JSON output for any command
asn auth status --json
```

## Configuration

Configuration is stored in:
- **macOS**: `~/Library/Application Support/asncli/config.json`
- **Linux**: `~/.config/asncli/config.json`
- **Windows**: `%APPDATA%\asncli\config.json`

### Set Default Workspace

```bash
# Interactive selection from your workspaces
asn config set-workspace

# View current default
asn config get-workspace

# Show all configuration
asn config show
```

### Environment Variables (Optional)

- `ASNCLI_TOKEN`: Personal access token (takes precedence over stored token)
- `ASNCLI_DEFAULT_WORKSPACE`: Default workspace GID

Useful for CI/CD or scripting:

```bash
export ASNCLI_TOKEN="your-token"
export ASNCLI_DEFAULT_WORKSPACE="123456789"
asn tasks search --text "bug" --json
```

## Development

```bash
# Clone and build
git clone https://github.com/michalvavra/asncli.git
cd asncli
go build -o bin/asn ./cmd/asn

# Run locally
./bin/asn --help

# Run tests
go test ./...
```

## References

Asana API: https://developers.asana.com/docs

OpenAPI spec: https://raw.githubusercontent.com/Asana/openapi/master/defs/app_components_oas.yaml
