package views

import (
	"github.com/sokryptk/metamorph/internal/cluster"

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
		Embedded: NewLayout(),
	}
}

func (m MetamorphicView) Init() tea.Cmd {
	return nil
}

func (m MetamorphicView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Embedded, _ = m.Embedded.Update(msg)
	return m, nil
}

func (m MetamorphicView) View() string {
	return m.Embedded.View()
}
