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
	//create a default ConfigMap
	c.ConfigMap = make(ConfigMap)
	c.EnvConfigMap = make(ConfigMap)
	return c
}

func (c *Config) SetConfigMap(cm ConfigMap) *Config {
	c.ConfigMap = cm
	return c
}

func (c *Config) SetConfigImpl(impl Configuration) *Config {
	c.configImpl = impl
	return c
}

// LoadConfig is a function to load the configurations in ConfigMap
func (c *Config) LoadConfigs(configFiles ...string) (err error) {
	// validate if required files exist to start reading the configs
	for _, configFile := range configFiles {
		if configFile == "" {
			return fmt.Errorf("configuration file should not be empty")
		}

		if _, err := os.Stat(configFile); err != nil {
			fmt.Printf("no able to find file %s \n", configFile)
		}
	}

	// set default values from the implementation
	if c.configImpl != nil {
		for key, val := range c.configImpl.SetDefaults() {
			c.SetDefault(key, val)
		}
	}

	// load the configs from file
	if err := c.getLocalConfigs(configFiles...); err != nil {
		return err
	}

	// merge the env Variables (replace the placeholders) if have values on EnvConfigMap
	if len(c.EnvConfigMap) > 0 {
		c.mergeEnvVariables()
	}

	return nil

}

func (c *Config) getLocalConfigs(configFiles ...string) error {
	for _, s := range configFiles {
		if err := c.ConfigFileMerge(s); err != nil {
			return fmt.Errorf("fail to load configs from file %s: %w", s, err)
		}
	}
	return nil
}

// ConfigFileMerge read configs from file and merge the config into ConfigMap
// if Key exist previosly in ConfigMap, the value will be overridden by the value from the file
func (c *Config) ConfigFileMerge(s string) error {
	if c.ConfigMap == nil {
		c.ConfigMap = make(ConfigMap)
	}
	config, err := ReadFile(s)
	if err != nil {
		return err
	}
	c.ConfigMap = MergeKeys(c.ConfigMap, config)
	return nil
}

// Get return value from given key, and return empty string if key don't exist
// key can be passed in `dot-notation`
func (c *Config) Get(k string) interface{} {
	return GetValue(c.ConfigMap, strings.Split(k, "."))
}

func (c *Config) getEnv(k string) interface{} {
	return GetValue(c.EnvConfigMap, []string{k})
}

// Set add or update value from given key
// key can be passed in `dot-notation`
func (c *Config) Set(k string, v interface{}) {
	SetValue(c.ConfigMap, strings.Split(k, "."), v)
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

// WithEnv Load env variables and add into ConfigMap
func (c *Config) WithEnv(envs ...string) *Config {
	if c.EnvConfigMap == nil {
		c.EnvConfigMap = make(ConfigMap)
	}
	for _, v := range os.Environ() {
		env := strings.SplitN(v, "=", 2)
		if canSave(envs, env[0]) {
			c.EnvConfigMap[env[0]] = env[1]
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

// Unmarshal function convert a ConfigMap type into a struct
// using mapStructure Decoder
func (c *Config) Unmarshal(s any) error {
	if err := mapStructureDecoder(c.ConfigMap, &s); err != nil {
		return fmt.Errorf("unable to unmarshal configurations: %w", err)
	}
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
	mergeENV := MergeEnvVar(c.ConfigMap, c.EnvConfigMap)
	c.ConfigMap = mergeENV
}
