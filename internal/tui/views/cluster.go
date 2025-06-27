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
		AlignHorizontal(lipgloss.Center).
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
	for cc, cluster := range clusters {
		row := table.Row{
			cc,
			strings.Join(cluster.Versions, "\n"),
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
	var batch tea.BatchMsg

	switch msg := msg.(type) {
	case messages.GetClustersMsg:
		log.Println("Received clusters message")
		c.SetRowsFromClusters(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			batch = append(batch, SwitchContentCmd(NewTopic(c.manager)))
		}
	case tea.WindowSizeMsg:
		log.Println("Setting Cluster height to ", msg.Height)
		c.table.SetHeight(msg.Height)
		c.table.SetWidth(msg.Width)
		for i, c2 := range headers {
			headers[i] = table.Column{
				Title: c2.Title,
				Width: (msg.Width - 5) / len(headers),
			}
		}

		c.table.SetColumns(headers)
	}

	var cmd tea.Cmd
	*c.table, cmd = (*c.table).Update(msg)
	batch = append(batch, cmd)

	return c, tea.Batch(batch...)
}

func (c Cluster) View() string {
	return c.table.View()
}
