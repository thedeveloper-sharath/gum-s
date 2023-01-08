// Package pager provides a pager (similar to less) for the terminal.
//
// $ cat file.txt | gum page
package pager

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/gum/timeout"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type model struct {
	content         string
	viewport        viewport.Model
	helpStyle       lipgloss.Style
	showLineNumbers bool
	lineNumberStyle lipgloss.Style
	softWrap        bool
	timeout         time.Duration
	hasTimeout      bool
}

func (m model) Init() tea.Cmd {
	return timeout.Init(m.timeout, nil)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)

	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - lipgloss.Height(m.helpStyle.Render("?")) - 1
		m.viewport.Width = msg.Width
		textStyle := lipgloss.NewStyle().Width(m.viewport.Width)
		var text strings.Builder

		// Determine max width of a line
		maxLineWidth := m.viewport.Width
		if m.softWrap {
			vpStyle := m.viewport.Style
			maxLineWidth -= vpStyle.GetHorizontalBorderSize() + vpStyle.GetHorizontalMargins() + vpStyle.GetHorizontalPadding()
			if m.showLineNumbers {
				maxLineWidth -= len("     │ ")
			}
		}

		for i, line := range strings.Split(m.content, "\n") {
			line = strings.ReplaceAll(line, "\t", "    ")
			if m.showLineNumbers {
				text.WriteString(m.lineNumberStyle.Render(fmt.Sprintf("%4d │ ", i+1)))
			}
			for m.softWrap && len(line) > maxLineWidth {
				truncatedLine := runewidth.Truncate(line, maxLineWidth, "")
				text.WriteString(textStyle.Render(truncatedLine))
				text.WriteString("\n")
				if m.showLineNumbers {
					text.WriteString(m.lineNumberStyle.Render("     │ "))
				}
				line = strings.Replace(line, truncatedLine, "", 1)
			}
			text.WriteString(textStyle.Render(runewidth.Truncate(line, maxLineWidth, "")))
			text.WriteString("\n")
		}

		diffHeight := m.viewport.Height - lipgloss.Height(text.String())
		if diffHeight > 0 && m.showLineNumbers {
			remainingLines := "   ~ │ " + strings.Repeat("\n   ~ │ ", diffHeight-1)
			text.WriteString(m.lineNumberStyle.Render(remainingLines))
		}
		m.viewport.SetContent(text.String())
	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			m.viewport.GotoTop()
		case "G":
			m.viewport.GotoBottom()
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var timeoutStr string
	if m.hasTimeout {
		timeoutStr = timeout.TimeoutStr(m.timeout) + " "
	}
	return m.viewport.View() + m.helpStyle.Render("\n"+timeoutStr+"↑/↓: Navigate • q: Quit")
}
