package views

import (
	"log"

	"github.com/sokryptk/metamorph/internal/cluster"
	"github.com/sokryptk/metamorph/internal/tui/messages"

	tea "github.com/charmbracelet/bubbletea"
)

// This is the main view
// We reroute to every other view from here basically

var _ tea.Model = MetamorphicView{}

type MetamorphicView struct {
	Manager  *cluster.Manager
	Embedded tea.Model
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

	log.Println("Received message", msg)
	switch msg := msg.(type) {
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
