package cluster

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokryptk/metamorph/internal/config"
	"github.com/sokryptk/metamorph/internal/kafka"
	"github.com/sokryptk/metamorph/internal/safemap"
	"github.com/sokryptk/metamorph/internal/tui/messages"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type status int

const (
	StatusDisconnected status = iota
	StatusConnecting
	StatusConnected
	StatusError
)

type State struct {
	Config config.ClusterConfig
	ADM    *kadm.Client
	Client *kgo.Client
	Status status
}

type Manager struct {
	ClusterStates *safemap.SafeMap[string, *State]
	ActiveCluster string
}

func NewManager(ctx context.Context, config *config.Config) (*Manager, error) {
	states := make(map[string]*State, len(config.Clusters))

	for _, cc := range config.Clusters {
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

	return &Manager{
		ClusterStates: safemap.NewFromMap(states),
	}, nil
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

func (s *State) GetClusterInfo(ctx context.Context) (kafka.Cluster, error) {
	cluster, err := kafka.GetCluster(ctx, s.ADM)
	if err != nil {
		return kafka.Cluster{}, err
	}

	return cluster, nil
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
		log.Print("Clusters:", messages.GetClustersMsg(clusters))
		return messages.GetClustersMsg(clusters)
	}
}
