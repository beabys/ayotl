package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestReadFileJSON(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.json", `{"a":"b","nested":{"x":"y"}}`)

	m, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned unexpected error: %v", err)
	}

	v, ok := m["a"].(string)
	if !ok || v != "b" {
		t.Fatalf("expected a=b, got %#v (ok=%v)", v, ok)
	}

	nested, ok := m["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested to be map[string]interface{}, got %#v", m["nested"])
	}
	if nx, _ := nested["x"].(string); nx != "y" {
		t.Fatalf("expected nested.x = y, got %#v", nested["x"])
	}
}

func TestReadFileYAMLAndYML(t *testing.T) {
	dir := t.TempDir()
	yamlContent := "a: b\nnested:\n  x: y\n"

	// .yaml
	pathYaml := writeTempFile(t, dir, "test.yaml", yamlContent)
	m1, err := ReadFile(pathYaml)
	if err != nil {
		t.Fatalf("ReadFile(.yaml) returned unexpected error: %v", err)
	}
	if v, _ := m1["a"].(string); v != "b" {
		t.Fatalf(".yaml: expected a=b, got %#v", m1["a"])
	}

	// .yml
	pathYml := writeTempFile(t, dir, "test.yml", yamlContent)
	m2, err := ReadFile(pathYml)
	if err != nil {
		t.Fatalf("ReadFile(.yml) returned unexpected error: %v", err)
	}
	if v, _ := m2["a"].(string); v != "b" {
		t.Fatalf(".yml: expected a=b, got %#v", m2["a"])
	}
}

func TestReadFileINI(t *testing.T) {
	dir := t.TempDir()
	iniContent := `[server]
host = localhost
port = 8080

[logger]
level = debug
`
	path := writeTempFile(t, dir, "test.ini", iniContent)

	m, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned unexpected error: %v", err)
	}

	server, ok := m["server"].(ConfigMap)
	if !ok {
		t.Fatalf("expected server to be ConfigMap, got %#v", m["server"])
	}
	if v, _ := server["host"].(string); v != "localhost" {
		t.Fatalf("expected server.host = localhost, got %#v", server["host"])
	}
	if v, _ := server["port"].(string); v != "8080" {
		t.Fatalf("expected server.port = 8080, got %#v", server["port"])
	}

	logger, ok := m["logger"].(ConfigMap)
	if !ok {
		t.Fatalf("expected logger to be ConfigMap, got %#v", m["logger"])
	}
	if v, _ := logger["level"].(string); v != "debug" {
		t.Fatalf("expected logger.level = debug, got %#v", logger["level"])
	}
}

func TestReadFileINIMalformed(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.ini", "[server\nhost = localhost")

	_, err := ReadFile(path)
	if err == nil {
		t.Fatalf("expected ini decode error, got nil")
	}
}

func TestReadFileInvalidExtension(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.txt", "a: b")

	_, err := ReadFile(path)
	if err == nil {
		t.Fatalf("expected error for invalid extension, got nil")
	}
	if !strings.Contains(err.Error(), "invalid extension type") {
		t.Fatalf("expected invalid extension error, got: %v", err)
	}
}

func TestReadFileNonExistent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doesnotexist.json")

	_, err := ReadFile(path)
	if err == nil {
		t.Fatalf("expected error for non-existent file, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected a not-exist error, got: %v", err)
	}
}

func TestReadFileMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.json", `{"a":`)

	_, err := ReadFile(path)
	if err == nil {
		t.Fatalf("expected json decode error, got nil")
	}
	// error message can vary; ensure it's a json unmarshal related error
	if !strings.Contains(err.Error(), "unexpected end of JSON input") && !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("expected json unmarshal error, got: %v", err)
	}
}
