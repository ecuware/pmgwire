# Getting Started

## Installation

Download the latest binary for your platform from [Releases](https://github.com/ecuware/pmgwire/releases).

```bash
# Linux amd64
chmod +x pmgwire-linux-amd64
sudo mv pmgwire-linux-amd64 /usr/local/bin/pmgwire

# macOS (Apple Silicon)
chmod +x pmgwire-darwin-arm64
sudo mv pmgwire-darwin-arm64 /usr/local/bin/pmgwire
```

Or build from source:

```bash
git clone https://github.com/ecuware/pmgwire.git
cd pmgwire
go build -o pmgwire ./cmd/pmgwire/
```

## First Workflow

### 1. Create a workflow template

```bash
pmgwire init my-first-workflow
```

This creates `~/.pmgwire/workflows/my-first-workflow.yaml`:

```yaml
name: my-first-workflow
description: "Description of your workflow"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: false

vars: {}

steps:
  - id: example-step
    action: quarantine.list
    params:
      type: spam
    output: result

  - id: summary
    action: report.console
    input: example-step
```

### 2. Edit the workflow

Open the file and customize it. See [Workflow Reference](workflow-reference.md) for all options.

### 3. Validate

```bash
pmgwire validate ~/.pmgwire/workflows/my-first-workflow.yaml
```

### 4. Run

```bash
# Dry run first (no changes)
pmgwire apply ~/.pmgwire/workflows/my-first-workflow.yaml --dry-run

# Execute for real
pmgwire apply ~/.pmgwire/workflows/my-first-workflow.yaml

# With interactive TUI
pmgwire apply ~/.pmgwire/workflows/my-first-workflow.yaml --tui
```

## Built-in Workflow Templates

PMGWire ships with ready-to-use workflows in `workflows/builtin/`:

| Template | Description |
|----------|-------------|
| `deliver-spam.yaml` | Filter and deliver quarantined spam emails |
| `blacklist-bulk.yaml` | Bulk add emails/domains to Who Objects |
| `quarantine-stats.yaml` | Show quarantine statistics |

## Connecting to a Remote PMG

By default, PMGWire connects to `localhost:8006`. To connect to a remote server:

```bash
# Via command-line flags
pmgwire apply workflow.yaml --host https://pmg.example.com --token YOUR_API_TOKEN

# Via workflow YAML
auth:
  host: "https://pmg.example.com"
  token: "{{ .PMG_TOKEN }}"
  insecure: false
```

Then run:

```bash
export PMG_TOKEN=your-token-here
pmgwire apply workflow.yaml
```

## Next Steps

- [Workflow Reference](workflow-reference.md) — Full YAML schema and options
- [Actions](actions.md) — All available actions and their parameters
- [Examples](examples.md) — Ready-to-use workflow examples