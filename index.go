// Package neosearch is a full-text-search library focused in fast search on multiple
// indices, doing data join between queries.
//
// Create and add documents to NeoSearch is very easy, see below:
//
//   func main() {
//       config := neosearch.NewConfig()
//       config.Option(neosearch.DataDir("/tmp"))
//       config.Option(neosearch.Debug(false))
//
//       neo := neosearch.New(config)
//
//       index, err := neosearch.CreateIndex("test")
//
//       if err != nil {
//           panic(err)
//       }
//
//       err = index.Add(1, `{"name": "Neoway Business Solution", "type": "company"}`)
//
//       if err != nil {
//           panic(err)
//       }
//
//       err = index.Add(2, `{"name": "Facebook Inc", "type": "company"}`)
//
//       if err != nil {
//           panic(err)
//       }
//
//       values, err := index.MatchPrefix([]byte("name"), []byte("neoway"))
//
//       if err != nil {
//           panic(err)
//       }
//
//       for _, value := range values {
//           fmt.Println(value)
//       }
//   }
//
// NeoSearch supports the features below:
//
// * Create/Delete index
// * Index JSON documents (No schema)
// * Bulk writes
// * Analysers
//   - Tokenizer
// * Search
//   - MatchPrefix
//   - FilterTerm
//
// This project is in active development stage, it is not recommended for
// production environments.
package neosearch

import (
	"errors"

	"github.com/neowaylabs/neosearch/engine"
	"github.com/neowaylabs/neosearch/index"
)

// NeoSearch is the core of the neosearch package.
// This structure handles all of the user's interactions with the indices,
// like CreateIndex, DeleteIndex, UpdateIndex and others.
type NeoSearch struct {
	Indices []*index.Index

	config *Config
	engine *engine.Engine
}

// New creates the NeoSearch high-level interface.
// Use that for index/update/delete JSON documents.
func New(cfg *Config) *NeoSearch {
	if cfg.dataDir == "" {
		panic(errors.New("DataDir is required for NeoSearch interface"))
	}

	if cfg.dataDir[len(cfg.dataDir)-1] == '/' {
		cfg.dataDir = cfg.dataDir[0 : len(cfg.dataDir)-1]
	}

	if cfg.cacheSize == 0 && cfg.enableCache {
		cfg.cacheSize = 3 << 30
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
			DataDir:     neo.config.dataDir,
			Debug:       neo.config.debug,
			CacheSize:   neo.config.cacheSize,
			EnableCache: neo.config.enableCache,
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
			DataDir:     neo.config.dataDir,
			Debug:       neo.config.debug,
			CacheSize:   neo.config.cacheSize,
			EnableCache: neo.config.enableCache,
		},
		false,
	)

	if err != nil {
		return nil, err
	}

	neo.Indices = append(neo.Indices, index)
	return index, nil
}

// Close all of the open indices
func (neo *NeoSearch) Close() {
	for idx := range neo.Indices {
		index := neo.Indices[idx]
		index.Close()
	}

	neo.Indices = make([]*index.Index, 0)
}
