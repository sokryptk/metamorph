package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokryptk/metamorph/internal/kafka"
)

type InitMsg struct{}

type GetClustersMsg map[string]kafka.Cluster

type SwitchContentMsg struct {
	Model tea.Model
}
