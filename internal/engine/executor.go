package engine

import (
	"context"
	"fmt"

	"github.com/ecuware/pmgwire/internal/actions"
	"github.com/ecuware/pmgwire/internal/pmg"
	tui "github.com/ecuware/pmgwire/internal/tui"
)

type Executor struct {
	Client      *pmg.Client
	StepOutputs map[string]actions.Data
	Vars        map[string]interface{}
	DryRun      bool
}

func NewExecutor(client *pmg.Client, vars map[string]interface{}, dryRun bool) *Executor {
	return &Executor{
		Client:      client,
		StepOutputs: make(map[string]actions.Data),
		Vars:        vars,
		DryRun:      dryRun,
	}
}

func (e *Executor) Execute(ctx context.Context, wf *Workflow) error {
	fmt.Println()

	for i := range wf.Steps {
		step := &wf.Steps[i]
		if err := e.executeStep(ctx, step); err != nil {
			switch step.OnError {
			case "continue":
				tui.PrintStepFail(step.ID, err)
				fmt.Printf("  %s Continuing...\n", tui.WarningStyle.Render(tui.WarnIcon))
				continue
			case "retry":
				if err := e.retryStep(ctx, step); err != nil {
					return fmt.Errorf("step %s failed after retries: %w", step.ID, err)
				}
				tui.PrintStepSuccess(step.ID)
			default:
				tui.PrintStepFail(step.ID, err)
				return fmt.Errorf("step %s failed: %w", step.ID, err)
			}
		} else {
			tui.PrintStepSuccess(step.ID)
		}
	}

	fmt.Println()
	fmt.Println(tui.SuccessStyle.Render(fmt.Sprintf("  %s Workflow '%s' completed successfully!", tui.OKIcon, wf.Name)))
	return nil
}

func (e *Executor) executeStep(ctx context.Context, step *Step) error {
	if err := ResolveStepTemplates(step, e.Vars, e.StepOutputs); err != nil {
		return fmt.Errorf("resolving templates: %w", err)
	}

	tui.PrintStepStart(step.ID, step.Action)

	action, ok := actions.Get(step.Action)
	if !ok {
		return fmt.Errorf("unknown action: %s", step.Action)
	}

	if step.Confirm && !e.DryRun {
		if !tui.PrintConfirm(fmt.Sprintf("Step '%s' requires confirmation. Proceed?", step.ID)) {
			fmt.Printf("  %s Step %s skipped by user\n", tui.WarningStyle.Render(tui.WarnIcon), step.ID)
			return nil
		}
	}

	if e.DryRun {
		fmt.Printf("  %s [DRY-RUN] Would execute %s\n", tui.InfoStyle.Render(tui.InfoIcon), step.Action)
		e.StepOutputs[step.ID] = actions.Data{"dry_run": true}
		return nil
	}

	var input actions.Data
	if step.Input != "" {
		var ok bool
		input, ok = e.StepOutputs[step.Input]
		if !ok {
			return fmt.Errorf("input step %s not found or has no output", step.Input)
		}
	}

	output, err := action.Execute(ctx, e.Client, input, step.Params, step.Filters)
	if err != nil {
		return err
	}

	if step.Output != "" {
		e.StepOutputs[step.ID] = output
	}

	return nil
}

func (e *Executor) retryStep(ctx context.Context, step *Step) error {
	maxRetries := step.RetryCount
	if maxRetries <= 0 {
		maxRetries = 3
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("  %s Retry attempt %d/%d for step %s\n",
			tui.WarningStyle.Render(tui.WarnIcon), attempt, maxRetries, step.ID)
		if err := e.executeStep(ctx, step); err != nil {
			if attempt == maxRetries {
				return err
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("step %s failed after %d retries", step.ID, maxRetries)
}