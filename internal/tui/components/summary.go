package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tui "github.com/ecuware/pmgwire/internal/tui"
)

type Summary struct {
	Title   string
	Success int
	Failed  int
	Total   int
	Details map[string]string
	Items   []SummaryItem
}

type SummaryItem struct {
	Label  string
	Status string
	Error  string
}

func (s Summary) Render() string {
	var b strings.Builder

	title := tui.TitleStyle.Render(fmt.Sprintf(" %s %s ", tui.AppIcon, s.Title))
	b.WriteString(title + "\n\n")

	successBar := ""
	if s.Total > 0 {
		ratio := float64(s.Success) / float64(s.Total)
		barWidth := 30
		filled := int(ratio * float64(barWidth))
		successBar = fmt.Sprintf("  %s%s  %.0f%%",
			tui.SuccessStyle.Render(strings.Repeat("█", filled)),
			tui.DimStyle.Render(strings.Repeat("░", barWidth-filled)),
			ratio*100)
	}

	b.WriteString(fmt.Sprintf("  %s Success: %d  %s Failed: %d  %s Total: %d\n",
		tui.SuccessStyle.Render(tui.OKIcon), s.Success,
		tui.ErrorStyle.Render(tui.FailIcon), s.Failed,
		tui.InfoStyle.Render(tui.ArrowIcon), s.Total))

	if successBar != "" {
		b.WriteString(successBar + "\n")
	}

	if len(s.Items) > 0 {
		b.WriteString("\n")
		for _, item := range s.Items {
			status := ""
			switch item.Status {
			case "success", "ok", "delivered", "added":
				status = tui.SuccessStyle.Render(tui.OKIcon)
			case "failed", "error":
				status = tui.ErrorStyle.Render(tui.FailIcon)
			case "skipped":
				status = tui.WarningStyle.Render(tui.WarnIcon)
			default:
				status = tui.InfoStyle.Render(tui.ArrowIcon)
			}

			line := fmt.Sprintf("  %s %s", status, item.Label)
			if item.Error != "" {
				line += tui.ErrorStyle.Render(fmt.Sprintf(" (%s)", item.Error))
			}
			b.WriteString(line + "\n")
		}
	}

	if len(s.Details) > 0 {
		b.WriteString("\n")
		maxKeyLen := 0
		for k := range s.Details {
			if len(k) > maxKeyLen {
				maxKeyLen = len(k)
			}
		}
		for k, v := range s.Details {
			b.WriteString(fmt.Sprintf("  %-*s  %s\n", maxKeyLen, tui.DimStyle.Render(k+":"), v))
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C56FF")).
		Padding(1, 2).
		Render(b.String())
}

func RenderDeliverSummary(delivered, failed, total int) string {
	s := Summary{
		Title:   "Delivery Summary",
		Success: delivered,
		Failed:  failed,
		Total:   total,
		Details: map[string]string{
			"Delivered": fmt.Sprintf("%d", delivered),
			"Failed":    fmt.Sprintf("%d", failed),
			"Total":     fmt.Sprintf("%d", total),
		},
	}
	return s.Render()
}

func RenderBlacklistSummary(added, failed int, emails, domains []string) string {
	s := Summary{
		Title:   "Blacklist Bulk Summary",
		Success: added,
		Failed:  failed,
		Total:   added + failed,
		Details: map[string]string{
			"Added":   fmt.Sprintf("%d", added),
			"Failed":  fmt.Sprintf("%d", failed),
			"Emails":  fmt.Sprintf("%d", len(emails)),
			"Domains": fmt.Sprintf("%d", len(domains)),
		},
	}
	return s.Render()
}