package badm

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	filePermission   = 0644
	defaultKeyType   = "string"
	defaultValueType = "hex"
)

// Configuration handles the badm configuration file.
type Configuration struct {
	SelectedPath string             `json:"selected_path"`
	Buckets      map[string]*Bucket `json:"buckets"`
	Plugins      []*Plugin          `json:"plugins"`
}

// Bucket handles the bucket configuration.
type Bucket struct {
	KeyType   string `json:"key_type,omitempty"`
	ValueType string `json:"value_type,omitempty"`
}

// Plugin hanles the plugin configuration.
type Plugin struct {
	Path string `json:"path"`
}

func updateConfiguration(path string, fn func(*Configuration) error) error {
	c, err := readConfiguration(path)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	if err := fn(c); err != nil {
		return err
	}

	if err := writeConfiguration(path, c); err != nil {
		return fmt.Errorf("write configuration: %v", err)
	}

	return nil
}

func readConfiguration(path string) (*Configuration, error) {
	c := &Configuration{}

	file, err := os.OpenFile(path, os.O_RDONLY, filePermission)
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(c); err != nil {
		return nil, fmt.Errorf("decode json: %v", err)
	}

	return c, nil
}

func writeConfiguration(path string, c *Configuration) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePermission)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("encode json: %v", err)
	}

	return nil
}

// KeyType returns the key type for the provided bucket. The default type is returned,
// if no key type is specified.
func (c *Configuration) KeyType(name string) string {
	if c.Buckets == nil {
		return defaultKeyType
	}
	if b, ok := c.Buckets[name]; ok && b.KeyType != "" {
		return b.KeyType
	}
	return defaultKeyType
}

// ValueType returns the value type for the provided bucket. The default type is returned,
// if no value type is specified.
func (c *Configuration) ValueType(name string) string {
	if c.Buckets == nil {
		return defaultValueType
	}
	if b, ok := c.Buckets[name]; ok && b.ValueType != "" {
		return b.ValueType
	}
	return defaultValueType
}
