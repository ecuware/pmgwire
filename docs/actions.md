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

## Rule Database — Who Objects

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

## Rule Database — What Objects

### ruledb.what.list

List What Objects or show details of a specific one.

```yaml
# List all What Objects
- id: list-all
  action: ruledb.what.list
  output: what_objects

# Show details of a specific What Object
- id: list-detail
  action: ruledb.what.list
  params:
    what_id: 3               # Optional. Show content types, patterns, fields
  output: what_detail
```

Output fields (no what_id): `what_objects` (list), `count` (int)
Output fields (with what_id): `what_id` (int), `content_types` (list), `patterns` (list), `fields` (list)

### ruledb.what.add

Add content types, patterns, or fields to a What Object. Or create a new one.

```yaml
# Add entries to existing What Object
- id: add-patterns
  action: ruledb.what.add
  params:
    what_id: 3
    mode: pattern              # contenttype | pattern | field | auto
  input: raw_entries
  on_error: continue

# Create a new What Object
- id: create-what
  action: ruledb.what.add
  params:
    create: true
    name: "Blocked Attachments"
```

Mode behavior:
- **contenttype** — Add as MIME content type (e.g. `application/zip`)
- **pattern** — Add as filename pattern (e.g. `*.exe`)
- **field** — Add as header field name (e.g. `X-Spam-Flag`)
- **auto** — Auto-detect: `/` → content type, `X-` prefix → field, else → pattern

Output fields: `what_id` (int), `added_content_types` (list), `added_patterns` (list), `added_fields` (list), `success` (int), `failed` (int)

### ruledb.what.remove

Remove entries from a What Object or delete it entirely.

```yaml
# Remove specific patterns
- id: remove-patterns
  action: ruledb.what.remove
  params:
    what_id: 3
    mode: pattern              # contenttype | pattern | field
  input: entries

# Delete entire What Object
- id: delete-what
  action: ruledb.what.remove
  params:
    what_id: 3
    delete: true
```

Output fields: `what_id` (int), `success` (int), `failed` (int)

---

## Rule Database — Rules

### ruledb.rule.list

List all rules or show details of a specific rule.

```yaml
# List all rules
- id: list-rules
  action: ruledb.rule.list
  output: rules

# Show single rule details
- id: list-rule
  action: ruledb.rule.list
  params:
    rule_id: 5               # Optional. Show rule details
  output: rule_detail
```

Output fields (no rule_id): `rules` (list), `count` (int)
Output fields (with rule_id): `rule_id` (int), `raw` (JSON)

### ruledb.rule.create

Create a new mail filter rule.

```yaml
- id: create-rule
  action: ruledb.rule.create
  params:
    priority: 100            # Rule priority (lower = higher priority)
    direction: in             # in | out | both
    action: block             # block | accept | quarantine
    who_id: 2                # Who Object ID
    what_id: 3                # What Object ID
    name: "Block Executables" # Optional rule name
  output: new_rule
```

Output fields: `priority` (int), `direction` (string), `action` (string), `who_id` (int), `what_id` (int), `raw` (JSON)

### ruledb.rule.remove

Delete a rule.

```yaml
- id: remove-rule
  action: ruledb.rule.remove
  params:
    rule_id: 5
```

Output fields: `rule_id` (int)

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