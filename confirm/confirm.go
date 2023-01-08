// Package confirm provides an interface to ask a user to confirm an action.
// The user is provided with an interface to choose an affirmative or negative
// answer, which is then reflected in the exit code for use in scripting.
//
// If the user selects the affirmative answer, the program exits with 0. If the
// user selects the negative answer, the program exits with 1.
//
// I.e. confirm if the user wants to delete a file
//
// $ gum confirm "Are you sure?" && rm file.txt
package confirm

import (
	"github.com/charmbracelet/gum/timeout"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type defaultVal struct {
	Value bool
}

type model struct {
	prompt       string
	affirmative  string
	negative     string
	quitting     bool
	aborted      bool
	hasTimeout   bool
	timeout      time.Duration
	confirmation bool
	defvalue     bool
	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return timeout.Init(m.timeout, defaultVal{Value: m.confirmation})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.confirmation = false
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case "q", "n", "N":
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		case "left", "h", "ctrl+p", "tab",
			"right", "l", "ctrl+n", "shift+tab":
			if m.negative == "" {
				break
			}
			m.confirmation = !m.confirmation
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case "y", "Y":
			m.quitting = true
			m.confirmation = true
			return m, tea.Quit
		}
	case timeout.TimeoutMsg:

		if msg.TimeoutValue <= 0 {
			m.quitting = true
			if v, ok := msg.Data.(defaultVal); ok {
				m.confirmation = v.Value
			} else {
				m.confirmation = false
			}

			return m, tea.Quit
		}

		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var aff, neg, timeoutStrYes, timeoutStrNo string
	timeoutStrNo = ""
	timeoutStrYes = ""
	if m.hasTimeout {
		if m.defvalue {
			timeoutStrYes = timeout.TimeoutStr(m.timeout)
		} else {
			timeoutStrNo = timeout.TimeoutStr(m.timeout)
		}
	}

	if m.confirmation {
		aff = m.selectedStyle.Render(m.affirmative + timeoutStrYes)
		neg = m.unselectedStyle.Render(m.negative + timeoutStrNo)
	} else {
		aff = m.unselectedStyle.Render(m.affirmative + timeoutStrYes)
		neg = m.selectedStyle.Render(m.negative + timeoutStrNo)
	}

	// If the option is intentionally empty, do not show it.
	if m.negative == "" {
		neg = ""
	}

	return lipgloss.JoinVertical(lipgloss.Center, m.promptStyle.Render(m.prompt), lipgloss.JoinHorizontal(lipgloss.Left, aff, neg))
}
