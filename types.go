package config

type ConfigMap map[string]interface{}

type Configuration interface {
	SetDefaults() ConfigMap
}

type Config struct {
	ConfigMap    ConfigMap
	EnvConfigMap ConfigMap
	configImpl   Configuration
}
