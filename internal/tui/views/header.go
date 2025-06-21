package views

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = Header{}

type Header struct{}

func (h Header) Init() tea.Cmd {
	return nil
}

func (h Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

func (h Header) View() string {
	return "Header"
}
