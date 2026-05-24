package components

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	tui "github.com/ecuware/pmgwire/internal/tui"
)

type ConfirmModel struct {
	Message  string
	Yes      bool
	Confirmed bool
	Quitting  bool
}

func NewConfirm(message string) ConfirmModel {
	return ConfirmModel{
		Message: message,
		Yes:     true,
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "tab":
			m.Yes = !m.Yes
		case "enter":
			m.Confirmed = m.Yes
			m.Quitting = true
			return m, tea.Quit
		case "q", "ctrl+c":
			m.Confirmed = false
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	yesBtn := tui.ButtonStyle.Render(" Yes ")
	noBtn := tui.ButtonStyle.Render(" No ")
	if m.Yes {
		yesBtn = tui.ButtonFocused.Render(" Yes ")
	} else {
		noBtn = tui.ButtonFocused.Render(" No ")
	}

	content := fmt.Sprintf("\n  %s\n\n  %s  %s\n", m.Message, yesBtn, noBtn)
	return tui.BorderStyle.Render(content)
}

func RunConfirm(message string) (bool, error) {
	p := tea.NewProgram(NewConfirm(message))
	m, err := p.Run()
	if err != nil {
		return false, err
	}
	if model, ok := m.(ConfirmModel); ok {
		return model.Confirmed, nil
	}
	return false, nil
}