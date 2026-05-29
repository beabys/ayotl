package config

type ConfigMap map[string]interface{}
type ConfigEnvAlias map[string]string

type Config struct {
	ConfigMap    ConfigMap
	EnvConfigMap ConfigMap
	Defaults     ConfigMap
	EnvAlias     ConfigEnvAlias
	immutable    bool
	envLoaded    bool
}

type Params struct {
	Defaults ConfigMap
	EnvAlias ConfigEnvAlias
}
