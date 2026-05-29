# Ayotl

Ayotl is a lightweight Go library for loading configuration.
It supports **config files** (JSON / YAML / INI), **environment variable placeholders** within those files,
and **environment-only mode** for when no config file is needed — without forcing you
to implement any interface.

```go
import "github.com/beabys/ayotl"
```

---

## Quick Start

```go
// 1. Define your config struct with mapstructure tags
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type Config struct {
	Server ServerConfig `mapstructure:"server"`
}
//                                                 
// 2. Create a new config, load from file, unmarshal
cfg := config.New().WithEnv().LoadConfigs("config.json")
cfg.Unmarshal(&myConfig)
```

No struct methods to implement. No interface to satisfy.
Just set fields on `config.Config` and call load.

---

## Modes

### 1. Config File Mode — JSON / YAML / INI

Load values from a `.json`, `.yaml`, `.yml`, or `.ini` file:

```go
cfg := config.New().
    WithEnv().
    LoadConfigs("config.json")
```

Config files support **optional** `${PLACEHOLDER}` substitution.
Use `WithEnv()` to load env vars for placeholder resolution:

```json
{
    "stage": "${STAGE}",
    "app": {
        "host": "127.0.0.1",
        "port": "${APP_PORT}"
    }
}
```

```go
cfg := config.New().
    WithEnv().
    LoadConfigs("config.json")
// stage = os.Getenv("STAGE"), app.port = os.Getenv("APP_PORT")
```

Restrict which env vars are available for substitution for security:

```go
cfg := config.New().
    WithEnv("STAGE", "APP_PORT").  // ← only these two
    LoadConfigs("config.json")
```

If a referenced placeholder is not set (or filtered out), it resolves to empty string.

---

### 2. Environment-Only Mode — no files

Call `LoadConfigs()` with **no arguments**.
Env vars are loaded into `ConfigMap` via an explicit alias:

```go
cfg := config.New()
cfg.EnvAlias = config.ConfigEnvAlias{
    "SERVER_HOST": "server.host",
    "SERVER_PORT": "server.port",
}
cfg.WithEnv().LoadConfigs()

fmt.Println(cfg.MustString("server.host", ""))  // "localhost"
```

`EnvAlias` maps env var names to dot-notation config keys.
Only env vars listed in the alias are placed into `ConfigMap`.

> **Note:** `LoadConfigs()` calls `WithEnv()` internally if env vars
> haven't been loaded yet, so you can omit the explicit `WithEnv()` call:
>
> ```go
> cfg := config.New()
> cfg.EnvAlias = config.ConfigEnvAlias{
>     "SERVER_HOST": "server.host",
> }
> cfg.LoadConfigs()  // ← WithEnv called internally
> ```

---

## Default Values

Set fallback values via the `Defaults` field.
Defaults are applied **before** env alias overrides,
so env vars always take precedence:

```go
cfg := config.New()
cfg.Defaults = config.ConfigMap{
    "server.host":   "127.0.0.1",
    "server.port":   3000,
    "logger.level":  "info",
}
cfg.LoadConfigs()
```

Defaults and env alias can be combined freely.
Env vars that match an alias key override the default.
Keys not set by env keep their default.

---

## Constructor: `NewWithParams`

Batch-set `Defaults` and `EnvAlias` at construction time:

```go
cfg := config.NewWithParams(&config.Params{
    Defaults: config.ConfigMap{
        "server.port":  9090,
        "logger.level": "info",
    },
    EnvAlias: config.ConfigEnvAlias{
        "SERVER_HOST": "server.host",
        "SERVER_PORT": "server.port",
    },
})
cfg.WithEnv().LoadConfigs()
```

`NewWithParams` does not modify `ConfigMap` or `EnvConfigMap` —
it only sets `Defaults` and `EnvAlias`. Loading still happens
when you call `LoadConfigs()`.

---

## Immutability

Lock the config after loading to prevent accidental mutations at runtime:

```go
cfg := config.New()
cfg.Defaults = config.ConfigMap{"server.host": "original"}
cfg.LoadConfigs()

cfg.Immutable()

// All mutations become no-ops:
cfg.Set("server.host", "will-not-stick")
cfg.SetConfigMap(newMap)
cfg.WithEnv("SOME_SECRET")
```

`Immutable()` can be called **before** `LoadConfigs()` —
loading still works (env vars and defaults are applied),
but all post-load writes are blocked:

```go
cfg := config.New()
cfg.EnvAlias = config.ConfigEnvAlias{
    "SERVER_HOST": "server.host",
}
cfg.Immutable().
    LoadConfigs()         // ← still works

cfg.Set("server.host", "blocked")  // ← no-op
```

After `Immutable()`, only `Get`, `Must*`, and `Unmarshal` remain available.

---

## Unmarshal

After loading, decode the resolved `ConfigMap` into your typed struct:

```go
type Config struct {
    Server ServerConfig `mapstructure:"server"`
}
type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

var myConfig Config
cfg := config.New().WithEnv().LoadConfigs("config.json")
if err := cfg.Unmarshal(&myConfig); err != nil {
    log.Fatal(err)
}
fmt.Println(myConfig.Server.Host)
```

All fields need `mapstructure` struct tags to be recognized.

---

## Access values directly

Use the `Must*` helpers to read values by dot-notation key:

```go
host  := cfg.MustString("server.host", "localhost")
port  := cfg.MustInt("server.port", 3000)
debug := cfg.MustBool("server.debug", false)
```

`Must*` checks `EnvConfigMap` first (env vars loaded by `WithEnv`),
then falls back to `ConfigMap`. If neither has the key, the default value is returned.

---

## Struct tag requirement

All configurable fields **must** have a `mapstructure` struct tag.
Fields without a tag are ignored in both file and env-only modes.
