package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadFile is a function to read a file and decode,
// supporting only yaml and json at the moment
func ReadFile(file string) (ConfigMap, error) {
	content, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	ext := getFileExt(file)
	dataMap := make(ConfigMap)

	switch true {
	case ext == "json":
		err = jsonDecode(content, &dataMap)
	case ext == "yaml" || ext == "yml":
		err = yamlDecode(content, &dataMap)
	default:
		err = fmt.Errorf("invalid extension type: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	return dataMap, nil
}

func jsonDecode(j []byte, d *ConfigMap) error {
	return json.Unmarshal(j, d)
}

func yamlDecode(j []byte, d *ConfigMap) error {
	return yaml.Unmarshal(j, d)
}

func getFileExt(s string) (ext string) {
	fullExt := filepath.Ext(s)
	return strings.TrimPrefix(fullExt, ".")
}
