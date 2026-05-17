# Ayotl

Ayotl is a lightweight Go library for loading configuration into well-defined structs.
It supports **config files** (JSON/YAML), **environment variable placeholders** within those files, and **pure environment-variable mode** when no config file is needed.

---

## Quick Start

```go
import "github.com/beabys/ayotl"

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}
```

## Modes

### 1. Config File Mode

Load values from a `.json`, `.yaml`, or `.yml` file:

```go
config := config.New().
    SetConfigImpl(&myConfig).
    LoadConfigs("config.json")
```

### 2. Environment-Only Mode (no file)

Call `LoadConfigs()` with **no arguments**.
The library walks your struct's `mapstructure` tags and reads matching env vars directly from the OS environment — **no need to call `.WithEnv()`**:

| Struct path             | Env var              |
|------------------------|----------------------|
| `server.host`          | `SERVER_HOST`        |
| `server.port`          | `SERVER_PORT`        |
| `logger.level`         | `LOGGER_LEVEL`       |

```go
os.Setenv("SERVER_HOST", "localhost")
os.Setenv("SERVER_PORT", "8080")
os.Setenv("LOGGER_LEVEL", "debug")

config := config.New().
    SetConfigImpl(&myConfig).
    LoadConfigs()   // ← no files, no WithEnv needed
```

Env var values are automatically converted to the target field type (string, int, bool, etc.).

### 3. File Mode with Optional Placeholder Substitution

Config values can be literal or reference env vars with `${VAR_NAME}`.
The `${...}` notation is **optional** — you can mix literal values and placeholders freely in the same file:

```json
{
    "stage": "${STAGE}",
    "app": {
        "host": "127.0.0.1",
        "port": "${APPLICATION_PORT}"
    }
}
```

Here `host` is a hardcoded literal, while `stage` and `port` will be replaced from env vars.

**`.WithEnv()` is only needed for this mode.** It loads environment variables into an internal map so `${PLACEHOLDER}` values can be resolved. By default it loads **all** environment variables:

```go
config := config.New().
    SetConfigImpl(&myConfig).
    WithEnv().                      // ← needed for ${...} substitution
    LoadConfigs("config.json")
```

For security, restrict which env vars are loaded by passing explicit names.
Only those vars will be available for substitution:

```go
config := config.New().
    SetConfigImpl(&myConfig).
    WithEnv("STAGE", "APPLICATION_PORT").   // ← only these two
    LoadConfigs("config.json")
```

If a referenced env var is not set (or not included in the filter), the placeholder resolves to an empty string.

> **Note:** `.WithEnv()` is **not** used in environment-only mode (mode 2) — that mode reads `os.Getenv` directly via struct reflection.

---

## Default Values

Implement `SetDefaults()` on your struct to provide fallback values in dot-notation.
Defaults are applied only when the key is not already set (from file or env):

```go
func (c *Config) SetDefaults() ConfigMap {
    return ConfigMap{
        "server.host":   "127.0.0.1",
        "server.port":   "3000",
        "logger.level":  "info",
    }
}
```

---

## Unmarshal

After loading, unmarshal the resolved values into your struct:

```go
config := config.New().
    SetConfigImpl(&myConfig).
    LoadConfigs("config.json")

if err := config.Unmarshal(&myConfig); err != nil {
    log.Fatal(err)
}
```

## Access values directly

Use the `Must*` helpers to read values by dot-notation key:

```go
host := config.MustString("server.host", "localhost")
port := config.MustInt("server.port", 3000)
debug := config.MustBool("server.debug", false)
```

---

## Struct tag requirement

All configurable fields **must** have a `mapstructure` struct tag.
Fields without a tag are ignored in both file and env-only modes.
