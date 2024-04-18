package config

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockConfig struct {
	App     MockApplicationConfig `mapstructure:"application"`
	Logger  MockLoggerConfig      `mapstructure:"logger"`
	Logger2 MockLoggerConfig      `mapstructure:"logger2"`
}
type MockApplicationConfig struct {
	Port int `mapstructure:"port"`
}
type MockLoggerConfig struct {
	LogOutput   string `mapstructure:"log_output_to"`
	ErrorOutput string `mapstructure:"log_errors_to"`
	Level       string `mapstructure:"log_level"`
}

func TestConfig(t *testing.T) {

	path := "./testConfig/"
	defer os.RemoveAll(path)
	t.Run("test new should return a new config", func(t *testing.T) {
		mock := &MockConfig{}
		configMap := make(ConfigMap)
		envConfigMap := make(ConfigMap)
		config := New().SetConfigImpl(mock).SetConfigMap(configMap).WithEnv()
		var want = &Config{configMap, envConfigMap, mock}
		want.WithEnv()
		areEqual := assert.ObjectsAreEqual(config, want)
		assert.True(t, areEqual)
	})

	t.Run("test Error Loading configs path without CONFIG_FILE", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		mock := &MockConfig{}
		config := New()
		assert.ErrorContains(t, config.LoadConfigs(mock, ""), "configuration file")
	})

	t.Run("test Error Loading configs file", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock)
		assert.ErrorContains(t, config.LoadConfigs(mock, "./../env.configuration.json"), "Fail to load configs")
	})

	t.Run("Test Loading configs", func(t *testing.T) {
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		data := `{"server": {"enabled": false},"second": {"config": {"enabled": false}}}`
		testPath, err := createTestConfigFile(path, "/config.json", data)
		assert.NoError(t, err)
		assert.NoError(t, config.LoadConfigs(mock, testPath+"/config.json"))
		assert.NoError(t, config.Unmarshal(mock))
	})

	t.Run("Test Loading configs with config placeholders", func(t *testing.T) {
		mock := &MockConfig{}
		data := `{"server": {"enabled": "${IS_CONFIG_FOR_TEST_ENABLED}"},"second": {"config": {"enabled": false}}}`
		testPath, err := createTestConfigFile(path, "/config.json", data)
		assert.NoError(t, err)
		os.Setenv("IS_CONFIG_FOR_TEST_ENABLED", "true")
		// c := New()
		config := New().SetConfigImpl(mock).WithEnv()
		assert.NoError(t, config.LoadConfigs(mock, testPath+"/config.json"))
		assert.NoError(t, config.Unmarshal(mock))
		os.Unsetenv("IS_CONFIG_FOR_TEST_ENABLED")
	})

	t.Run("Test Loading configs with config placeholders and no Value", func(t *testing.T) {
		mock := &MockConfig{}
		data := `{"server": {"enabled": "${IS_CONFIG_FOR_TEST_ENABLED_SECOND}"},"second": {"config": {"enabled": false}}}`
		testPath, err := createTestConfigFile(path, "/config.json", data)
		assert.NoError(t, err)
		// c := &Config{}
		config := New().SetConfigImpl(mock).WithEnv()
		assert.NoError(t, config.LoadConfigs(mock, testPath+"/config.json"))
		assert.NoError(t, config.Unmarshal(mock))
	})
}

func TestMustFunctions(t *testing.T) {
	t.Run("test MustBool", func(t *testing.T) {
		mock := &MockConfig{}
		os.Setenv("ANY", "false")
		config := New().SetConfigImpl(mock).WithEnv()
		required := true
		val := config.MustBool("ANY", required)
		assert.False(t, val)
		os.Unsetenv("ANY")
	})
	t.Run("test MustBool exist", func(t *testing.T) {
		os.Setenv("ANY", "true")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		required := true
		val := config.MustBool("ANY", required)
		assert.True(t, val)
		os.Unsetenv("ANY")
	})
	t.Run("test MustString", func(t *testing.T) {
		os.Setenv("ANY", "required")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		required := "required"
		val := config.MustString("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt", func(t *testing.T) {
		os.Setenv("ANY", "0")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		required := 0
		val := config.MustInt("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt32", func(t *testing.T) {
		os.Setenv("ANY", "0")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		integer := 0
		required := int32(integer)
		val := config.MustInt32("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt64", func(t *testing.T) {
		os.Setenv("ANY", "0")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		integer := 0
		required := int64(integer)
		val := config.MustInt64("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustString with variable", func(t *testing.T) {
		os.Setenv("ANY", "required")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		required := "required"
		val := config.MustString("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt with variable", func(t *testing.T) {
		os.Setenv("ANY", "10")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		required := 10
		val := config.MustInt("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt32 with variable", func(t *testing.T) {
		os.Setenv("ANY", "10")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		integer := 10
		required := int32(integer)
		val := config.MustInt32("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt64 with variable", func(t *testing.T) {
		os.Setenv("ANY", "10")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).WithEnv()
		integer := 10
		required := int64(integer)
		val := config.MustInt64("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})
}

func (mc *MockConfig) SetDefaults() ConfigMap {
	defaults := make(ConfigMap)
	defaults["this.is.a.very.nested.config"] = true
	defaults["this.is.a.very.nested.config.with"] = "one"
	defaults["this.is.a.very.nested.config.with.second"] = int(2)
	return defaults
}

func createTestConfigFile(path, filename, data string) (string, error) {
	testpath := path + randomString(8)
	config := testpath + filename
	err := os.MkdirAll(testpath, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(config, []byte(data), 0644)
	if err != nil {
		fmt.Printf("error creating file")
		return "", err
	}
	return testpath, nil
}

// randomString return a random string
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
