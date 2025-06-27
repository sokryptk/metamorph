package cluster

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokryptk/metamorph/internal/config"
	"github.com/sokryptk/metamorph/internal/kafka"
	"github.com/sokryptk/metamorph/internal/safemap"
	"github.com/sokryptk/metamorph/internal/tui/messages"
	"github.com/twmb/franz-go/pkg/kadm"
)

type Manager struct {
	ClusterStates *safemap.SafeMap[string, *State]
	ActiveCluster string
	ActiveState   *State
}

func NewManager(ctx context.Context, config *config.Config) (*Manager, error) {
	states := make(map[string]*State, len(config.Clusters))
	manager := &Manager{}

	for _, cc := range config.Clusters {
		if manager.ActiveCluster == "" {
			manager.ActiveCluster = cc.Name
			manager.ActiveState = &State{
				Config: cc,
				Client: nil,
				ADM:    nil,
				Status: StatusDisconnected,
			}
		}

		client, err := kafka.CreateClientFromConfig(ctx, cc)
		if err != nil {
			return nil, err
		}

		states[cc.Name] = &State{
			Config: cc,
			Client: client,
			ADM:    kadm.NewClient(client),
			Status: StatusConnecting,
		}

		err = client.Ping(ctx)
		if err != nil {
			states[cc.Name].Status = StatusError
			continue
		}

		states[cc.Name].Status = StatusConnected
	}

	manager.ClusterStates = safemap.NewFromMap(states)

	return manager, nil
}

func (m *Manager) GetClusterState(name string) (*State, bool) {
	return m.ClusterStates.Get(name)
}

func (m *Manager) GetCurrentClusterState() (*State, bool) {
	return m.ClusterStates.Get(m.ActiveCluster)
}

func (m *Manager) Range(f func(string, *State) bool) {
	m.ClusterStates.Range(f)
}

func (m *Manager) MustGetCurrentClusterState() *State {
	state, ok := m.GetCurrentClusterState()
	if !ok {
		panic("no active cluster")
	}
	return state
}

func (m *Manager) GetClusters() map[string]kafka.Cluster {
	info := make(map[string]kafka.Cluster, m.ClusterStates.Len())
	m.Range(func(name string, state *State) bool {
		cluster, err := state.GetClusterInfo(context.Background())
		if err != nil {
			return true
		}

		info[name] = cluster
		return true
	})

	return info
}

func (m *Manager) GetClustersCmd() func() tea.Msg {
	return func() tea.Msg {
		clusters := m.GetClusters()
		return messages.GetClustersMsg(clusters)
	}
}

func (m *Manager) GetTopicsCmd() func() tea.Msg {
	return func() tea.Msg {
		topics, err := m.MustGetCurrentClusterState().GetTopics(context.Background())
		if err != nil {
			return messages.GetTopicsMsg(nil)
		}
		return messages.GetTopicsMsg(topics)
	}
}
