package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cast"
)

// New return  a New Config
func New() *Config {
	c := &Config{}
	//create a default configMap
	c.configMap = make(ConfigMap)
	c.envConfigMap = make(ConfigMap)
	return c
}

func (c *Config) SetConfigMap(cm ConfigMap) *Config {
	c.configMap = cm
	return c
}

func (c *Config) SetConfigImpl(impl Configuration) *Config {
	c.configImpl = impl
	return c
}

// LoadConfig is a function to load the configuration, stored on the config files
// Unmarshalling in the Struct given
func (c *Config) LoadConfigs(configuration interface{}, configFile string) (err error) {
	// validate if required environment variables exist to start reading the configs

	if configFile == "" {
		return fmt.Errorf("configuration file should not be empty")
	}

	if _, err := os.Stat(configFile); err != nil {
		fmt.Printf("no able to find file %s \n", configFile)
	}

	// set default values from the implementation
	if c.configImpl != nil {
		for key, val := range c.configImpl.SetDefaults() {
			c.SetDefault(key, val)
		}
	}

	// load the configs from file
	if err := c.getLocalConfigs(configFile); err != nil {
		return err
	}

	// merge the env Variables (replace the placeholders) if mergEnv is true
	if len(c.envConfigMap) > 0 {
		c.mergeEnvVariables()
	}

	// Unmarshall configs into Config struct
	return c.Unmarshal(&configuration)

}

func (c *Config) getLocalConfigs(s string) (err error) {
	if err := c.ConfigFileMerge(s); err != nil {
		return fmt.Errorf(fmt.Sprintf("Fail to load configs: %s", err.Error()))
	}
	return nil
}

func (c *Config) ConfigFileRead(s string) error {
	config, err := ReadFile(s)
	if err != nil {
		return err
	}
	c.configMap = config
	return nil
}

func (c *Config) ConfigFileMerge(s string) error {
	config, err := ReadFile(s)
	if err != nil {
		return err
	}
	c.configMap = MergeKeys(c.configMap, config)
	return nil
}

func (c *Config) MergeConfigMap(m map[string]interface{}) error {
	c.configMap = MergeKeys(c.configMap, m)
	return nil
}

func (c *Config) Get(k string) interface{} {
	return GetValue(c.configMap, strings.Split(k, "."))
}

func (c *Config) getEnv(k string) interface{} {
	return GetValue(c.envConfigMap, []string{k})
}

func (c *Config) Set(k string, v interface{}) {
	SetValue(c.configMap, strings.Split(k, "."), v)
}

func (c *Config) isSet(k string) bool {
	value := c.Get(k)
	return value != nil
}

func (c *Config) isSetEnv(k string) bool {
	value := c.getEnv(k)
	return value != nil
}

func (c *Config) SetDefault(key string, val interface{}) {
	// if key don't exist we add it
	if !c.isSet(key) {
		c.Set(key, val)
	}
}

// WithEnv Load env variables and add into configmap
func (c *Config) WithEnv(envs ...string) *Config {
	if c.envConfigMap == nil {
		c.envConfigMap = make(ConfigMap)
	}
	for _, v := range os.Environ() {
		env := strings.SplitN(v, "=", 2)
		if canSave(envs, env[0]) {
			c.envConfigMap[env[0]] = env[1]
		}
	}
	return c
}

func canSave(e []string, k string) bool {
	if len(e) < 1 {
		return true
	}
	for _, n := range e {
		if k == n {
			return true
		}
	}
	return false
}

func (c *Config) Unmarshal(s *interface{}) error {
	if err := UnMarshall(c.configMap, s); err != nil {
		return fmt.Errorf(fmt.Sprintf("Unable to unmarshall configurations: %s", err.Error()))
	}
	c.envConfigMap, c.configMap = nil, nil
	return nil
}

// MustString returns the value associated with the key as a string or a default value if empty string.
func (c *Config) MustString(key, must string) string {
	// first search in the env Variables loaded
	switch true {
	case c.isSetEnv(key):
		must = cast.ToString(c.getEnv(key))
	case c.isSet(key):
		must = cast.ToString(c.Get(key))
	}
	return must
}

// MustString returns the value associated with the key as a string or a default value if empty string.
func (c *Config) MustEnvString(key, must string) string {
	if c.isSetEnv(key) {
		must = cast.ToString(c.Get(key))
	}
	return must
}

// MustInt returns the value associated with the key int or a default value if 0.
func (c *Config) MustInt(key string, must int) int {
	val := must
	// first search in the env Variables loaded
	switch true {
	case c.isSetEnv(key):
		val = cast.ToInt(c.getEnv(key))
	case c.isSet(key):
		val = cast.ToInt(c.Get(key))
	}
	return val
}

// MustInt32 returns the value associated with the key as a int32 or a default value if 0.
func (c *Config) MustInt32(key string, must int32) int32 {
	val := must
	// first search in the env Variables loaded
	switch true {
	case c.isSetEnv(key):
		val = cast.ToInt32(c.getEnv(key))
	case c.isSet(key):
		val = cast.ToInt32(c.Get(key))
	}
	return val
}

// MustInt64 returns the value associated with the key as a int64 or a default value if 0.
func (c *Config) MustInt64(key string, must int64) int64 {
	val := must
	switch true {
	case c.isSetEnv(key):
		val = cast.ToInt64(c.getEnv(key))
	case c.isSet(key):
		val = cast.ToInt64(c.Get(key))
	}
	return val
}

// MustBool returns the value associated with the key as a int64 or a default value if 0.
func (c *Config) MustBool(key string, must bool) bool {
	val := must
	switch true {
	case c.isSetEnv(key):
		val = cast.ToBool(c.getEnv(key))
	case c.isSet(key):
		val = cast.ToBool(c.Get(key))
	}
	return val
}

// mergeEnvVariables replace placeholders on config files
func (c *Config) mergeEnvVariables() {
	mergeENV := MergeEnvVar(c.configMap, c.envConfigMap)
	c.configMap = mergeENV
}
