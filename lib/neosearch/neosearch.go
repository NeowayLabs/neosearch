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
//     - snappy (optional, only required for compressed data)
//     - Go > 1.3
//
// Install
//
//     go get -u -v github.com/NeowayLabs/neosearch
//     cd $GOPATH/src/github/NeowayLabs/neosearch
//     go test -v .
//
// Create and add documents to NeoSearch is very easy, see below:
//
//   func main() {
//       cfg := NewConfig()
//       cfg.Option(neosearch.DataDir("/tmp"))
//       cfg.Option(neosearch.Debug(false))
//
//       neo := neosearch.New(cfg)
//
//       index, err := neosearch.CreateIndex("test")
//       if err != nil {
//           panic(err)
//       }
//
//       err = index.Add(1, `{"name": "Neoway Business Solution", "type": "company"}`)
//       if err != nil {
//           panic(err)
//       }
//
//       err = index.Add(2, `{"name": "Facebook Inc", "type": "company"}`)
//       if err != nil {
//           panic(err)
//       }
//
//       values, err := index.MatchPrefix([]byte("name"), []byte("neoway"))
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
//   - Create/Delete index
//   - Index JSON documents (No schema)
//   - Bulk writes
//   - Analysers
//     - Tokenizer
//   - Search
//     - MatchPrefix
//     - FilterTerm
//
// This project is in active development stage, it is not recommended for
// production environments.
package neosearch

import (
	"errors"
	"expvar"
	"fmt"
	"os"

	"github.com/NeowayLabs/neosearch/lib/neosearch/cache"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/NeowayLabs/neosearch/lib/neosearch/index"
)

var (
	cachedIndices = expvar.NewInt("cachedIndices")
)

// NeoSearch is the core of the neosearch package.
// This structure handles all of the user's interactions with the indices,
// like CreateIndex, DeleteIndex, UpdateIndex and others.
type NeoSearch struct {
	indices cache.Cache
	config  *config.Config
}

// New creates the NeoSearch high-level interface.
// Use that for index/update/delete JSON documents.
func New(cfg *config.Config) *NeoSearch {
	if cfg.DataDir == "" {
		panic(errors.New("DataDir is required for NeoSearch interface"))
	}

	if cfg.DataDir[len(cfg.DataDir)-1] == '/' {
		cfg.DataDir = cfg.DataDir[0 : len(cfg.DataDir)-1]
	}

	if cfg.MaxIndicesOpen == 0 {
		cfg.MaxIndicesOpen = config.DefaultMaxIndicesOpen
	}

	neo := &NeoSearch{
		config:  cfg,
		indices: cache.NewLRUCache(cfg.MaxIndicesOpen),
	}

	neo.indices.OnRemove(func(key string, value interface{}) {
		v, ok := value.(*index.Index)

		if ok {
			v.Close()
		}
	})

	return neo
}

// GetIndices returns cache indices
func (neo *NeoSearch) GetIndices() cache.Cache {
	return neo.indices
}

// CreateIndex creates and setup a new index
func (neo *NeoSearch) CreateIndex(name string) (*index.Index, error) {
	indx, err := index.New(name, neo.config, true)
	if err != nil {
		return nil, err
	}

	neo.indices.Add(name, indx)
	cachedIndices.Set(int64(neo.indices.Len()))
	return indx, nil
}

// DeleteIndex does exactly what the name says.
func (neo *NeoSearch) DeleteIndex(name string) error {
	// closes the index on remove
	neo.indices.Remove(name)
	idxLen := neo.indices.Len()
	cachedIndices.Set(int64(idxLen))

	if exists, err := neo.IndexExists(name); exists == true && err == nil {
		err := os.RemoveAll(neo.config.DataDir + "/" + name)
		return err
	}

	return errors.New("Index '" + name + "' not found.")
}

// OpenIndex open a existing index for read/write operations.
func (neo *NeoSearch) OpenIndex(name string) (*index.Index, error) {
	var (
		ok         bool
		cacheIndex interface{}
		indx       *index.Index
		err        error
	)

	cacheIndex, ok = neo.indices.Get(name)

	if ok && cacheIndex != nil {
		indx, ok = cacheIndex.(*index.Index)

		if !ok {
			panic("NeoSearch has inconsistent data stored in cache")
		}

		return indx, nil
	}

	ok, err = neo.IndexExists(name)

	if err == nil && !ok {
		return nil, fmt.Errorf("Index '%s' not found in directory '%s'.", name, neo.config.DataDir)
	} else if err != nil {
		return nil, err
	}

	indx, err = index.New(name, neo.config, false)
	if err != nil {
		return nil, err
	}

	neo.indices.Add(name, indx)
	cachedIndices.Set(int64(neo.indices.Len()))
	return indx, nil
}

// IndexExists verifies if the directory of the index given by name exists
func (neo *NeoSearch) IndexExists(name string) (bool, error) {
	indexPath := neo.config.DataDir + "/" + name
	_, err := os.Stat(indexPath)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Close all of the open indices
func (neo *NeoSearch) Close() {
	neo.indices.Clean()
	cachedIndices.Set(int64(neo.indices.Len()))
}
