package config

type Runtime interface {
	Environment() string
	Version() string
	Locales() []string
}

type runtime struct {
	environment string
	version     string
	locales     []string
}

type YamlRuntimeConfig struct {
	Environment string   `yaml:"environment"`
	Version     string   `yaml:"version"`
	Locales     []string `yaml:"locales"`
}

const (
	EnvironmentLocal   = "local"
	EnvironmentDev     = "dev"
	EnvironmentStaging = "staging"
	EnvironmentProd    = "prod"
)

func NewRuntime(env, version string, locales []string) Runtime {
	return &runtime{
		environment: env,
		version:     version,
		locales:     locales,
	}
}

func (d runtime) Environment() string {
	return d.environment
}

func (d runtime) Version() string {
	return d.version
}

func (d runtime) Locales() []string {
	return d.locales
}
