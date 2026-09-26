package config

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-yaml"
)

// ByteSize represents a byte-size configuration value.
type ByteSize int64

// Bytes returns the configured size as a number of bytes.
func (s *ByteSize) Bytes() int64 {
	return int64(*s)
}

// UnmarshalYAML parses a human-readable byte size from YAML.
func (s *ByteSize) UnmarshalYAML(data []byte) error {
	var raw string

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode byte size: %w", err)
	}

	bytes, err := humanize.ParseBytes(raw)

	if err != nil {
		return err
	}

	*s = ByteSize(bytes)

	return nil
}
