// Package neosearch is a full-text-search library focused in fast search on
// multiple indices, doing data join between queries.
//
// NeoSearch - Neoway Full Text Search Index
//
// NeoSearch is a feature-limited full-text-search library with focus on
// indices relationships, its main goal is provide very fast JOIN
// operations between information stored on different indices.
// It's not a complete FTS (Full Text Search) engine, in the common sense, but
// aims to solve very specific problems of FTS. At the moment, NeoSearch is a
// laboratory for research, not recommended for production usage, here we will
// test various technology for fast storage and search algorithms. In the
// future, maybe, we can proud of a very nice tech for solve search in big
// data companies.
//
// NeoSearch is like a Lucene library but without all of the complexities of
// complete FTS engine, written in Go, focusing on high performance search
// with data relationships.
//
// It's not yet complete, still in active development, then stay tuned for
// updates.
//
// Dependencies
//
//     - leveldb
//     - snappy (optional, only required for compressed data)
//     - Go > 1.3
//
// Install
//
//     git clone git@bitbucket.org:i4k/neosearch.git
//     cd neosearch
//     go get -u -v .
//     go build -tags leveldb -v .
//     go test -tags leveldb -v .
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

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/index"
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

// Close all of the open indices
func (neo *NeoSearch) Close() {
	for idx := range neo.Indices {
		index := neo.Indices[idx]
		index.Close()
	}

	neo.Indices = make([]*index.Index, 0)
}
