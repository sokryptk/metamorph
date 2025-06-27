package cluster

import (
	"context"

	"github.com/sokryptk/metamorph/internal/config"
	"github.com/sokryptk/metamorph/internal/kafka"
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

func (s *State) GetClusterInfo(ctx context.Context) (kafka.Cluster, error) {
	cluster, err := kafka.GetCluster(ctx, s.ADM)
	if err != nil {
		return kafka.Cluster{}, err
	}

	return cluster, nil
}

func (s *State) GetTopics(ctx context.Context) ([]kafka.Topic, error) {
	topics, err := kafka.GetTopics(ctx, s.ADM)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
