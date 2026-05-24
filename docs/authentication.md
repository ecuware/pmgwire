# Authentication

PMGWire supports two authentication modes: local and remote.

## Local Authentication (Default)

When running on the PMG server itself, no auth configuration is needed:

```yaml
auth:
  host: "https://localhost:8006"
```

Or simply omit the `auth` block entirely — `localhost:8006` is the default.

In this mode, PMGWire connects using the local PMG ticket system.

## Remote Authentication

To connect to a remote PMG server, provide an API token:

```yaml
auth:
  host: "https://pmg.example.com:8006"
  token: "{{ .PMG_TOKEN }}"
```

Then set the environment variable:

```bash
export PMG_TOKEN=your-api-token
pmgwire apply workflow.yaml
```

You can also pass the token via CLI flag:

```bash
pmgwire apply workflow.yaml --token your-api-token
```

## Self-Signed Certificates

If your PMG server uses a self-signed certificate, enable insecure mode:

```yaml
auth:
  host: "https://pmg.example.com:8006"
  token: "{{ .PMG_TOKEN }}"
  insecure: true
```

Or via CLI:

```bash
pmgwire apply workflow.yaml --insecure
```

## Priority Order

Auth settings are resolved in this order (highest priority wins):

1. **CLI flags** — `--host`, `--token`, `--insecure`
2. **Environment variables** — `PMGWIRE_*` prefix
3. **Workflow YAML** — `auth` block
4. **Defaults** — `localhost:8006`, no token, secure TLS