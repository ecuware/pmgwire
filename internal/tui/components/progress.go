package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	tui "github.com/ecuware/pmgwire/internal/tui"
)

type ProgressModel struct {
	Total     int
	Current   int
	Label     string
	Done      bool
	Quitting  bool
}

func NewProgress(total int, label string) ProgressModel {
	return ProgressModel{
		Total:   total,
		Current: 0,
		Label:   label,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		m.Current = msg.Current
		m.Done = msg.Done
		if m.Done {
			m.Quitting = true
			return m, tea.Quit
		}
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ProgressModel) View() string {
	width := 40
	percent := 0.0
	if m.Total > 0 {
		percent = float64(m.Current) / float64(m.Total)
	}
	filled := int(percent * float64(width))
	empty := width - filled

	bar := tui.ProgressFill.Render(strings.Repeat("█", filled)) +
		tui.ProgressEmpty.Render(strings.Repeat("░", empty))

	status := fmt.Sprintf("%s [%d/%d]", m.Label, m.Current, m.Total)
	if m.Done {
		status = tui.SuccessStyle.Render(fmt.Sprintf("  %s %s [%d/%d] Complete!", tui.OKIcon, m.Label, m.Current, m.Total))
	}

	return fmt.Sprintf("\n  %s\n  %s %.1f%%\n", status, bar, percent*100)
}

type ProgressMsg struct {
	Current int
	Done    bool
}

func RunProgress(total int, label string, items []func() error) (int, int) {
	success := 0
	failed := 0

	for i, fn := range items {
		if err := fn(); err != nil {
			failed++
			fmt.Printf("  %s [%d/%d] %v\n", tui.FailIcon, i+1, total, err)
		} else {
			success++
			fmt.Printf("  %s [%d/%d]\n", tui.OKIcon, i+1, total)
		}
	}

	return success, failed
}