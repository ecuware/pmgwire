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
		Use:   "pmgwire [command]",
		Short: "PMGWire — Declarative workflow engine for Proxmox Mail Gateway",
		Long: heredoc(`
			PMGWire — Declarative workflow engine for Proxmox Mail Gateway

			Automate PMG operations using YAML workflow definitions.
			Chain quarantine delivery, blacklist management, reporting
			and more into reproducible, shareable workflows.

			Quick start:
			  pmgwire init my-task              Create a new workflow template
			  pmgwire apply my-task.yaml        Execute a workflow
			  pmgwire apply my-task.yaml --tui  Execute with interactive TUI
			  pmgwire validate my-task.yaml     Check a workflow for errors

			Get started:
			  https://github.com/ecuware/pmgwire
		`),
	}

	rootCmd.PersistentFlags().StringVar(&hostFlag, "host", "", "PMG host URL (overrides workflow config)")
	rootCmd.PersistentFlags().StringVar(&tokenFlag, "token", "", "PMG API token (overrides workflow config)")
	rootCmd.PersistentFlags().BoolVar(&insecureFlag, "insecure", false, "Skip TLS certificate verification")

	applyCmd := &cobra.Command{
		Use:   "apply <workflow.yaml>",
		Short: "Execute a workflow",
		Long: heredoc(`
			Execute a workflow defined in a YAML file.

			The workflow is parsed, variables are resolved (from flags,
			environment variables, or interactive prompts), and each step
			is executed in order.

			Use --dry-run to preview what would happen without making changes.
			Use --tui to run with an interactive terminal interface.
		`),
		Example: heredoc(`
			pmgwire apply workflows/builtin/deliver-spam.yaml
			pmgwire apply my-workflow.yaml --dry-run
			pmgwire apply my-workflow.yaml --tui
			pmgwire apply my-workflow.yaml --host https://pmg.example.com --token MYTOKEN
		`),
		Args: cobra.ExactArgs(1),
		RunE: runApply,
	}
	applyCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Preview changes without executing them")
	applyCmd.Flags().BoolVar(&tuiFlag, "tui", false, "Run with interactive TUI")
	rootCmd.AddCommand(applyCmd)

	tuiCmd := &cobra.Command{
		Use:   "tui <workflow.yaml>",
		Short: "Run workflow with interactive TUI",
		Long: heredoc(`
			Run a workflow using the interactive terminal interface.

			Shows step-by-step progress with colored output, spinners,
			and summaries. Equivalent to: pmgwire apply <file> --tui
		`),
		Example: heredoc(`
			pmgwire tui workflows/builtin/deliver-spam.yaml
			pmgwire tui my-workflow.yaml
		`),
		Args: cobra.ExactArgs(1),
		RunE: runTUI,
	}
	rootCmd.AddCommand(tuiCmd)

	validateCmd := &cobra.Command{
		Use:   "validate <workflow.yaml>",
		Short: "Validate a workflow file",
		Long: heredoc(`
			Check a workflow YAML file for errors.

			Verifies: required fields, unique step IDs, known actions,
			and template syntax. Does not connect to a PMG server.
		`),
		Example: heredoc(`
			pmgwire validate my-workflow.yaml
		`),
		Args: cobra.ExactArgs(1),
		RunE: runValidate,
	}
	rootCmd.AddCommand(validateCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available actions and workflow directory",
		Long: heredoc(`
			Display all registered actions that can be used in workflow
			steps, along with the directory where local workflows are stored.
		`),
		Example: heredoc(`
			pmgwire list
		`),
		RunE: runList,
	}
	rootCmd.AddCommand(listCmd)

	initCmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Create a new workflow template",
		Long: heredoc(`
			Create a new workflow template file in the workflow directory.

			A skeleton YAML file is generated with the most common
			fields pre-filled. Edit it to define your workflow.
		`),
		Example: heredoc(`
			pmgwire init deliver-monthly
			pmgwire init sync-blacklist
		`),
		Args: cobra.ExactArgs(1),
		RunE: runInit,
	}
	rootCmd.AddCommand(initCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Display the PMGWire version number.",
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

	actionGroups := map[string][]string{
		"Quarantine": {"quarantine.list", "quarantine.deliver", "quarantine.delete"},
		"Rule Database": {"ruledb.who.list", "ruledb.who.add", "ruledb.who.remove"},
		"Transform":   {"transform.deduplicate", "transform.filter", "transform.extract"},
		"Report":      {"report.console", "report.file", "report.json"},
	}

	registered := actions.All()
	for group, actionNames := range actionGroups {
		fmt.Printf("  %s\n", tui.BrandBold.Render(group))
		for _, name := range actionNames {
			icon := "  "
			if _, ok := registered[name]; ok {
				icon = tui.SuccessStyle.Render(tui.OKIcon + " ")
			} else {
				icon = tui.DimStyle.Render("○ ")
			}
			fmt.Printf("    %s %s\n", icon, name)
		}
		fmt.Println()
	}

	cfg := config.DefaultConfig()
	fmt.Println(tui.TitleStyle.Render(fmt.Sprintf(" %s Workflow Directory ", tui.AppIcon)))
	fmt.Printf("  %s\n\n", cfg.WorkflowsDir)

	fmt.Println(tui.DimStyle.Render("  Use 'pmgwire init <name>' to create a new workflow template."))
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
	fmt.Println()
	fmt.Println(tui.DimStyle.Render("  Edit the file and run:"))
	fmt.Printf("  pmgwire apply %s\n", filePath)
	return nil
}

func heredoc(s string) string {
	s = strings.TrimPrefix(s, "\n")
	lines := strings.Split(s, "\n")
	var trimmed []string
	for _, line := range lines {
		if strings.HasPrefix(line, "\t\t") {
			trimmed = append(trimmed, line[2:])
		} else {
			trimmed = append(trimmed, line)
		}
	}
	return strings.Join(trimmed, "\n")
}