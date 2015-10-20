// +build leveldb

package leveldb

import (
	"sync"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/jmhodges/levigo"
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
func (w *LVDBWriter) Set(key []byte, value []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.isBatch {
		// isBatch == true, we can safely access _writeBatch pointer
		w.store.writeBatch.Put(key, value)
		return nil
	}

	options := defaultWriteOptions()
	err := w.store.db.Put(options, key, value)
	options.Close()
	return err
}

// Get returns the value of the given key
func (w *LVDBWriter) Get(key []byte) ([]byte, error) {
	options := defaultReadOptions()
	b, e := w.store.db.Get(options, key)
	options.Close()
	return b, e
}

// MergeSet add value to a ordered set of integers stored in key. If value
// is already on the key, than the set will be skipped.
func (w *LVDBWriter) MergeSet(key []byte, value uint64) error {
	return store.MergeSet(w, key, value, w.store.debug)
}

// Delete remove the given key
func (w *LVDBWriter) Delete(key []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.isBatch {
		w.store.writeBatch.Delete(key)
		return nil
	}

	options := defaultWriteOptions()
	err := w.store.db.Delete(options, key)
	options.Close()
	return err
}

// StartBatch start a new batch write processing
func (w *LVDBWriter) StartBatch() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.store.writeBatch == nil {
		w.store.writeBatch = levigo.NewWriteBatch()
	} else {
		w.store.writeBatch.Clear()
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
		err = w.store.db.Write(options, w.store.writeBatch)
		options.Close()
		// After flush, release the writeBatch for future uses
		w.store.writeBatch.Clear()
		w.isBatch = false
	}

	return err
}
