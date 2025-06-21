package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokryptk/metamorph/internal/cluster"
)

var _ tea.Model = Cluster{}

type Cluster struct {
	manager *cluster.Manager
}

func NewCluster(manager *cluster.Manager) Cluster {
	return Cluster{
		manager: manager,
	}
}

func (c Cluster) Init() tea.Cmd {
	return nil
}

func (c Cluster) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//TODO implement me
	panic("implement me")
}

func (c Cluster) View() string {
	return c.View()
}
