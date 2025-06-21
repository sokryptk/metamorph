package cluster

import (
	"context"

	"github.com/sokryptk/metamorph/internal/config"
	"github.com/sokryptk/metamorph/internal/kafka"
	"github.com/sokryptk/metamorph/internal/safemap"
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

func (m *Manager) GetCurrentCluster() (*State, bool) {
	return m.ClusterStates.Get(m.ActiveCluster)
}
