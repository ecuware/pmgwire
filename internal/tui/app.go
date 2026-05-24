package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/ecuware/pmgwire/internal/actions"
	"github.com/ecuware/pmgwire/internal/pmg"
)

type WorkflowStepInfo struct {
	ID     string
	Action string
}

type AppModel struct {
	WorkflowName    string
	WorkflowDesc    string
	Steps           []WorkflowStepInfo
	Client          *pmg.Client
	Vars            map[string]interface{}
	DryRun          bool
	StepActions     []EngineStep
	StepOutputs     map[string]actions.Data

	Phase       string
	CurrentStep int
	StepStates  []string
	StepErrors  []string

	Width    int
	Height   int
	Quitting bool
	Output   string
	Err      error
}

type EngineStep struct {
	ID       string
	Action   string
	Params   map[string]interface{}
	Filters  map[string]string
	Input    string
	Output   string
	Confirm  bool
	OnError  string
	RetryCnt int
}

const (
	PhaseVars      = "vars"
	PhaseConfirm   = "confirm"
	PhaseExecuting = "executing"
	PhaseDone      = "done"
	PhaseError     = "error"
)

func NewAppModel(wfName, wfDesc string, steps []EngineStep, client *pmg.Client, vars map[string]interface{}, dryRun bool) AppModel {
	stepInfos := make([]WorkflowStepInfo, len(steps))
	states := make([]string, len(steps))
	errors := make([]string, len(steps))
	for i, s := range steps {
		stepInfos[i] = WorkflowStepInfo{ID: s.ID, Action: s.Action}
		states[i] = "pending"
		errors[i] = ""
	}

	return AppModel{
		WorkflowName: wfName,
		WorkflowDesc: wfDesc,
		Steps:        stepInfos,
		StepActions:  steps,
		Client:       client,
		Vars:         vars,
		DryRun:       dryRun,
		StepOutputs:  make(map[string]actions.Data),
		Phase:        PhaseExecuting,
		StepStates:   states,
		StepErrors:   errors,
	}
}

func (m AppModel) Init() tea.Cmd {
	if len(m.StepActions) > 0 {
		return executeStepCmd(0, m)
	}
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			if m.Phase == PhaseDone || m.Phase == PhaseError {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	case stepResultMsg:
		if msg.Err != nil {
			m.StepStates[msg.Index] = "failed"
			m.StepErrors[msg.Index] = msg.Err.Error()

			step := m.StepActions[msg.Index]
			switch step.OnError {
			case "continue":
				m.StepStates[msg.Index] = "failed"
				if msg.Index+1 < len(m.StepActions) {
					return m, executeStepCmd(msg.Index+1, m)
				}
				m.Phase = PhaseDone
			case "retry":
				m.StepStates[msg.Index] = "failed"
				if msg.Index+1 < len(m.StepActions) {
					return m, executeStepCmd(msg.Index+1, m)
				}
				m.Phase = PhaseDone
			default:
				m.Err = msg.Err
				m.Phase = PhaseError
				return m, nil
			}
		}
		m.StepStates[msg.Index] = "success"
		if msg.Output != nil {
			m.StepOutputs[msg.StepID] = msg.Output
		}
		m.CurrentStep = msg.Index + 1
		if msg.Index+1 < len(m.StepActions) {
			return m, executeStepCmd(msg.Index+1, m)
		}
		m.Phase = PhaseDone
	}
	return m, nil
}

func (m AppModel) View() string {
	var b strings.Builder

	header := TitleStyle.Render(fmt.Sprintf(" %s PMGWire: %s ", AppIcon, m.WorkflowName))
	if m.WorkflowDesc != "" {
		header += "\n" + SubtitleStyle.Render("  "+m.WorkflowDesc)
	}
	b.WriteString(header + "\n\n")

	if m.DryRun {
		b.WriteString(WarningStyle.Render("  ⚠ DRY RUN MODE") + "\n\n")
	}

	for i, step := range m.Steps {
		icon := "○"
		style := DimStyle
		switch m.StepStates[i] {
		case "success":
			icon = OKIcon
			style = SuccessStyle
		case "failed":
			icon = FailIcon
			style = ErrorStyle
		case "running":
			icon = "◎"
			style = BrandNormal
		}
		b.WriteString(fmt.Sprintf("  %s %s %s\n", style.Render(icon), style.Render(step.ID), DimStyle.Render(step.Action)))
		if m.StepErrors[i] != "" {
			b.WriteString(ErrorStyle.Render(fmt.Sprintf("    └─ %s\n", m.StepErrors[i])))
		}
	}

	if m.Phase == PhaseDone {
		b.WriteString("\n" + SuccessStyle.Render(fmt.Sprintf("  %s Workflow completed successfully!", OKIcon)) + "\n")
		b.WriteString(DimStyle.Render("\n  Press q to exit.\n"))
	}

	if m.Phase == PhaseError {
		b.WriteString("\n" + ErrorStyle.Render(fmt.Sprintf("  %s Workflow failed!", FailIcon)) + "\n")
		b.WriteString(DimStyle.Render("\n  Press q to exit.\n"))
	}

	return b.String()
}

type stepResultMsg struct {
	Index  int
	StepID string
	Err    error
	Output  actions.Data
}

func executeStepCmd(index int, m AppModel) tea.Cmd {
	return func() tea.Msg {
		if index >= len(m.StepActions) {
			return stepResultMsg{Index: index, Err: fmt.Errorf("step index out of range")}
		}

		step := m.StepActions[index]
		action, ok := actions.Get(step.Action)
		if !ok {
			return stepResultMsg{
				Index:  index,
				StepID: step.ID,
				Err:    fmt.Errorf("unknown action: %s", step.Action),
			}
		}

		var input actions.Data
		if step.Input != "" {
			var found bool
			input, found = m.StepOutputs[step.Input]
			if !found {
				return stepResultMsg{
					Index:  index,
					StepID: step.ID,
					Err:    fmt.Errorf("input step %s not found", step.Input),
				}
			}
		}

		output, err := action.Execute(context.Background(), m.Client, input, step.Params, step.Filters)
		if err != nil {
			return stepResultMsg{
				Index:  index,
				StepID: step.ID,
				Err:    err,
			}
		}

		return stepResultMsg{
			Index:  index,
			StepID: step.ID,
			Output:  output,
		}
	}
}

func RunApp(wfName, wfDesc string, steps []EngineStep, client *pmg.Client, vars map[string]interface{}, dryRun bool) error {
	m := NewAppModel(wfName, wfDesc, steps, client, vars, dryRun)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}