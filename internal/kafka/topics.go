package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kadm"
)

type Topic struct {
	Name       string
	IsInternal isinternal
	Partitions kadm.PartitionDetails
	ID         string
}

type isinternal bool

func (i isinternal) String() string {
	if i {
		return "true"
	}
	return "false"
}

func GetTopics(ctx context.Context, client *kadm.Client) (results []Topic, err error) {
	var ts []Topic
	topics, err := client.ListTopicsWithInternal(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range topics {
		ts = append(ts, Topic{
			Name:       t.Topic,
			IsInternal: isinternal(t.IsInternal),
			Partitions: t.Partitions,
			ID:         t.ID.String(),
		})
	}

	return ts, nil
}
