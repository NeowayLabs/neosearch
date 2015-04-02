package neosearch

import (
	"io/ioutil"
	"log"
	"os"

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

	// CacheSize is the length of LRU cache used by the storage engine
	// Default is 1GB
	KVCacheSize int `yaml:"cacheSize"`

	// IndicesCacheSize is the max number of indices maintained open
	MaxIndicesOpen int `yaml:"maxIndicesOpen"`

	// EnableCache enable/disable cache support
	EnableCache bool `yaml:"enableCache"`
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

// EnableCache enables or disable index cache
func EnableCache(t bool) Option {
	return func(c *Config) Option {
		previous := c.EnableCache
		c.EnableCache = t

		return EnableCache(previous)
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

// KVCacheSize set the size of the cache for the LRU storage cache.
// The link below have more details on the cache mechanism of the
// leveldb (or any other LSM compatible):
// More details in the section "Performance" of the link below:
// http://htmlpreview.github.io/?https://github.com/google/leveldb/blob/master/doc/index.html
func KVCacheSize(size int) Option {
	return func(c *Config) Option {
		previous := c.KVCacheSize
		c.KVCacheSize = size

		return KVCacheSize(previous)
	}
}

func MaxIndicesOpen(size int) Option {
	return func(c *Config) Option {
		previous := c.MaxIndicesOpen
		c.MaxIndicesOpen = size

		return MaxIndicesOpen(previous)
	}
}

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
