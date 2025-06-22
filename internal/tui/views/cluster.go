package views

import (
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sokryptk/metamorph/internal/cluster"
	"github.com/sokryptk/metamorph/internal/kafka"
	"github.com/sokryptk/metamorph/internal/tui/messages"
)

var headers = []table.Column{
	{Title: "Name"},
	{Title: "Version"},
	{Title: "Broker Count"},
}

var _ tea.Model = Cluster{}

type Cluster struct {
	width, height int
	manager       *cluster.Manager
	table         *table.Model
}

func NewCluster(manager *cluster.Manager) Cluster {
	ti := table.New(
		table.WithColumns(headers),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	ti.SetStyles(s)

	return Cluster{
		table:   &ti,
		manager: manager,
	}
}

func (c Cluster) SetRowsFromClusters(clusters map[string]kafka.Cluster) {
	rows := make([]table.Row, 0)
	for _, cluster := range clusters {
		row := table.Row{
			cluster.Name,
			strings.Join(cluster.Versions, ", "),
			strconv.Itoa(cluster.BrokerCount),
		}

		rows = append(rows, row)
	}

	c.table.SetRows(rows)
}

func (c Cluster) Init() tea.Cmd {
	return c.manager.GetClustersCmd()
}

func (c Cluster) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.GetClustersMsg:
		log.Println("Received clusters message")
		c.SetRowsFromClusters(msg)
	case tea.WindowSizeMsg:
		c.table.SetHeight(msg.Height)
		c.table.SetWidth(msg.Width)
		for i, c2 := range headers {
			headers[i] = table.Column{
				Title: c2.Title,
				Width: msg.Width / len(headers),
			}
		}

		c.table.SetColumns(headers)
	}

	*c.table, cmd = (*c.table).Update(msg)
	return c, cmd
}

func (c Cluster) View() string {
	return c.table.View()
}
