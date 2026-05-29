package config

import (
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// Helper — create a temp config file
// ---------------------------------------------------------------------------

func benchWriteFile(b *testing.B, dir, name, content string) string {
	b.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		b.Fatal(err)
	}
	return path
}

// ---------------------------------------------------------------------------
// ayotl benchmarks
// ---------------------------------------------------------------------------

func BenchmarkAyotlNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := New()
		_ = c
	}
}

func BenchmarkAyotlLoadEnvOnly(b *testing.B) {
	b.ReportAllocs()
	os.Setenv("BENCH_HOST", "localhost")
	os.Setenv("BENCH_PORT", "8080")
	os.Setenv("BENCH_DEBUG", "true")

	alias := ConfigEnvAlias{
		"BENCH_HOST":  "server.host",
		"BENCH_PORT":  "server.port",
		"BENCH_DEBUG": "server.debug",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.EnvAlias = alias
		c.WithEnv()
		c.LoadConfigs()
	}
}

func BenchmarkAyotlWithEnvAll(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.WithEnv()
	}
}

func BenchmarkAyotlWithEnvFiltered(b *testing.B) {
	b.ReportAllocs()
	os.Setenv("BENCH_FILTERED", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.WithEnv("BENCH_FILTERED")
	}
}

func BenchmarkAyotlLoadJSON(b *testing.B) {
	b.ReportAllocs()
	dir := b.TempDir()
	path := benchWriteFile(b, dir, "cfg.json",
		`{"server":{"host":"localhost","port":8080,"debug":true},"logger":{"level":"debug"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		if err := c.LoadConfigs(path); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAyotlLoadYAML(b *testing.B) {
	b.ReportAllocs()
	dir := b.TempDir()
	path := benchWriteFile(b, dir, "cfg.yaml",
		"server:\n  host: localhost\n  port: 8080\n  debug: true\nlogger:\n  level: debug\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		if err := c.LoadConfigs(path); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAyotlLoadINI(b *testing.B) {
	b.ReportAllocs()
	dir := b.TempDir()
	path := benchWriteFile(b, dir, "cfg.ini",
		"[server]\nhost = localhost\nport = 8080\ndebug = true\n[logger]\nlevel = debug\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		if err := c.LoadConfigs(path); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAyotlLoadJSONWithPlaceholders(b *testing.B) {
	b.ReportAllocs()
	os.Setenv("BENCH_PLACEHOLDER", "overridden")
	dir := b.TempDir()
	path := benchWriteFile(b, dir, "cfg.json",
		`{"server":{"host":"localhost","port":"${BENCH_PLACEHOLDER}"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.WithEnv()
		if err := c.LoadConfigs(path); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAyotlUnmarshal(b *testing.B) {
	b.ReportAllocs()
	type Server struct {
		Host  string `mapstructure:"host"`
		Port  int    `mapstructure:"port"`
		Debug bool   `mapstructure:"debug"`
	}
	type Cfg struct {
		Server Server `mapstructure:"server"`
	}

	dir := b.TempDir()
	path := benchWriteFile(b, dir, "cfg.json",
		`{"server":{"host":"localhost","port":8080,"debug":true}}`)
	c := New()
	c.LoadConfigs(path)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out Cfg
		if err := c.Unmarshal(&out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAyotlMustString(b *testing.B) {
	b.ReportAllocs()
	dir := b.TempDir()
	path := benchWriteFile(b, dir, "cfg.json",
		`{"server":{"host":"localhost","port":8080}}`)
	c := New()
	c.LoadConfigs(path)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.MustString("server.host", "")
	}
}

// ---------------------------------------------------------------------------
// Memory allocation breakdown
// ---------------------------------------------------------------------------

func BenchmarkAyotlMemoryBreakdown(b *testing.B) {
	b.ReportAllocs()

	alias := ConfigEnvAlias{
		"BENCH_HOST": "server.host",
		"BENCH_PORT": "server.port",
	}
	defaults := ConfigMap{
		"server.host": "default",
		"server.port": 3000,
	}

	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			New()
		}
	})

	b.Run("New+Defaults", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := New()
			c.Defaults = defaults
			c.EnvAlias = alias
		}
	})

	b.Run("New+Defaults+Load", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c := New()
			c.Defaults = defaults
			c.EnvAlias = alias
			c.LoadConfigs()
		}
	})
}
