package goleveldb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type LVDBWriter struct {
	store   *LVDB
	mutex   sync.Mutex
	isBatch bool
}

// newWriter returns a new writer
func newWriter(lvdb *LVDB) *LVDBWriter {
	return &LVDBWriter{
		store: lvdb,
	}
}

// Set put or update the key with the given value
func (w *LVDBWriter) Set(key, value []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.isBatch {
		// isBatch == true, we can safely access _writeBatch pointer
		w.store.writeBatch.Put(key, value)
		return nil
	}

	options := defaultWriteOptions()
	return w.store.db.Put(key, value, options)
}

// MergeSet add value to a ordered set of integers stored in key. If value
// is already on the key, than the set will be skipped.
func (w *LVDBWriter) MergeSet(key []byte, value uint64) error {
	var (
		buf      *bytes.Buffer
		err      error
		v        uint64
		i        uint64
		inserted bool
		reader   = w.store.Reader()
	)

	data, err := reader.Get(key)
	if err != nil {
		return err
	}

	if w.store.debug {
		fmt.Printf("[INFO] %d ids == %d GB of ids\n", len(data)/8, len(data)/(1024*1024*1024))
	}

	buf = new(bytes.Buffer)
	lenBytes := uint64(len(data))

	// O(n)
	for i = 0; i < lenBytes; i += 8 {
		v = utils.BytesToUint64(data[i : i+8])

		// returns if value is already stored
		if v == value {
			return nil
		}

		if value < v {
			err = binary.Write(buf, binary.BigEndian, value)
			if err != nil {
				return err
			}

			inserted = true
		}

		err = binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}

	if lenBytes == 0 || !inserted {
		err = binary.Write(buf, binary.BigEndian, value)
		if err != nil {
			return err
		}
	}

	return w.Set(key, buf.Bytes())
}

func (w *LVDBWriter) Delete(key []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.isBatch {
		w.store.writeBatch.Delete(key)
		return nil
	}

	options := defaultWriteOptions()
	return w.store.db.Delete(key, options)
}

// StartBatch start a new batch write processing
func (w *LVDBWriter) StartBatch() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.store.writeBatch == nil {
		w.store.writeBatch = new(leveldb.Batch)
	} else {
		w.store.writeBatch.Reset()
	}

	w.isBatch = true
}

// IsBatch returns true if LVDB is in batch mode
func (w *LVDBWriter) IsBatch() bool {
	return w.isBatch
}

// FlushBatch writes the batch to disk
func (w *LVDBWriter) FlushBatch() error {
	var err error

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.store.writeBatch != nil {
		options := defaultWriteOptions()
		err = w.store.db.Write(w.store.writeBatch, options)
		// After flush, release the writeBatch for future uses
		w.store.writeBatch.Reset()
		w.isBatch = false
	}

	return err
}
