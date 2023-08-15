package config

import (
	"gopkg.in/yaml.v3"

	"common/containers/deep_map"
)

type ServiceProvider interface {
	Credentials(...string) string
}

type serviceProvider struct {
	providersConfig *deep_map.DeepMap
}

type yamlServiceProviderConfig struct {
	Services yaml.Node `yaml:"services"`
}

func NewServiceProvider(providersCreds yamlServiceProviderConfig) ServiceProvider {
	return &serviceProvider{
		providersConfig: deep_map.NewDeepMap(providersCreds.Services),
	}
}

func (s serviceProvider) Credentials(keys ...string) string {
	currNode := s.providersConfig

	var err error
	for _, key := range keys {
		currNode, err = currNode.Get(key)
		if err != nil {
			panic(err)
		}
	}
	return currNode.GetScalar()
}
