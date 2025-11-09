# Ayotl

Ayotl is a lightweight library designed to simplify the management of local configuration files by mapping their contents into well-defined Go structs.


## Config Files

### supported config files
- yaml
- json

### Example of a config json file:

```json
{
    "stage": "development",
    "app": {
        "host": "127.0.0.1",
        "port": 3001
    },
    "services": {
        "login": {
            "host": "127.0.0.1",
            "port": 3002,
            "user": "123",
            "password": "4567",
        },
        "enabled": true,
    }
}
```
## Environment variables
Configuration values can be sourced from environment variables using placeholders in your config file:

```json
{
    "stage": "${STAGE}",
    "app": {
        "host": "${APPLICATION_HOST}",
        "port": "${APPLICATION_PORT}"
    },
   
    "services": {
        "login": {
            "host": "${LOGIN_SERVICE_HOST}",
            "port": "${LOGIN_SERVICE_PORT}",
            "user": "${LOGIN_SERVICE_USER}",
            "password": "${LOGIN_SERVICE_PASSWORD}",
        },
        "enabled": true,
}
```

Set your environment variables as follows:

```bash
STAGE=development
APPLICATION_PORT=3001
APPLICATION_HOST=127.0.0.1
LOGIN_SERVICE_HOST=127.0.0.1
LOGIN_SERVICE_PORT=3002
LOGIN_SERVICE_USER=1234
LOGIN_SERVICE_PASSWORD=5678
```

If an environment variable is not set, its value will default to an empty string.

## Integration

Define your configuration struct using the `mapstructure` struct tag. For example:

```go
type Config struct {
	Stage    string        `mapstructure:"stage"`
	App      App           `mapstructure:"app"`
	Services Services      `mapstructure:"services"`
}

// ApplicationConfig is a struct to define configurations for the http server
type App struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LoggerConfig is a struct to define configurations for Logger
type Services struct {
	Login   LoginService `mapstructure:"login"`
}
// LoggerConfig is a struct to define configurations for Logger
type LoginService struct {
	Host string     `mapstructure:"host"`
	Port int        `mapstructure:"port"`
	User string     `mapstructure:"user"`
	Pass string     `mapstructure:"password"`
}
```

### Default Values

To specify default values for configuration fields, implement the `SetDefaults` function:

```go
func (c *Config) SetDefaults() ConfigMap {
    // create a new configMap
	defaults := make(ConfigMap)
	// application defaults values
	defaults["services.login.host"] = "127.0.0.1"
	defaults["services.login.port"] = "3002"
    // return the default configMap with our mapping
	return defaults
}
```
Using the `dot-notation` can set the default value, this will be appended in the struct if this config value doesn't exist in our config file.

if those values are defied in our config file, those will be overridden for the one existing in the config file
