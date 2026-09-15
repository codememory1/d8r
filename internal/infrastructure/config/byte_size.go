package config

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-yaml"
)

type ByteSize int64

func (s *ByteSize) Bytes() int64 {
	return int64(*s)
}

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
