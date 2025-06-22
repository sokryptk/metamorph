package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = Footer{}

type Footer struct {
	input textinput.Model
}

func NewFooter() Footer {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Footer{
		input: ti,
	}
}

func (f Footer) Init() tea.Cmd {
	return textinput.Blink
}

func (f Footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	f.input, _ = f.input.Update(msg)

	switch msg.(type) {
	case tea.KeyMsg:
		value := f.input.Value()
		if value == ":topics" {
			return f, SwitchContentCmd(Footer{})
		}
	}
	return f, nil
}

func (f Footer) View() string {
	box := lipgloss.NewStyle()

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			f.input.View(),
		),
	)
}
