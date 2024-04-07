package config

type ConfigMap map[string]interface{}

type Configuration interface {
	SetDefaults() ConfigMap
}

type Config struct {
	configMap    ConfigMap
	envConfigMap ConfigMap
	configImpl   Configuration
}
