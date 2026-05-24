# CLI Reference

## Global Flags

These flags apply to all commands.

| Flag | Type | Description |
|------|------|-------------|
| `--host` | string | PMG host URL (overrides workflow config) |
| `--token` | string | PMG API token (overrides workflow config) |
| `--insecure` | bool | Skip TLS certificate verification |

---

## pmgwire apply

Execute a workflow.

```bash
pmgwire apply <workflow.yaml> [flags]
```

| Flag | Type | Description |
|------|------|-------------|
| `--dry-run` | bool | Preview changes without executing them |
| `--tui` | bool | Run with interactive terminal UI |

Examples:

```bash
pmgwire apply deliver-spam.yaml
pmgwire apply deliver-spam.yaml --dry-run
pmgwire apply deliver-spam.yaml --tui
pmgwire apply deliver-spam.yaml --host https://pmg.example.com --token MYTOKEN
```

---

## pmgwire tui

Run a workflow with the interactive terminal interface. Equivalent to `pmgwire apply --tui`.

```bash
pmgwire tui <workflow.yaml>
```

---

## pmgwire validate

Validate a workflow YAML file. Checks for required fields, unique step IDs, template syntax, and action names. Does not connect to a PMG server.

```bash
pmgwire validate <workflow.yaml>
```

---

## pmgwire list

Display all registered actions grouped by category, and the local workflow directory path.

```bash
pmgwire list
```

---

## pmgwire init

Create a new workflow template in `~/.pmgwire/workflows/`.

```bash
pmgwire init <name>
```

The name must not contain spaces, slashes, or special characters.

After creation, edit the generated file and run:

```bash
pmgwire apply ~/.pmgwire/workflows/<name>.yaml
```

---

## pmgwire version

Print the current PMGWire version.

```bash
pmgwire version
```