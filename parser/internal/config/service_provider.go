package config

import "github.com/pkg/errors"

type ServiceProvider interface {
	Credentials(string) map[string]string
	CrawlersCount() int
}

type serviceProvider struct {
	providersConfig yamlServiceProviderConfig
}

type yamlServiceProviderConfig struct {
	Services map[string]map[string]string `yaml:"services"`
}

func NewServiceProvider(providersCreds yamlServiceProviderConfig) ServiceProvider {
	return &serviceProvider{
		providersConfig: providersCreds,
	}
}

func (s serviceProvider) Credentials(providerName string) map[string]string {
	if _, ok := s.providersConfig.Services[providerName]; !ok {
		panic(errors.Errorf("unknown provider %s", providerName))
	}
	return s.providersConfig.Services[providerName]
}

func (s serviceProvider) CrawlersCount() int {
	return len(s.providersConfig.Services)
}
