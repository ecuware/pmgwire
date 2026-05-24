package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ecuware/pmgwire/internal/pmg"
)

type ReportConsoleAction struct{}

func (a *ReportConsoleAction) Name() string { return "report.console" }

func (a *ReportConsoleAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	template := getStringParam(params, "template")

	fmt.Println("\n--- Report ---")

	switch template {
	case "deliver-summary":
		printDeliverSummary(input)
	case "blacklist-summary":
		printBlacklistSummary(input)
	default:
		printGenericReport(input)
	}

	fmt.Println("--- End ---")

	return input, nil
}

func printDeliverSummary(input Data) {
	if total, ok := input["total"]; ok {
		fmt.Printf("  Total: %v\n", total)
	}
	if delivered, ok := input["delivered"]; ok {
		fmt.Printf("  Delivered: %v\n", delivered)
	}
	if failed, ok := input["failed"]; ok {
		fmt.Printf("  Failed: %v\n", failed)
	}
	if count, ok := input["count"]; ok {
		fmt.Printf("  Mails found: %v\n", count)
	}
}

func printBlacklistSummary(input Data) {
	fmt.Printf("  Who Object ID: %v\n", input["who_id"])
	if success, ok := input["success"]; ok {
		fmt.Printf("  Successfully added: %v\n", success)
	}
	if failed, ok := input["failed"]; ok {
		fmt.Printf("  Failed: %v\n", failed)
	}
	if emails, ok := input["added_emails"]; ok {
		fmt.Printf("  Added emails: %v\n", emails)
	}
	if domains, ok := input["added_domains"]; ok {
		fmt.Printf("  Added domains: %v\n", domains)
	}
}

func printGenericReport(input Data) {
	data, err := json.MarshalIndent(input, "  ", "  ")
	if err != nil {
		fmt.Printf("  %v\n", input)
		return
	}
	fmt.Printf("  %s\n", string(data))
}

type ReportFileAction struct{}

func (a *ReportFileAction) Name() string { return "report.file" }

func (a *ReportFileAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	filePath := getStringParam(params, "path")
	format := getStringParam(params, "format")
	if format == "" {
		format = "json"
	}

	var content []byte
	var err error

	switch format {
	case "json":
		content, err = json.MarshalIndent(input, "", "  ")
	default:
		content, err = json.MarshalIndent(input, "", "  ")
	}

	if err != nil {
		return nil, fmt.Errorf("marshaling report: %w", err)
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return nil, fmt.Errorf("writing report file: %w", err)
	}

	fmt.Printf("  Report written to %s\n", filePath)
	return input, nil
}

type ReportJSONAction struct{}

func (a *ReportJSONAction) Name() string { return "report.json" }

func (a *ReportJSONAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling JSON: %w", err)
	}

	fmt.Println(string(data))
	return input, nil
}

func init() {
	Register(&ReportConsoleAction{})
	Register(&ReportFileAction{})
	Register(&ReportJSONAction{})
}

var _ = time.Time{}