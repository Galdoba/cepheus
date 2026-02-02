package tttable

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// Serialization defines supported data formats for table persistence.
type Serialization string

const (
	ToJSON       Serialization = "json"        // JSON format without indentation
	ToJSONIndent Serialization = "json_indent" // JSON format with indentation (human-readable)
	ToTOML       Serialization = "toml"        // TOML format (Tom's Obvious Minimal Language)
	ToYAML       Serialization = "yaml"        // YAML format (YAML Ain't Markup Language)
)

type Serializer interface {
	Serialize(Serialization) ([]byte, error)
}

// Serialize converts a Table to the specified serialization format.
// Returns the byte representation of the table or an error if the format is unsupported
// or marshaling fails.
func (t *Table) Serialize(method Serialization) ([]byte, error) {
	switch method {
	case ToJSON:
		return t.toJSON()
	case ToJSONIndent:
		return t.toJSONIndent()
	case ToTOML:
		return t.toTOML()
	case ToYAML:
		return t.toYAML()
	case "":
		return nil, fmt.Errorf("serialization method not provided")
	default:
		return nil, fmt.Errorf("unsupported serialization method (%v)", method)
	}
}

// Deserialize attempts to load table data from bytes in JSON, TOML, or YAML format.
// It tries each format sequentially and returns on first successful decoding.
// Returns an error if all supported formats fail to decode the data.
func (t *Table) Deserialize(data []byte) error {
	// var t Table
	var errs []string
	t.parsed = make(map[string]*RangeKey)
	t.Rows = make(map[string]TableEntry)
	t.ModsFlat = make(map[string]int)
	t.ModsCumulative = make(map[string]int)
	t.ModsMax = make(map[string]int)
	t.ModsMin = make(map[string]int)

	// Try JSON
	if err := json.Unmarshal(data, &t); err == nil {
		return nil
	} else {
		errs = append(errs, fmt.Sprintf("json: %v", err))
	}

	// Try TOML
	if err := toml.Unmarshal(data, &t); err == nil {
		return nil
	} else {
		errs = append(errs, fmt.Sprintf("toml: %v", err))
	}

	// Try YAML
	if err := yaml.Unmarshal(data, &t); err == nil {
		return nil
	} else {
		errs = append(errs, fmt.Sprintf("yaml: %v", err))
	}

	return fmt.Errorf("failed to deserialize data as any supported format:\n  - %s",
		strings.Join(errs, "\n  - "))
}

// SaveAs writes the table to a file, determining the format by file extension.
// Supported extensions: .json, .toml, .yaml, .yml.
// Returns an error if the format is unsupported, serialization fails, or file operations fail.
func SaveAs(s Serializer, path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	var method Serialization

	switch ext {
	case ".json":
		method = ToJSONIndent
	case ".toml":
		method = ToTOML
	case ".yaml", ".yml":
		method = ToYAML
	default:
		return fmt.Errorf("unsupported file extension %q, use .json, .toml, .yaml, or .yml", ext)
	}

	data, err := s.Serialize(method)
	if err != nil {
		return fmt.Errorf("serialization failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// Load creates a Table from a file containing JSON, TOML, or YAML data.
// The format is automatically detected. Returns an error if the file cannot be read
// or the data cannot be deserialized.
func Load(path string) (RollableTable, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	t := &Table{}
	err = t.Deserialize(data)
	if err == nil {
		return t, nil
	}
	return nil, fmt.Errorf("failed to deserealize data from %q", path)
}

// toJSON converts the table to compact JSON format.
func (t *Table) toJSON() ([]byte, error) {
	return json.Marshal(t)
}

// toJSONIndent converts the table to indented JSON format for human readability.
func (t *Table) toJSONIndent() ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}

// toTOML converts the table to TOML format.
func (t *Table) toTOML() ([]byte, error) {
	return toml.Marshal(t)
}

// toYAML converts the table to YAML format.
func (t *Table) toYAML() ([]byte, error) {
	return yaml.Marshal(t)
}
