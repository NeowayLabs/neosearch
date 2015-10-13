package neosearch

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"gopkg.in/yaml.v2"
)

// Option is a internal type used for setting config options.
// This need be exported only because the usage of the
// self-referential functions in the option assignments.
// More details on the Rob Pike blog post below:
// http://commandcenter.blogspot.com.br/2014/01/self-referential-functions-and-design.html
type Option func(c *Config) Option

// Config stores NeoSearch configurations
type Config struct {
	// Root directory where all of the indices will be written.
	DataDir string `yaml:"dataDir"`

	// Enables debug in every sub-module
	Debug bool `yaml:"debug"`

	// IndicesCacheSize is the max number of indices maintained open
	MaxIndicesOpen int `yaml:"maxIndicesOpen"`

	// Specific options to store
	KVConfig store.KVConfig `yaml:"store"`
}

// NewConfig creates new config
func NewConfig() *Config {
	return &Config{}
}

// Option configures the config struct
func (c *Config) Option(opts ...Option) (previous Option) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

// Debug enables or disable debug
func Debug(t bool) Option {
	return func(c *Config) Option {
		previous := c.Debug
		c.Debug = t
		return Debug(previous)
	}
}

// DataDir set the data directory for neosearch database and indices
func DataDir(path string) Option {
	return func(c *Config) Option {
		previous := c.DataDir
		c.DataDir = path
		return DataDir(previous)
	}
}

// MaxIndicesOpen set the maximum number of open indices
func MaxIndicesOpen(size int) Option {
	return func(c *Config) Option {
		previous := c.MaxIndicesOpen
		c.MaxIndicesOpen = size
		return MaxIndicesOpen(previous)
	}
}

// Specific configurations to store
func KVConfig(kvconfig store.KVConfig) Option {
	return func(c *Config) Option {
		previous := c.KVConfig
		c.KVConfig = kvconfig
		return KVConfig(previous)
	}
}

// ConfigFromFile loads configuration from YAML file
func ConfigFromFile(filename string) (*Config, error) {
	// Load config from file
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Failed to open config file: %s", err.Error())
	}

	fileContent, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	cfg := NewConfig()

	err = yaml.Unmarshal(fileContent, &cfg)
	return cfg, err
}
