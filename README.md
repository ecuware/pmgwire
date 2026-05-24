# PMGWire

A declarative workflow engine for Proxmox Mail Gateway. Define automation tasks in YAML, run them from the CLI or an interactive TUI.

## Why

PMG ships with `pmgsh`, a command-line tool for individual API calls. But real operations are repetitive — bulk quarantine delivery, mass blacklist additions, statistical reports. These require chaining multiple calls, filtering results, and handling errors. PMGWire turns these into reproducible, shareable YAML workflows.

## How It Works

```
YAML Workflow → Parser → Executor → PMG REST API → PMG
                            ↓
                      Report / TUI Output
```

1. You write a YAML file describing steps (fetch quarantine, filter, deliver, report).
2. `pmgwire apply workflow.yaml` parses and executes each step sequentially.
3. Steps can reference previous step outputs, use variables, and handle errors.

## Installation

```bash
go build -o pmgwire ./cmd/pmgwire/
```

Or download a pre-built binary from Releases.

## Quick Start

### Apply a built-in workflow

```bash
pmgwire apply workflows/builtin/deliver-spam.yaml
```

### With interactive TUI

```bash
pmgwire apply workflows/builtin/deliver-spam.yaml --tui
# or
pmgwire tui workflows/builtin/deliver-spam.yaml
```

### Dry run (no changes)

```bash
pmgwire apply workflows/builtin/deliver-spam.yaml --dry-run
```

### Validate a workflow

```bash
pmgwire validate my-workflow.yaml
```

### Create a new workflow

```bash
pmgwire init my-task
```

### List available actions

```bash
pmgwire list
```

## Workflow Reference

```yaml
name: my-workflow
description: "What this does"
version: "1.0"

auth:                              # Optional. Defaults to localhost:8006
  host: "https://pmg.example.com"
  token: "{{ .PMG_TOKEN }}"         # Env var reference
  insecure: false

vars:                              # Interactive or default variables
  sender:
    default: "*"
    prompt: "Sender filter"
    required: false

steps:
  - id: fetch                      # Unique step ID
    action: quarantine.list        # Built-in action
    params:
      type: spam
    filters:
      sender: "{{ .vars.sender }}" # Template reference
    output: spam_list              # Output name for chaining

  - id: deliver
    action: quarantine.deliver
    input: spam_list               # References previous step output
    confirm: true                  # Ask before executing
    on_error: continue             # stop | continue | retry

  - id: summary
    action: report.console
    input: deliver
    params:
      template: "deliver-summary"
```

## Built-in Actions

| Action | Description |
|--------|-------------|
| `quarantine.list` | List spam/virus quarantine with filters |
| `quarantine.deliver` | Bulk deliver quarantined mails |
| `quarantine.delete` | Bulk delete quarantined mails |
| `ruledb.who.list` | List emails and domains in a Who Object |
| `ruledb.who.add` | Add emails/domains to a Who Object |
| `ruledb.who.remove` | Remove emails/domains from a Who Object |
| `ruledb.what.list` | List What Objects and their content |
| `ruledb.what.add` | Add content types, patterns, or fields to a What Object |
| `ruledb.what.remove` | Remove entries from or delete a What Object |
| `ruledb.rule.list` | List mail filter rules |
| `ruledb.rule.create` | Create a new mail filter rule |
| `ruledb.rule.remove` | Delete a mail filter rule |
| `transform.deduplicate` | Remove duplicate entries |
| `transform.filter` | Filter entries by pattern |
| `transform.extract` | Load entries from a file |
| `report.console` | Print formatted summary to terminal |
| `report.file` | Write report to file (JSON/CSV/YAML) |
| `report.json` | Output raw JSON |

## Project Structure

```
pmgwire/
├── cmd/pmgwire/main.go            # CLI entry point
├── internal/
│   ├── engine/
│   │   ├── parser.go               # YAML parsing, variable resolution, templates
│   │   └── executor.go             # Step execution, retry, confirm, dry-run
│   ├── pmg/
│   │   ├── client.go               # HTTP client, auth, request handling
│   │   ├── quarantine.go           # Quarantine API + filtering
│   │   ├── ruledb.go               # Rule database API
│   │   ├── config.go               # Config API
│   │   ├── statistics.go           # Statistics API
│   │   └── nodes.go                # Node management API
│   ├── actions/
│   │   ├── registry.go             # Action interface and global registry
│   │   ├── quarantine.go           # quarantine.* actions
│   │   ├── ruledb.go               # ruledb.who.* actions
│   │   ├── transform.go            # transform.* actions
│   │   └── report.go               # report.* actions
│   ├── tui/
│   │   ├── app.go                  # Bubble Tea TUI application model
│   │   ├── theme.go                # Lipgloss styles, colors, icons
│   │   └── components/             # Reusable TUI components
│   │       ├── table.go            # Interactive data table
│   │       ├── confirm.go          # Yes/No confirmation dialog
│   │       ├── form.go             # Variable input form
│   │       ├── progress.go         # Progress bar
│   │       └── summary.go          # Result summary panel
│   └── config/config.go            # App configuration
├── workflows/builtin/              # Ready-to-use workflow templates
│   ├── deliver-spam.yaml
│   ├── blacklist-bulk.yaml
│   └── quarantine-stats.yaml
├── go.mod
└── go.sum
```

## Auth

- **Local (default):** No auth config needed. Connects to `localhost:8006` using PMG's ticket API.
- **Remote:** Set `auth.host` and `auth.token` in YAML, or pass `--host` and `--token` flags.
- **Insecure:** Set `auth.insecure: true` or `--insecure` flag to skip TLS verification.

## Error Handling

Each step supports `on_error`:

- `stop` (default) — Halt on first error.
- `continue` — Log the error, proceed to next step.
- `retry` — Retry up to `retry_count` times (default 3).

Steps with `confirm: true` prompt for user approval before executing.

## License

MIT