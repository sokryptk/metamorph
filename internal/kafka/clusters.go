package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kadm"
)

type Cluster struct {
	Name        string
	Versions    []string
	BrokerCount int
}

func GetCluster(ctx context.Context, client *kadm.Client) (info Cluster, err error) {
	brokerVersions, err := client.ApiVersions(ctx)
	if err != nil {
		return Cluster{}, err
	}

	met, err := client.BrokerMetadata(ctx)
	if err != nil {
		return Cluster{}, err
	}

	info.Name = met.Cluster
	info.BrokerCount = len(brokerVersions)
	info.Versions = make([]string, 0, len(brokerVersions))

	brokerVersions.Each(func(bav kadm.BrokerApiVersions) {
		info.Versions = append(info.Versions, bav.VersionGuess())
	})

	return info, nil
}
