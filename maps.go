package config

import (
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// MergeKeys merge 2 map[string]interface{} given
func MergeKeys(m1, m2 map[string]interface{}) map[string]interface{} {
	for key, m2Val := range m2 {
		// first we validate if key exist
		m1Val, ok := m1[key]
		if !ok {
			// If key no exist , add it and continue
			m1[key] = m2Val
			continue
		}
		switch v := m1Val.(type) {
		case map[string]interface{}:
			// Recursive Call
			m1[key] = MergeKeys(v, m2Val.(map[string]interface{}))
		default:
			m1[key] = m2Val
		}
	}
	return m1
}

// MergeEnvVar merge Env variables into placeholders
func MergeEnvVar(m map[string]interface{}) map[string]interface{} {
	for key, val := range m {
		switch value := val.(type) {
		case map[string]interface{}:
			// Recursive Call
			m[key] = MergeEnvVar(value)
		case string:
			if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
				p := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
				m[key] = os.Getenv(p)
			}
		default:
			continue
		}
	}
	return m
}

// Flatten  is a init wrapper for flatten
func Flatten(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	out = flatten(m, nil, out)
	return out
}

// flatten  is a recursive function to convert map[string]interface{} into a dot-notation
func flatten(m map[string]interface{}, keys []string, out map[string]interface{}) map[string]interface{} {
	for key, val := range m {
		// Copy the incoming key paths into a new map
		// and append the current key in the iteration.
		keyPaths := make([]string, 0, len(keys)+1)
		// append the keys received
		keyPaths = append(keyPaths, keys...)
		// append the new key
		keyPaths = append(keyPaths, key)

		// verify the type of the current value
		switch cur := val.(type) {
		case map[string]interface{}:
			// Empty map. only add as is it
			if len(cur) == 0 {
				newKey := strings.Join(keyPaths, ".")
				out[newKey] = val
				continue
			}

			// Recursive call if value is not empty
			out = flatten(cur, keyPaths, out)
		default:
			newKey := strings.Join(keyPaths, ".")
			out[newKey] = val
		}
	}
	return out
}

// GetValue is a function to search recursively a key in a map[string]interface{}
func GetValue(m map[string]interface{}, keysToFind []string) interface{} {
	lenkeysToFind := len(keysToFind)
	// if keysToFind is empty we return empty string
	if lenkeysToFind < 1 {
		return nil
	}
	// takes the first postion
	keyVal := keysToFind[0]
	next := keysToFind[1:]
	val, ok := m[keyVal]
	if !ok {
		return nil
	}
	// validate the type, if still a map[string]interface{}
	// we should do a recursive call
	switch v := val.(type) {
	case map[string]interface{}:
		// Recusrive call
		return GetValue(v, next)
	default:
		// if 'next' has keys inside means, the key don't exist
		if len(next) < 1 {
			return val
		}
	}
	return nil
}

// SetValue is a function to search recursively a key in a map[string]interface{} and add it with a set of keys given
func SetValue(m map[string]interface{}, keysToFind []string, value interface{}) map[string]interface{} {
	lenkeysToFind := len(keysToFind)
	// if keysToFind is empty we return empty string
	if lenkeysToFind < 1 {
		return m
	}
	// takes the first position
	keyVal := keysToFind[0]
	next := keysToFind[1:]
	val, ok := m[keyVal]
	if !ok {
		// by default set the key as the value
		val = value
		// if key not exist, but have more keys to search
		if len(next) > 0 {
			val = make(map[string]interface{})
		}
	}
	// validate the type, if still a map[string]interface
	// we should do recursion call, in not validate if still
	// theres keys to find
	switch v := val.(type) {
	case map[string]interface{}:
		// Recusrive call
		m[keyVal] = SetValue(v, next, value)
	default:
		// if still has more keys, override the current value with
		// with a new nested map[string]interface{}
		if len(next) > 0 {
			val = SetValue(make(map[string]interface{}), next, value)
		}
		// create the new key in the map
		m[keyVal] = val

	}
	return m
}

// Unmarshall function convert a map[string]interface{} into a struct using mapstructure from external package
func UnMarshall(configMap ConfigMap, out *interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           out,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(configMap)
}
