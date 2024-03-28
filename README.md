# Config

Config is a simple implementation intended to reduce the complexity to manage local configuration files, mapping into defined structs.


## Config Files

### supported config files
- yaml
- json

Example of a config json file:

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
Environment variables depend on a [config file](#config-files).
you can create your config file with placeholders similar to this:
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
assuming your environment variables are like this:
```bash
STAGE=development
APPLICATION_PORT=3001
APPLICATION_HOST=127.0.0.1
LOGIN_SERVICE_HOST=127.0.0.1
LOGIN_SERVICE_PORT=3002
LOGIN_SERVICE_USER=1234
LOGIN_SERVICE_PASSWORD=5678
```
If any environment doesn't exist or can not be found, this will be replaced as an empty string.

## Integration

You should have a struct with a struct tag `mapstructure`

following with the example of

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

You can define a struct similar to:
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
Sometimes is required to have default values in case we don't need to set them inside a config file
for those cases, you should implement those default values in the function `Setdefaults`.

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

if those values are defied in our config file, those will be overrided for the one existing in the config file

### Direct Mapping using environment variables
You can define a struct similar to:
```go
type Config struct {
	Stage    string    `mapstructure:"APPLICATION_STAGE"`
}
```
Assuming your environment variables are like this:
```bash
APPLICATION_STAGE=development
```
