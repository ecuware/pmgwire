package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tui "github.com/ecuware/pmgwire/internal/tui"
)

type Table struct {
	Headers   []string
	Rows      [][]string
	Cursor    int
	Width     int
	Selectable bool
}

func NewTable(headers []string) *Table {
	return &Table{
		Headers:    headers,
		Cursor:     0,
		Selectable: true,
	}
}

func (t *Table) AddRow(row ...string) {
	t.Rows = append(t.Rows, row)
}

func (t *Table) MoveUp() {
	if t.Cursor > 0 {
		t.Cursor--
	}
}

func (t *Table) MoveDown() {
	if t.Cursor < len(t.Rows)-1 {
		t.Cursor++
	}
}

func (t *Table) Selected() []string {
	if len(t.Rows) == 0 {
		return nil
	}
	return t.Rows[t.Cursor]
}

func (t *Table) View() string {
	if len(t.Rows) == 0 {
		return tui.DimStyle.Render("  (no data)")
	}

	colWidths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		colWidths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	maxWidth := 0
	for _, w := range colWidths {
		maxWidth += w + 3
	}
	termWidth := maxWidth
	if t.Width > 0 && t.Width < termWidth {
		termWidth = t.Width
	}
	_ = termWidth

	var headerCells []string
	for i, h := range t.Headers {
		if i < len(colWidths) {
			headerCells = append(headerCells, pad(h, colWidths[i]))
		}
	}
	headerLine := tui.TableHeaderStyle.Render(strings.Join(headerCells, " │ "))

	var rows []string
	for i, row := range t.Rows {
		var cells []string
		for j, cell := range row {
			if j < len(colWidths) {
				cells = append(cells, pad(cell, colWidths[j]))
			}
		}

		rowText := strings.Join(cells, " │ ")
		if t.Selectable && i == t.Cursor {
			rowText = tui.SelectedRowStyle.Render(" " + tui.ArrowIcon + " " + rowText)
		} else if t.Selectable {
			rowText = tui.TableRowStyle.Render("   " + rowText)
		} else {
			rowText = tui.TableRowStyle.Render(" " + rowText)
		}
		rows = append(rows, rowText)
	}

	return lipgloss.JoinVertical(lipgloss.Left, headerLine, strings.Join(rows, "\n"))
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func FormatQuarantineTable(mails []map[string]interface{}) string {
	if len(mails) == 0 {
		return tui.DimStyle.Render("  No quarantine mails found.")
	}

	t := NewTable([]string{"ID", "From", "Receiver", "Subject"})
	for _, m := range mails {
		id := fmt.Sprintf("%v", m["ID"])
		from := fmt.Sprintf("%v", m["From"])
		receiver := fmt.Sprintf("%v", m["Receiver"])
		subject := fmt.Sprintf("%v", m["Subject"])
		maxLen := 30
		if len(subject) > maxLen {
			subject = subject[:maxLen-3] + "..."
		}
		t.AddRow(id, from, receiver, subject)
	}
	t.Selectable = false
	return t.View()
}