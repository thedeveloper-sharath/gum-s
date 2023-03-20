package pager

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

type search struct {
	active   bool
	input    textinput.Model
	matches  []int
	curMatch int
}

func (s *search) new() {
	input := textinput.New()
	input.Placeholder = "search"
	input.Prompt = "/"
	s.input = input
}

func (s *search) Begin() {
	s.new()
	s.matches = s.matches[0:0]
	s.active = true
	s.input.Focus()
}

func (s *search) Execute(m *model) {
	defer s.Done()
	if s.input.Value() == "" {
		return
	}

	query := regexp.MustCompile(s.input.Value())
	for i, line := range strings.Split(m.content, "\n") {
		if query.Match([]byte(line)) {
			s.matches = append(s.matches, i)
		}
	}

	matches := unique(query.FindAllString(m.content, -1))
	for _, match := range matches {
		replacement := m.matchStyle.Render(match)
		m.content = strings.ReplaceAll(m.content, match, replacement)
	}
}

func unique(strings []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, s := range strings {
		if _, uniq := keys[s]; !uniq {
			keys[s] = true
			list = append(list, s)
		}
	}
	return list
}

func (s *search) Done() {
	s.active = false
	s.curMatch = 0
}

func (s *search) NextMatch(m *model) {
	switch {
	case len(s.matches) <= 0:
		return
	case s.curMatch == len(s.matches)-1:
		(*m).viewport.GotoTop()
		s.curMatch = 0
	case (*m).viewport.AtBottom():
		s.curMatch++
	default:
		for i, match := range s.matches {
			if match > (*m).viewport.YOffset {
				s.curMatch = i
				break
			}
		}
	}

	m.viewport.SetYOffset(m.search.matches[s.curMatch])
}

func (s *search) PrevMatch(m *model) {
	switch {
	case len(s.matches) <= 0:
		return
	case s.curMatch == 0:
		(*m).viewport.GotoBottom()
		s.curMatch = len(s.matches) - 1
	case (*m).viewport.AtBottom():
		s.curMatch--
	default:
		for i := len(s.matches) - 1; i >= 0; i-- {
			if s.matches[i] < (*m).viewport.YOffset {
				s.curMatch = i
				break
			}
		}
	}

	m.viewport.SetYOffset(m.search.matches[s.curMatch])
}
