package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	tui "github.com/ecuware/pmgwire/internal/tui"
)

type FormField struct {
	Name     string
	Label    string
	Default  string
	Required bool
	Value    string
}

type FormModel struct {
	Fields    []FormField
	Cursor    int
	Submitting bool
	Submitted bool
	Quitting  bool
	inputs    []string
	focused   int
}

func NewForm(fields []FormField) FormModel {
	inputs := make([]string, len(fields))
	for i, f := range fields {
		inputs[i] = f.Default
	}
	return FormModel{
		Fields:  fields,
		inputs:  inputs,
		focused: 0,
	}
}

func (m FormModel) Init() tea.Cmd {
	return nil
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.focused > 0 {
				m.focused--
			}
		case "down", "tab":
			if m.focused < len(m.Fields)-1 {
				m.focused++
			}
		case "enter":
			if m.focused == len(m.Fields)-1 {
				m.Submitting = true
				m.Submitted = true
				m.Quitting = true
				return m, tea.Quit
			}
			m.focused++
		case "backspace":
			if len(m.inputs[m.focused]) > 0 {
				m.inputs[m.focused] = m.inputs[m.focused][:len(m.inputs[m.focused])-1]
			}
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		default:
			if len(msg.String()) == 1 {
				m.inputs[m.focused] += msg.String()
			}
		}
	}
	return m, nil
}

func (m FormModel) View() string {
	var b strings.Builder

	b.WriteString(tui.TitleStyle.Render(" PMGWire - Workflow Variables ") + "\n\n")

	for i, field := range m.Fields {
		label := fmt.Sprintf("  %s: ", field.Label)
		if field.Required {
			label = fmt.Sprintf("  %s *: ", field.Label)
		}

		value := m.inputs[i]
		if i == m.focused {
			value = tui.HighlightStyle.Render(value) + "█"
		} else {
			value = tui.DimStyle.Render(value)
		}

		prefix := "  "
		if i == m.focused {
			prefix = tui.BrandNormal.Render("► ")
		}

		b.WriteString(fmt.Sprintf("%s%s%s\n", prefix, label, value))
	}

	b.WriteString("\n  " + tui.DimStyle.Render("[Enter] Submit  [Ctrl+C] Cancel"))

	return b.String()
}

func (m FormModel) GetValues() map[string]string {
	values := make(map[string]string)
	for i, field := range m.Fields {
		values[field.Name] = m.inputs[i]
	}
	return values
}

func RunForm(fields []FormField) (map[string]string, error) {
	m := NewForm(fields)
	p := tea.NewProgram(m)
	model, err := p.Run()
	if err != nil {
		return nil, err
	}
	if fm, ok := model.(FormModel); ok {
		if fm.Quitting && !fm.Submitted {
			return nil, fmt.Errorf("form cancelled")
		}
		return fm.GetValues(), nil
	}
	return nil, fmt.Errorf("unexpected model type")
}