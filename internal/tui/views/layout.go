package views

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	_ "github.com/charmbracelet/lipgloss"
)

var (
	_          tea.Model = Layout{}
	contentBox           = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	box                  = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

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
	var batch tea.BatchMsg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		l.width = msg.Width - box.GetHorizontalFrameSize() - contentBox.GetHorizontalFrameSize()
		l.height = msg.Height - lipgloss.Height(l.Header.View()) - lipgloss.Height(l.Footer.View()) - box.GetVerticalFrameSize() - contentBox.GetVerticalFrameSize() - 1

		nm := tea.WindowSizeMsg{Height: l.height, Width: l.width}
		l.Content, cmd = l.Content.Update(nm)

		batch = append(batch, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			l.ActivePanel = FooterActive
		case "esc":
			l.ActivePanel = ContentActive
		}

		var cmd tea.Cmd
		switch l.ActivePanel {
		case HeaderActive:
			l.Header, cmd = l.Header.Update(msg)
		case ContentActive:
			l.Content, cmd = l.Content.Update(msg)
		case FooterActive:
			l.Footer, cmd = l.Footer.Update(msg)
		default:
			l.Content, cmd = l.Content.Update(msg)
		}

		batch = append(batch, cmd)
	default:
		var c1 tea.Cmd
		var c2 tea.Cmd

		l.Content, c1 = l.Content.Update(msg)
		l.Footer, c2 = l.Footer.Update(msg)

		batch = append(batch, c1, c2)
	}

	return l, tea.Batch(batch...)
}

func (l Layout) View() string {
	ct := contentBox.Width(l.width).Height(l.height)

	slog.Info(
		"Dimensions",
		slog.Int("Height", lipgloss.Height(l.Header.View())),
		slog.Int("ContentHeight", lipgloss.Height(ct.Render(""))),
		slog.Int("FooterView", lipgloss.Height(l.Footer.View())),
	)
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		l.Header.View(),
		ct.Render(l.Content.View()),
		l.Footer.View(),
	)

	return box.Render(view)
}

func (l Layout) SwitchContent(model tea.Model) (tea.Model, tea.Cmd) {
	l.Content = model
	cmds := []tea.Cmd{
		l.Content.Init(),
	}

	nm := tea.WindowSizeMsg{Height: l.height, Width: l.width}

	var cmd tea.Cmd
	l.Content, cmd = l.Content.Update(nm)

	cmds = append(cmds, cmd)

	return l, tea.Batch(cmds...)
}
