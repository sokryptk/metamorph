package kafka

import (
	"context"

	"github.com/sokryptk/metamorph/internal/config"

	"github.com/twmb/franz-go/pkg/kgo"
)

func CreateClientFromConfig(ctx context.Context, config config.ClusterConfig) (*kgo.Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.BootstrapServers...),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}
