# Workflow Reference

## Full Schema

```yaml
name: my-workflow              # Required. Workflow name
description: "What it does"    # Optional. Human-readable description
version: "1.0"                 # Required. Schema version

auth:                          # Optional. Defaults to localhost:8006
  host: "https://localhost:8006"
  token: "{{ .PMG_TOKEN }}"   # Optional. Env var reference or static token
  insecure: false              # Optional. Skip TLS verification

vars:                          # Optional. Variable definitions
  sender:
    default: "*"               # Default value when user presses Enter
    prompt: "Sender filter"    # Question shown to user
    required: false            # Fail if not provided
  who_id:
    prompt: "Target Who Object ID"
    required: true

steps:                         # Required. Ordered list of steps
  - id: fetch                  # Required. Unique step identifier
    action: quarantine.list    # Required. Action name
    params:                    # Optional. Action parameters
      type: spam
    filters:                   # Optional. Result filters (wildcards supported)
      sender: "{{ .vars.sender }}"
      receiver: "*@example.com"
    output: spam_list          # Optional. Store step result for later steps
    input: previous_step       # Optional. Read output from a previous step
    confirm: true              # Optional. Ask user before executing (default: false)
    on_error: continue          # Optional. stop | continue | retry (default: stop)
    retry_count: 3              # Optional. Retry attempts (default: 3, only when on_error: retry)
```

## Auth

| Field | Default | Description |
|-------|---------|-------------|
| `host` | `https://localhost:8006` | PMG server URL |
| `token` | (empty) | API token. When empty, connects locally. |
| `insecure` | `false` | Skip TLS certificate verification |

Overrides via CLI:

```bash
pmgwire apply workflow.yaml --host https://pmg.example.com --token MYTOKEN --insecure
```

## Variables

Variables are resolved in this order (highest priority first):

1. **Environment variable** — `PMGWIRE_<VARNAME>` (uppercase, dashes to underscores)
2. **Interactive prompt** — if `prompt` is set, user is asked at runtime
3. **Default value** — from the `default` field

Example:

```yaml
vars:
  sender_filter:
    default: "*"
    prompt: "Sender filter"
```

```
$ pmgwire apply workflow.yaml
Sender filter (default: *): info@example.com
```

Or via environment:

```bash
export PMGWIRE_SENDER_FILTER="info@example.com"
pmgwire apply workflow.yaml  # No prompt, uses env var
```

## Template Syntax

Step parameters, filters, and inputs support Go template syntax:

```yaml
steps:
  - id: fetch
    action: quarantine.list
    params:
      type: "{{ .vars.mail_type }}"
    filters:
      sender: "{{ .vars.sender }}"
```

Available template variables:

| Variable | Source |
|----------|--------|
| `{{ .vars.<name> }}` | Workflow vars |
| `{{ .PMG_TOKEN }}` | Environment variables starting with `PMG_` |
| `{{ .steps.<id>.<field> }}` | Output from previous steps |

## Step Chaining

Steps reference each other's outputs using `input` and `output`:

```yaml
steps:
  - id: fetch-quarantine
    action: quarantine.list
    params:
      type: spam
    output: spam_list          # Store result

  - id: deliver-matched
    action: quarantine.deliver
    input: spam_list           # Use previous step's output
    confirm: true
```

## Error Handling

| Value | Behavior |
|-------|----------|
| `stop` | Halt immediately on error (default) |
| `continue` | Log error, proceed to next step |
| `retry` | Retry up to `retry_count` times, then stop |

```yaml
steps:
  - id: risky-operation
    action: quarantine.delete
    input: mail_list
    on_error: retry
    retry_count: 5
```

## Confirmation

Steps with `confirm: true` prompt the user before executing. Useful for destructive operations:

```yaml
steps:
  - id: delete-all
    action: quarantine.delete
    input: spam_list
    confirm: true
```