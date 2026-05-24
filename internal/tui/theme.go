package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	BrandBold = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C56FF"))

	BrandNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C56FF"))

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C56FF")).
			Padding(0, 2)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECB71"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1C40F"))

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3498DB"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	HighlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C56FF"))

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C56FF")).
			Padding(1, 2)

	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#444444")).
				Padding(0, 1)

	TableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	SelectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#7C56FF")).
				Padding(0, 1)

	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7C56FF"))

	ProgressFill = lipgloss.NewStyle().
			Background(lipgloss.Color("#7C56FF"))

	ProgressEmpty = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333"))

	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C56FF")).
			Padding(0, 3)

	ButtonFocused = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C56FF")).
			Background(lipgloss.Color("#FFFFFF")).
			Padding(0, 3).
			Bold(true)
)

const (
	AppIcon   = "⚡"
	OKIcon    = "✓"
	FailIcon  = "✗"
	WarnIcon  = "⚠"
	InfoIcon  = "ℹ"
	ArrowIcon = "►"
)

func PrintColoredHeader(name, version, description string) {
	fmt.Println(TitleStyle.Render(fmt.Sprintf(" %s PMGWire: %s ", AppIcon, name)))
	if version != "" {
		fmt.Println(SubtitleStyle.Render("  v" + version))
	}
	if description != "" {
		fmt.Println(SubtitleStyle.Render("  " + description))
	}
	fmt.Println()
}

func PrintStepStart(id, action string) {
	fmt.Printf("  %s %s %s\n", BrandNormal.Render(ArrowIcon), BrandBold.Render(id), DimStyle.Render(action))
}

func PrintStepSuccess(id string) {
	fmt.Printf("  %s %s\n", SuccessStyle.Render(OKIcon), SuccessStyle.Render(id))
}

func PrintStepFail(id string, err error) {
	fmt.Printf("  %s %s %s\n", ErrorStyle.Render(FailIcon), ErrorStyle.Render(id), ErrorStyle.Render(err.Error()))
}

func PrintConfirm(message string) bool {
	fmt.Printf("\n  %s %s (y/N): ", WarningStyle.Render(WarnIcon), message)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}