package actions

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ecuware/pmgwire/internal/pmg"
)

type TransformDeduplicateAction struct{}

func (a *TransformDeduplicateAction) Name() string { return "transform.deduplicate" }

func (a *TransformDeduplicateAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	mode := getStringParam(params, "mode")
	if mode == "" {
		mode = "auto"
	}

	entries := getEntriesFromInput(input)
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries provided for deduplication")
	}

	seen := make(map[string]bool)
	var unique []string
	for _, entry := range entries {
		normalized := normalizeEntry(entry, mode)
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		unique = append(unique, normalized)
	}

	fmt.Printf("  Deduplication: %d entries -> %d unique\n", len(entries), len(unique))

	return Data{
		"entries":    unique,
		"duplicates": len(entries) - len(unique),
	}, nil
}

type TransformFilterAction struct{}

func (a *TransformFilterAction) Name() string { return "transform.filter" }

func (a *TransformFilterAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	entries := getEntriesFromInput(input)
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries provided for filtering")
	}

	pattern := getStringParam(params, "pattern")

	var filtered []string
	for _, entry := range entries {
		if pattern == "" || pattern == "*" || strings.Contains(strings.ToLower(entry), strings.ToLower(pattern)) {
			filtered = append(filtered, entry)
		}
	}

	fmt.Printf("  Filter: %d entries -> %d matched\n", len(entries), len(filtered))

	return Data{
		"entries": filtered,
		"removed": len(entries) - len(filtered),
	}, nil
}

type TransformExtractAction struct{}

func (a *TransformExtractAction) Name() string { return "transform.extract" }

func (a *TransformExtractAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	filePath := getStringParam(params, "file")
	if filePath == "" {
		return nil, fmt.Errorf("file parameter is required for extract")
	}

	entries, err := readEntriesFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading entries from file: %w", err)
	}

	fmt.Printf("  Extracted %d entries from %s\n", len(entries), filePath)

	return Data{
		"entries": entries,
	}, nil
}

func normalizeEntry(entry, mode string) string {
	entry = strings.TrimSpace(entry)
	switch mode {
	case "email":
		if !strings.Contains(entry, "@") {
			return ""
		}
		return entry
	case "domain":
		entry = strings.TrimPrefix(entry, "@")
		if !strings.Contains(entry, ".") {
			return ""
		}
		return entry
	default:
		return entry
	}
}

func readEntriesFromFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			entries = append(entries, line)
		}
	}
	return entries, nil
}

func init() {
	Register(&TransformDeduplicateAction{})
	Register(&TransformFilterAction{})
	Register(&TransformExtractAction{})
}