package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	_ "github.com/charmbracelet/lipgloss"
)

var _ tea.Model = Layout{}

type ActivePanel int

const (
	HeaderActive ActivePanel = iota
	ContentActive
	FooterActive
)

type Layout struct {
	width, height int
	Header        tea.Model
	Content       tea.Model
	Footer        tea.Model
	ActivePanel   ActivePanel
}

func NewLayoutWithContent(view tea.Model) Layout {
	return Layout{
		Header:      Header{},
		Content:     view,
		Footer:      NewFooter(),
		ActivePanel: ContentActive,
	}
}

func (l Layout) Init() tea.Cmd {
	return l.Content.Init()
}

func (l Layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height

		l.height = l.height - lipgloss.Height(l.Header.View()) - lipgloss.Height(l.Footer.View()) - 2

		nm := tea.WindowSizeMsg{Height: l.height, Width: l.width}
		l.Content, _ = l.Content.Update(nm)
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			l.ActivePanel = FooterActive
		case "esc":
			l.ActivePanel = ContentActive
		}

		switch l.ActivePanel {
		case HeaderActive:
			l.Header, _ = l.Header.Update(msg)
		case ContentActive:
			l.Content, _ = l.Content.Update(msg)
		case FooterActive:
			l.Footer, cmd = l.Footer.Update(msg)
		default:
			l.Content, _ = l.Content.Update(msg)
		}
	default:
		l.Content, _ = l.Content.Update(msg)
		l.Footer, cmd = l.Footer.Update(msg)
	}

	return l, cmd
}

func (l Layout) View() string {
	box := lipgloss.NewStyle().Width(l.width).Height(l.height)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		l.Header.View(),
		l.Content.View(),
		l.Footer.View(),
	)

	return box.Render(view)
}

func (l Layout) SwitchContent(model tea.Model) (tea.Model, tea.Cmd) {
	l.Content = model
	cmds := []tea.Cmd{
		l.Content.Init(),
	}

	contentHeight := l.height - lipgloss.Height(l.Header.View()) - lipgloss.Height(l.Footer.View())
	msg := tea.WindowSizeMsg{
		Height: contentHeight,
		Width:  l.width,
	}

	var cmd tea.Cmd
	l.Content, cmd = l.Content.Update(msg)

	cmds = append(cmds, cmd)

	return l, tea.Batch(cmds...)
}
