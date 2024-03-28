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
		config := New().IsMergeEnv(true).SetConfigImpl(mock).SetConfigMap(configMap)
		var want = &Config{true, configMap, mock}
		areEqual := assert.ObjectsAreEqual(config, want)
		assert.True(t, areEqual)
	})

	t.Run("test Error Loading configs path without CONFIG_FILE", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		mock := &MockConfig{}
		config := New().IsMergeEnv(true)
		assert.ErrorContains(t, config.LoadConfigs(mock), "CONFIG_FILE")
	})

	t.Run("test Error Loading configs file", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock)
		os.Setenv("CONFIG_FILE", "./../env.configuration.json")
		assert.ErrorContains(t, config.LoadConfigs(mock), "Fail to load configs")
		os.Unsetenv("CONFIG_FILE")
	})

	t.Run("Test Loading configs", func(t *testing.T) {
		mock := &MockConfig{}
		config := New().IsMergeEnv(true).SetConfigImpl(mock)
		data := `{"server": {"enabled": false},"second": {"config": {"enabled": false}}}`
		testPath, err := createTestConfigFile(path, "/config.json", data)
		assert.NoError(t, err)
		os.Setenv("CONFIG_FILE", testPath+"/config.json")
		assert.NoError(t, config.LoadConfigs(mock))
		os.Unsetenv("CONFIG_FILE")
	})

	t.Run("Test Loading configs with config placeholders", func(t *testing.T) {
		mock := &MockConfig{}
		data := `{"server": {"enabled": "${IS_CONFIG_FOR_TEST_ENABLED}"},"second": {"config": {"enabled": false}}}`
		testPath, err := createTestConfigFile(path, "/config.json", data)
		assert.NoError(t, err)
		os.Setenv("IS_CONFIG_FOR_TEST_ENABLED", "true")
		os.Setenv("CONFIG_FILE", testPath+"/config.json")
		c := &Config{}
		config := c.IsMergeEnv(true).SetConfigImpl(mock)
		assert.NoError(t, config.LoadConfigs(mock))
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("IS_CONFIG_FOR_TEST_ENABLED")
	})

	t.Run("Test Loading configs with config placeholders and no Value", func(t *testing.T) {
		mock := &MockConfig{}
		data := `{"server": {"enabled": "${IS_CONFIG_FOR_TEST_ENABLED_SECOND}"},"second": {"config": {"enabled": false}}}`
		testPath, err := createTestConfigFile(path, "/config.json", data)
		assert.NoError(t, err)
		os.Setenv("CONFIG_FILE", testPath+"/config.json")
		c := &Config{}
		config := c.IsMergeEnv(true).SetConfigImpl(mock)
		assert.NoError(t, config.LoadConfigs(mock))
		os.Unsetenv("CONFIG_FILE")
	})
}

// func TestFetchAWSSecretes(t *testing.T) {

// 	path := "./../../testConfigAws/"
// 	defer os.RemoveAll(path)
// 	t.Run("Test AWS config", func(t *testing.T) {
// 		os.Setenv("AWS_REGION", "eu-west-1")
// 		mock := &MockConfig{}
// 		config := New().SetConfigImpl(mock)
// 		awsConfig := config.SetAWSConfigs()
// 		wants := &aws.Config{Region: aws.String("eu-west-1")}
// 		areEqual := assert.ObjectsAreEqual(awsConfig, wants)
// 		assert.True(t, areEqual)
// 	})

// 	// t.Run("Test Fail with empty data AWS Secrets", func(t *testing.T) {
// 	// 	mock := &MockConfig{}
// 	// 	fetcher := &MockFetcher{Content: []byte("")}
// 	// 	config := New().SetConfigImpl(mock).SetFetcher(fetcher)
// 	// 	os.Setenv("ENVIRONMENT_CONFIG", "myconfig")
// 	// 	os.Setenv("ENVIRONMENT_STORAGE", "AWS_SECRET")
// 	// 	assert.ErrorContains(t, config.LoadConfigs(mock), "error parsing json from secrets")
// 	// 	os.Unsetenv("ENVIRONMENT_STORAGE")
// 	// 	os.Unsetenv("ENVIRONMENT_CONFIG")
// 	// })

// 	t.Run("Test Fail fetch from AWS Secrets", func(t *testing.T) {
// 		mock := &MockConfig{}
// 		fetcher := &MockFetcherError{Content: []byte("")}
// 		config := New().SetConfigImpl(mock).SetFetcher(fetcher)
// 		os.Setenv("ENVIRONMENT_CONFIG", "myconfig")
// 		os.Setenv("ENVIRONMENT_STORAGE", "AWS_SECRET")
// 		assert.ErrorContains(t, config.LoadConfigs(mock), "error fetching general config")
// 		os.Unsetenv("ENVIRONMENT_STORAGE")
// 		os.Unsetenv("ENVIRONMENT_CONFIG")
// 	})

// 	t.Run("Test Success fetch from AWS Secrets", func(t *testing.T) {
// 		mock := &MockConfig{}
// 		randomFolderName := path + randomString(8)
// 		result := []byte("{\"application.port\": 8080,\"logger.log_errors_to\": \"stderr\",\"logger.log_output_to\": \"stderr\",\"logger.log_level\": \"log\"}")
// 		result2 := []byte("{\"logger.log_errors_to\": \"stdout\",\"logger.log_output_to\": \"stdout\",\"logger.log_level\": \"info\"}")
// 		fetcher := &MockFetcher2{Content1: result, Content2: result2}
// 		config := New().SetConfigImpl(mock).SetFetcher(fetcher)
// 		data := `{"grpc_server": {"enabled": false},"apollo_product": {"consumer": {"enabled": false}}}`
// 		errCreatingDir := os.MkdirAll(randomFolderName, os.ModePerm)
// 		if errCreatingDir != nil {
// 			fmt.Printf("error creating Directorie(s)")
// 		}
// 		createFile(randomFolderName+"/config.json", data)
// 		os.Setenv("ENVIRONMENT_CONFIG", "myconfig,myconfig2")
// 		os.Setenv("ENVIRONMENT_STORAGE", "AWS_SECRET")
// 		os.Setenv("LOCAL_CONFIG", randomFolderName+"/config.json")
// 		err := config.LoadConfigs(mock)
// 		assert.Equal(t, mock.App.Port, int(8080))
// 		assert.Equal(t, mock.Logger.LogOutput, "stderr")
// 		assert.Equal(t, mock.Logger.Level, "log")
// 		assert.NoError(t, err)
// 		os.Unsetenv("LOCAL_CONFIG")
// 		os.Unsetenv("ENVIRONMENT_STORAGE")
// 		os.Unsetenv("ENVIRONMENT_CONFIG")
// 	})
// }

func TestMustFunctions(t *testing.T) {
	t.Run("test MustBool", func(t *testing.T) {
		mock := &MockConfig{}
		os.Setenv("ANY", "false")
		config := New().SetConfigImpl(mock).LoadEnv()
		required := true
		val := config.MustBool("ANY", required)
		assert.False(t, val)
		os.Unsetenv("ANY")
	})
	t.Run("test MustBool exist", func(t *testing.T) {
		os.Setenv("ANY", "true")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		required := true
		val := config.MustBool("ANY", required)
		assert.True(t, val)
		os.Unsetenv("ANY")
	})
	t.Run("test MustString", func(t *testing.T) {
		os.Setenv("ANY", "required")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		required := "required"
		val := config.MustString("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt", func(t *testing.T) {
		os.Setenv("ANY", "0")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		required := 0
		val := config.MustInt("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt32", func(t *testing.T) {
		os.Setenv("ANY", "0")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		integer := 0
		required := int32(integer)
		val := config.MustInt32("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt64", func(t *testing.T) {
		os.Setenv("ANY", "0")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		integer := 0
		required := int64(integer)
		val := config.MustInt64("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustString with variable", func(t *testing.T) {
		os.Setenv("ANY", "required")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		required := "required"
		val := config.MustString("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt with variable", func(t *testing.T) {
		os.Setenv("ANY", "10")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		required := 10
		val := config.MustInt("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt32 with variable", func(t *testing.T) {
		os.Setenv("ANY", "10")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
		integer := 10
		required := int32(integer)
		val := config.MustInt32("ANY", required)
		assert.Equal(t, val, required)
		os.Unsetenv("ANY")
	})

	t.Run("test MustInt64 with variable", func(t *testing.T) {
		os.Setenv("ANY", "10")
		mock := &MockConfig{}
		config := New().SetConfigImpl(mock).LoadEnv()
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
