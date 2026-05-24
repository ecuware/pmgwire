# Actions

All available actions and their parameters.

## Quarantine

### quarantine.list

List quarantined emails with optional filters.

```yaml
- id: fetch
  action: quarantine.list
  params:
    type: spam               # spam | virus
  filters:
    sender: "*@example.com"  # Optional. Wildcard supported
    receiver: "admin@"      # Optional. Wildcard supported
  output: mail_list
```

Output fields: `count` (int), `mails` (list), `ids` (list of strings)

### quarantine.deliver

Deliver quarantined emails by ID.

```yaml
- id: deliver
  action: quarantine.deliver
  input: fetch               # Requires step with 'ids' in output
```

Output fields: `delivered` (int), `failed` (int), `total` (int)

### quarantine.delete

Delete quarantined emails by ID.

```yaml
- id: delete
  action: quarantine.delete
  input: fetch
```

Output fields: `deleted` (int), `failed` (int), `total` (int)

---

## Rule Database

### ruledb.who.list

List emails and domains in a Who Object.

```yaml
- id: load-existing
  action: ruledb.who.list
  params:
    who_id: 2               # Required. Who Object ID
  output: existing
```

Output fields: `who_id` (int), `emails` (list of strings), `domains` (list of strings)

### ruledb.who.add

Add emails or domains to a Who Object.

```yaml
- id: add-entries
  action: ruledb.who.add
  params:
    who_id: 2
    mode: auto               # auto | email | domain
  input: deduplicated_entries
```

Mode behavior:
- **auto** — Automatically detect whether each entry is an email or domain
- **email** — Treat all entries as emails
- **domain** — Treat all entries as domains (strips leading `@`)

Output fields: `who_id` (int), `added_emails` (list), `added_domains` (list), `success` (int), `failed` (int)

### ruledb.who.remove

Remove emails from a Who Object.

```yaml
- id: remove-entries
  action: ruledb.who.remove
  params:
    who_id: 2
  input: entries_with_emails
```

Output fields: `who_id` (int)

---

## Transform

Client-side data transformations. No API calls.

### transform.deduplicate

Remove duplicate entries from a list.

```yaml
- id: deduplicate
  action: transform.deduplicate
  input: raw_entries
  params:
    mode: auto               # auto | email | domain
  output: unique_entries
```

Output fields: `entries` (list), `duplicates` (int)

### transform.filter

Filter entries by pattern.

```yaml
- id: filter-sender
  action: transform.filter
  input: all_entries
  params:
    pattern: "example.com"
  output: matched
```

Output fields: `entries` (list), `removed` (int)

### transform.extract

Load entries from a file (one per line).

```yaml
- id: load-entries
  action: transform.extract
  params:
    file: "/path/to/emails.txt"
  output: raw_entries
```

Output fields: `entries` (list)

---

## Report

### report.console

Print a formatted summary to the terminal.

```yaml
- id: summary
  action: report.console
  input: deliver-step
  params:
    template: "deliver-summary"   # deliver-summary | blacklist-summary | (empty for JSON)
```

### report.file

Write report to a file.

```yaml
- id: save-report
  action: report.file
  input: fetch-step
  params:
    path: "/tmp/report.json"
    format: json               # json (only format currently)
```

### report.json

Output raw JSON to stdout.

```yaml
- id: json-out
  action: report.json
  input: fetch-step
```