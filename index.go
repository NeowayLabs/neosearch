// Package neosearch is a full-text-search library focused in fast search on multiple
// indices, doing data join between queries.
package neosearch

import (
	"errors"

	"bitbucket.org/i4k/neosearch/engine"
	"bitbucket.org/i4k/neosearch/index"
)

// NeoSearch is the core of the neosearch package.
// This structure handles all of the user's interactions with the indices,
// like CreateIndex, DeleteIndex, UpdateIndex and others.
type NeoSearch struct {
	Indices []*index.Index

	config Config
	engine *engine.Engine
}

// Config stores NeoSearch configurations
type Config struct {
	// Root directory where all of the indices will be written.
	DataDir string

	// Enables debug in every sub-module
	Debug bool

	// CacheSize is the length of LRU cache used by the storage engine
	// Default is 1GB
	CacheSize int

	// EnableCache enable/disable cache support
	EnableCache bool
}

// New creates the NeoSearch high-level interface.
// Use that for index/update/delete JSON documents.
func New(cfg Config) *NeoSearch {
	if cfg.DataDir == "" {
		panic(errors.New("DataDir is required for NeoSearch interface"))
	}

	if cfg.DataDir[len(cfg.DataDir)-1] == '/' {
		cfg.DataDir = cfg.DataDir[0 : len(cfg.DataDir)-1]
	}

	if cfg.CacheSize == 0 && cfg.EnableCache {
		cfg.CacheSize = 3 << 30
	}

	neo := &NeoSearch{
		config: cfg,
	}

	return neo
}

// CreateIndex creates and setup a new index
func (neo *NeoSearch) CreateIndex(name string) (*index.Index, error) {
	index, err := index.New(
		name,
		index.Config{
			DataDir:     neo.config.DataDir,
			Debug:       neo.config.Debug,
			CacheSize:   neo.config.CacheSize,
			EnableCache: neo.config.EnableCache,
		},
		true,
	)

	if err != nil {
		return nil, err
	}

	neo.Indices = append(neo.Indices, index)
	return index, nil
}

// OpenIndex open a existing index for read/write operations.
func (neo *NeoSearch) OpenIndex(name string) (*index.Index, error) {
	index, err := index.New(
		name,
		index.Config{
			DataDir:     neo.config.DataDir,
			Debug:       neo.config.Debug,
			CacheSize:   neo.config.CacheSize,
			EnableCache: neo.config.EnableCache,
		},
		false,
	)

	if err != nil {
		return nil, err
	}

	neo.Indices = append(neo.Indices, index)
	return index, nil
}

// Close all of the opened indices
func (neo *NeoSearch) Close() {
	for idx := range neo.Indices {
		index := neo.Indices[idx]
		index.Close()
	}

	neo.Indices = make([]*index.Index, 0)
}
