package views

import (
	"github.com/sokryptk/metamorph/internal/cluster"
	"github.com/sokryptk/metamorph/internal/tui/messages"

	tea "github.com/charmbracelet/bubbletea"
)

// This is the main view
// We reroute to every other view from here basically

var _ tea.Model = MetamorphicView{}

type MetamorphicView struct {
	width, height int
	Manager       *cluster.Manager
	Embedded      tea.Model
}

func NewMetamorph(manager *cluster.Manager) MetamorphicView {
	return MetamorphicView{
		Manager:  manager,
		Embedded: NewLayoutWithContent(Blank{}),
	}
}

func SwitchContentCmd(model tea.Model) tea.Cmd {
	return func() tea.Msg {
		return messages.SwitchContentMsg{Model: model}
	}
}

func (m MetamorphicView) Init() tea.Cmd {
	return tea.Batch(
		SwitchContentCmd(NewCluster(m.Manager)),
		m.Embedded.Init(),
	)
}

func (m MetamorphicView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Just keep a track of the current window size etc
		m.width = msg.Width
		m.height = msg.Height

		m.Embedded, cmd = m.Embedded.Update(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			m.Embedded, cmd = m.Embedded.Update(msg)
		}
	case messages.SwitchContentMsg:
		m.Embedded, cmd = m.Embedded.(Layout).SwitchContent(msg.Model)
	default:
		m.Embedded, cmd = m.Embedded.Update(msg)
	}

	return m, cmd
}

func (m MetamorphicView) View() string {
	return m.Embedded.View()
}
