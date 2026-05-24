package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseWorkflow(t *testing.T) {
	yamlContent := `name: test-workflow
description: "A test workflow"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: true

vars:
  sender:
    default: "*"
    prompt: "Sender filter"

steps:
  - id: fetch
    action: quarantine.list
    params:
      type: spam
    filters:
      sender: "{{ .vars.sender }}"
    output: result

  - id: summary
    action: report.console
    input: fetch
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	wf, err := ParseWorkflow(tmpFile)
	if err != nil {
		t.Fatalf("ParseWorkflow() error = %v", err)
	}

	if wf.Name != "test-workflow" {
		t.Errorf("Name = %q, want %q", wf.Name, "test-workflow")
	}
	if wf.Version != "1.0" {
		t.Errorf("Version = %q, want %q", wf.Version, "1.0")
	}
	if len(wf.Steps) != 2 {
		t.Errorf("Steps count = %d, want 2", len(wf.Steps))
	}
	if wf.Steps[0].ID != "fetch" {
		t.Errorf("Step 0 ID = %q, want %q", wf.Steps[0].ID, "fetch")
	}
	if wf.Steps[0].Action != "quarantine.list" {
		t.Errorf("Step 0 Action = %q, want %q", wf.Steps[0].Action, "quarantine.list")
	}
	if wf.Steps[0].OnError != "stop" {
		t.Errorf("Step 0 OnError = %q, want %q (default)", wf.Steps[0].OnError, "stop")
	}
	if wf.Auth.Insecure != true {
		t.Errorf("Auth.Insecure = %v, want true", wf.Auth.Insecure)
	}
}

func TestParseWorkflowMissingName(t *testing.T) {
	yamlContent := `version: "1.0"
steps:
  - id: test
    action: quarantine.list
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "bad.yaml")
	os.WriteFile(tmpFile, []byte(yamlContent), 0644)

	_, err := ParseWorkflow(tmpFile)
	if err == nil {
		t.Error("expected error for missing name, got nil")
	}
}

func TestParseWorkflowMissingSteps(t *testing.T) {
	yamlContent := `name: empty
version: "1.0"
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.yaml")
	os.WriteFile(tmpFile, []byte(yamlContent), 0644)

	_, err := ParseWorkflow(tmpFile)
	if err == nil {
		t.Error("expected error for missing steps, got nil")
	}
}

func TestParseWorkflowDefaults(t *testing.T) {
	yamlContent := `name: minimal
version: "1.0"
steps:
  - id: step1
    action: quarantine.list
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "minimal.yaml")
	os.WriteFile(tmpFile, []byte(yamlContent), 0644)

	wf, err := ParseWorkflow(tmpFile)
	if err != nil {
		t.Fatalf("ParseWorkflow() error = %v", err)
	}

	if wf.Auth.Host != "https://localhost:8006" {
		t.Errorf("Default host = %q, want %q", wf.Auth.Host, "https://localhost:8006")
	}
	if wf.Steps[0].OnError != "stop" {
		t.Errorf("Default OnError = %q, want %q", wf.Steps[0].OnError, "stop")
	}
}