package config

import (
	"errors"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-yaml"
)

type Config struct {
	HTTP     HTTP     `yaml:"http"`
	Postgres Postgres `yaml:"postgres"`
	Download Download `yaml:"httpdownload"`
	Workers  Workers  `yaml:"workers"`
}

type HTTP struct {
	Address      string        `yaml:"address"`
	Port         uint32        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Download struct {
	MinParallelSize  ByteSize `yaml:"min_parallel_size"`
	BufferSize       ByteSize `yaml:"buffer_size"`
	PartSize         ByteSize `yaml:"part_size"`
	MaxParallelParts int      `yaml:"max_parallel_parts"`
	RangeParts       int      `yaml:"range_parts"`
}

type Workers struct {
	Download DownloadWorker `yaml:"download"`
}

type DownloadWorker struct {
	Concurrency uint64 `yaml:"concurrency"`
}

func defaults() Config {
	return Config{
		HTTP: HTTP{
			Address:      "0.0.0.0",
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 20 * time.Second,
			IdleTimeout:  20 * time.Second,
		},
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "root",
			Password: "admin",
			Database: "postgres",
		},
		Download: Download{
			MinParallelSize:  ByteSize(50 * humanize.MiByte),  // 50MiB
			BufferSize:       ByteSize(256 * humanize.KiByte), // 256KiB
			PartSize:         ByteSize(5 * humanize.MiByte),   // 5MiB
			MaxParallelParts: 5,
			RangeParts:       5,
		},
		Workers: Workers{
			DownloadWorker{
				Concurrency: 5,
			},
		},
	}
}

// Load loads the configuration from the specified path and validates all parameters.
func Load(path string) (Config, error) {
	if path == "" {
		return Config{}, errors.New("configuration path is empty")
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	// A configuration is created with default parameters.
	configuration := defaults()

	if err := yaml.Unmarshal(data, &configuration); err != nil {
		return Config{}, err
	}

	if err := configuration.Validate(); err != nil {
		return Config{}, err
	}

	return configuration, nil
}

// Validate validates all configuration parameters.
func (c Config) Validate() error {
	if err := c.Postgres.Validate(); err != nil {
		return err
	}

	if err := c.Download.Validate(); err != nil {
		return err
	}

	if err := c.Workers.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates PostgreSQL configuration parameters.
func (c Postgres) Validate() error {
	if c.Host == "" {
		return errors.New("postgres.host is required")
	}

	if c.Port == 0 {
		return errors.New("postgres.port is required")
	}

	if c.User == "" {
		return errors.New("postgres.user is required")
	}

	if c.Database == "" {
		return errors.New("postgres.database is required")
	}

	return nil
}

// Validate validates Download configuration parameters.
func (c Download) Validate() error {
	if c.MinParallelSize.Bytes() <= 0 {
		return errors.New("httpdownload.min_parallel_size must be greater than 0")
	}

	return nil
}

// Validate orchestration and validation of all worker parameters.
func (c Workers) Validate() error {
	if err := c.Download.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates the download worker parameters.
func (c DownloadWorker) Validate() error {
	if c.Concurrency <= 0 {
		return errors.New("workers.download.concurrency must be greater than 0")
	}

	return nil
}
