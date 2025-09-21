package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Clusters []ClusterConfig `yaml:"clusters"`
}

type ClusterConfig struct {
	Name             string            `yaml:"name"`
	ReadOnly         bool              `yaml:"read_only"`
	BootstrapServers []string          `yaml:"bootstrap_servers"`
	TrustStore       ClusterTrustStore `yaml:"trust_store"`
}

type ClusterTrustStore struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

func ImportClustersFromConfig(path string) (*Config, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
