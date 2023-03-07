package config

type Runtime interface {
	Environment() string
	Version() string
}

type runtime struct {
	environment string
	version     string
}

type YamlRuntimeConfig struct {
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}

const (
	EnvironmentLocal   = "local"
	EnvironmentDev     = "dev"
	EnvironmentStaging = "staging"
	EnvironmentProd    = "prod"
)

func NewRuntime(env, version string) Runtime {
	return &runtime{
		environment: env,
		version:     version,
	}
}

func (d runtime) Environment() string {
	return d.environment
}

func (d runtime) Version() string {
	return d.version
}
