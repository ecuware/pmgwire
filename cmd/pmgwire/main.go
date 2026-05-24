package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ecuware/pmgwire/internal/actions"
	"github.com/ecuware/pmgwire/internal/config"
	"github.com/ecuware/pmgwire/internal/engine"
	"github.com/ecuware/pmgwire/internal/pmg"
	"github.com/ecuware/pmgwire/internal/tui"
	"github.com/spf13/cobra"
)

var (
	version      = "0.1.0"
	hostFlag     string
	tokenFlag    string
	insecureFlag bool
	dryRunFlag   bool
	tuiFlag      bool
)

func main() {
	initActions()

	rootCmd := &cobra.Command{
		Use:   "pmgwire",
		Short: "PMGWire - Proxmox Mail Gateway Workflow Engine",
		Long:  "A declarative workflow engine for Proxmox Mail Gateway. Automate PMG operations using YAML configuration files.",
	}

	rootCmd.PersistentFlags().StringVar(&hostFlag, "host", "", "PMG host URL (overrides workflow config)")
	rootCmd.PersistentFlags().StringVar(&tokenFlag, "token", "", "PMG API token (overrides workflow config)")
	rootCmd.PersistentFlags().BoolVar(&insecureFlag, "insecure", false, "Skip TLS certificate verification")

	applyCmd := &cobra.Command{
		Use:   "apply <workflow.yaml>",
		Short: "Execute a workflow",
		Args:  cobra.ExactArgs(1),
		RunE:  runApply,
	}
	applyCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Dry run mode (no changes)")
	applyCmd.Flags().BoolVar(&tuiFlag, "tui", false, "Run with interactive TUI")
	rootCmd.AddCommand(applyCmd)

	tuiCmd := &cobra.Command{
		Use:   "tui <workflow.yaml>",
		Short: "Run workflow with interactive TUI",
		Args:  cobra.ExactArgs(1),
		RunE:  runTUI,
	}
	rootCmd.AddCommand(tuiCmd)

	validateCmd := &cobra.Command{
		Use:   "validate <workflow.yaml>",
		Short: "Validate a workflow file",
		Args:  cobra.ExactArgs(1),
		RunE:  runValidate,
	}
	rootCmd.AddCommand(validateCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available workflows and actions",
		RunE:  runList,
	}
	rootCmd.AddCommand(listCmd)

	initCmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Create a new workflow template",
		Args:  cobra.ExactArgs(1),
		RunE:  runInit,
	}
	rootCmd.AddCommand(initCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("pmgwire v%s\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initActions() {
	_ = actions.All()
}

func createClient(wf *engine.Workflow) (*pmg.Client, error) {
	auth := pmg.AuthConfig{
		Host:     wf.Auth.Host,
		Token:    wf.Auth.Token,
		Insecure: wf.Auth.Insecure,
	}

	if hostFlag != "" {
		auth.Host = hostFlag
	}
	if tokenFlag != "" {
		auth.Token = tokenFlag
	}
	if insecureFlag {
		auth.Insecure = true
	}

	return pmg.NewClient(auth)
}

func runApply(cmd *cobra.Command, args []string) error {
	wf, err := engine.ParseWorkflow(args[0])
	if err != nil {
		return fmt.Errorf("parsing workflow: %w", err)
	}

	if tuiFlag {
		return runWithTUI(wf)
	}

	return runWithCLI(wf)
}

func runTUI(cmd *cobra.Command, args []string) error {
	wf, err := engine.ParseWorkflow(args[0])
	if err != nil {
		return fmt.Errorf("parsing workflow: %w", err)
	}
	return runWithTUI(wf)
}

func runWithTUI(wf *engine.Workflow) error {
	client, err := createClient(wf)
	if err != nil {
		return fmt.Errorf("creating PMG client: %w", err)
	}

	vars, err := engine.ResolveVars(wf)
	if err != nil {
		return fmt.Errorf("resolving variables: %w", err)
	}

	stepOutputs := make(map[string]map[string]interface{})
	if err := engine.ResolveTemplates(wf, vars, stepOutputs); err != nil {
		return fmt.Errorf("resolving templates: %w", err)
	}

	steps := make([]tui.EngineStep, len(wf.Steps))
	for i, s := range wf.Steps {
		steps[i] = tui.EngineStep{
			ID:       s.ID,
			Action:   s.Action,
			Params:   s.Params,
			Filters:  s.Filters,
			Input:    s.Input,
			Output:   s.Output,
			Confirm:  s.Confirm,
			OnError:  s.OnError,
			RetryCnt: s.RetryCount,
		}
	}

	return tui.RunApp(wf.Name, wf.Description, steps, client, vars, dryRunFlag)
}

func runWithCLI(wf *engine.Workflow) error {
	vars, err := engine.ResolveVars(wf)
	if err != nil {
		return fmt.Errorf("resolving variables: %w", err)
	}

	client, err := createClient(wf)
	if err != nil {
		return fmt.Errorf("creating PMG client: %w", err)
	}

	stepOutputs := make(map[string]map[string]interface{})
	if err := engine.ResolveTemplates(wf, vars, stepOutputs); err != nil {
		return fmt.Errorf("resolving templates: %w", err)
	}

	tui.PrintColoredHeader(wf.Name, wf.Version, wf.Description)

	if dryRunFlag {
		fmt.Println(tui.WarningStyle.Render("  ⚠ DRY RUN MODE - No changes will be made"))
		fmt.Println()
	}

	executor := engine.NewExecutor(client, vars, dryRunFlag)
	return executor.Execute(context.Background(), wf)
}

func runValidate(cmd *cobra.Command, args []string) error {
	wf, err := engine.ParseWorkflow(args[0])
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println(tui.TitleStyle.Render(fmt.Sprintf(" %s Workflow: %s ", tui.AppIcon, wf.Name)))
	fmt.Printf("  Version: %s\n", wf.Version)
	fmt.Printf("  Steps:   %d\n", len(wf.Steps))
	fmt.Println()

	for _, step := range wf.Steps {
		icon := tui.OKIcon
		if _, ok := actions.Get(step.Action); !ok {
			icon = tui.WarnIcon
		}
		fmt.Printf("  %s %s %s\n", tui.SuccessStyle.Render(icon), tui.BrandBold.Render(step.ID), tui.DimStyle.Render(step.Action))
	}

	availableActions := actions.All()
	for _, step := range wf.Steps {
		if _, ok := availableActions[step.Action]; !ok {
			fmt.Printf("\n  %s Unknown action '%s' in step '%s'\n", tui.ErrorStyle.Render(tui.FailIcon), step.Action, step.ID)
		}
	}

	fmt.Println()
	fmt.Println(tui.SuccessStyle.Render(fmt.Sprintf("  %s Workflow is valid!", tui.OKIcon)))
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	fmt.Println(tui.TitleStyle.Render(fmt.Sprintf(" %s Available Actions ", tui.AppIcon)))
	fmt.Println()
	for name := range actions.All() {
		fmt.Printf("  %s %s\n", tui.BrandNormal.Render(tui.ArrowIcon), name)
	}

	cfg := config.DefaultConfig()
	fmt.Println()
	fmt.Println(tui.TitleStyle.Render(fmt.Sprintf(" %s Workflow Directory ", tui.AppIcon)))
	fmt.Printf("  %s\n", cfg.WorkflowsDir)

	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	name := args[0]

	if strings.ContainsAny(name, "/\\^ \t\n\r") || strings.ContainsRune(name, 0) {
		return fmt.Errorf("invalid workflow name: %q", name)
	}

	cfg := config.DefaultConfig()
	if err := cfg.EnsureDirs(); err != nil {
		return fmt.Errorf("creating workflow directory: %w", err)
	}

	template := fmt.Sprintf(`name: %s
description: "Description of your workflow"
version: "1.0"

auth:
  host: "https://localhost:8006"
  insecure: false

vars: {}

steps:
  - id: example-step
    action: quarantine.list
    params:
      type: spam
    output: result

  - id: summary
    action: report.console
    input: example-step
`, name)

	filePath := fmt.Sprintf("%s/%s.yaml", cfg.WorkflowsDir, name)
	if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
		return fmt.Errorf("writing workflow template: %w", err)
	}

	fmt.Println(tui.SuccessStyle.Render(fmt.Sprintf("  %s Created workflow template:", tui.OKIcon)), filePath)
	return nil
}