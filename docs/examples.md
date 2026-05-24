# Examples

## Deliver Filtered Quarantine Emails

Fetch all spam from quarantine, filter by sender and receiver, then deliver matching emails.

```yaml
name: deliver-filtered-spam
description: "Deliver spam emails matching sender and receiver filters"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: false

vars:
  sender:
    default: "*"
    prompt: "Sender filter (domain or email)"
  receiver:
    default: "*"
    prompt: "Receiver domain filter"

steps:
  - id: fetch-quarantine
    action: quarantine.list
    params:
      type: spam
    filters:
      sender: "{{ .vars.sender }}"
      receiver: "{{ .vars.receiver }}"
    output: spam_list

  - id: deliver-matched
    action: quarantine.deliver
    input: spam_list
    confirm: true
    on_error: continue

  - id: summary
    action: report.console
    input: deliver-matched
    params:
      template: "deliver-summary"
```

## Bulk Blacklist Addition

Load emails/domains from a file, deduplicate against existing entries, and add to a Who Object.

```yaml
name: bulk-blacklist
description: "Bulk add emails/domains to a Who Object from a file"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: false

vars:
  who_id:
    prompt: "Target Who Object ID"
    required: true
  mode:
    default: "auto"
    prompt: "Mode (auto|email|domain)"
  entries_file:
    prompt: "Path to email/domain list file"
    required: true

steps:
  - id: load-existing
    action: ruledb.who.list
    params:
      who_id: "{{ .vars.who_id }}"
    output: existing_entries

  - id: load-entries
    action: transform.extract
    params:
      file: "{{ .vars.entries_file }}"
    output: raw_entries

  - id: deduplicate
    action: transform.deduplicate
    input: raw_entries
    params:
      mode: "{{ .vars.mode }}"
    output: new_entries

  - id: add-entries
    action: ruledb.who.add
    params:
      who_id: "{{ .vars.who_id }}"
      mode: "{{ .vars.mode }}"
    input: new_entries
    on_error: continue

  - id: summary
    action: report.console
    input: add-entries
    params:
      template: "blacklist-summary"
```

## Quarantine Statistics

Show a summary of quarantined emails.

```yaml
name: quarantine-stats
description: "Show quarantine statistics"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: false

steps:
  - id: list-spam
    action: quarantine.list
    params:
      type: spam
    output: spam_mails

  - id: list-virus
    action: quarantine.list
    params:
      type: virus
    output: virus_mails

  - id: summary
    action: report.console
    input: spam_mails
    params:
      template: "deliver-summary"
```

## Delete All Virus Quarantine (with confirmation)

Requires user confirmation before deleting.

```yaml
name: delete-virus
description: "Delete all virus-quarantined emails"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: false

steps:
  - id: fetch-virus
    action: quarantine.list
    params:
      type: virus
    output: virus_list

  - id: delete-all
    action: quarantine.delete
    input: virus_list
    confirm: true
    on_error: continue

  - id: summary
    action: report.console
    input: delete-all
    params:
      template: "deliver-summary"
```

## Remote PMG with Environment Variables

Connect to a remote PMG server using environment variables for credentials.

```yaml
name: remote-deliver
description: "Deliver spam from a remote PMG instance"
version: "1.0"

auth:
  host: "https://pmg.example.com:8006"
  token: "{{ .PMG_TOKEN }}"
  insecure: false

vars:
  sender:
    default: "*"
    prompt: "Sender filter"

steps:
  - id: fetch-quarantine
    action: quarantine.list
    params:
      type: spam
    filters:
      sender: "{{ .vars.sender }}"
    output: spam_list

  - id: deliver-matched
    action: quarantine.deliver
    input: spam_list
    on_error: continue

  - id: summary
    action: report.console
    input: deliver-matched
    params:
      template: "deliver-summary"
```

Run with:

```bash
export PMG_TOKEN=your-api-token-here
pmgwire apply remote-deliver.yaml
```