package views

import (
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sokryptk/metamorph/internal/cluster"
	"github.com/sokryptk/metamorph/internal/kafka"
	"github.com/sokryptk/metamorph/internal/tui/messages"
)

var topicsHeaders = []table.Column{
	{Title: "ID"},
	{Title: "Name"},
	{Title: "Partition Count"},
}

var _ tea.Model = Cluster{}

type Topic struct {
	width, height int
	manager       *cluster.Manager
	table         *table.Model
}

func NewTopic(manager *cluster.Manager) Topic {
	ti := table.New(
		table.WithColumns(topicsHeaders),
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

	return Topic{
		table:   &ti,
		manager: manager,
	}
}

func (t Topic) SetRows(topics []kafka.Topic) {
	rows := make([]table.Row, 0)
	for _, topics := range topics {
		row := table.Row{
			topics.ID,
			topics.Name,
			string(topics.Partitions.Numbers()),
		}

		rows = append(rows, row)
	}

	t.table.SetRows(rows)
}

func (t Topic) Init() tea.Cmd {
	return t.manager.GetTopicsCmd()
}

func (t Topic) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.GetTopicsMsg:
		log.Println("Received clusters message")
		t.SetRows(msg)
	case tea.WindowSizeMsg:
		t.table.SetHeight(msg.Height)
		t.table.SetWidth(msg.Width)
		for i, c2 := range topicsHeaders {
			topicsHeaders[i] = table.Column{
				Title: c2.Title,
				Width: (msg.Width - 5) / len(topicsHeaders),
			}
		}

		t.table.SetColumns(topicsHeaders)
	}

	*t.table, cmd = (*t.table).Update(msg)
	return t, cmd
}

func (t Topic) View() string {
	return t.table.View()
}
