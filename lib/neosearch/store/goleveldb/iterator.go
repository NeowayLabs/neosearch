package goleveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type LVDBIterator struct {
	iterator iterator.Iterator
}

func newIterator(store *LVDB) *LVDBIterator {
	options := defaultReadOptions()
	iter := store.db.NewIterator(nil, options)
	return &LVDBIterator{
		iterator: iter,
	}
}

func newIteratorWithSnapshot(store *LVDB, snapshot *leveldb.Snapshot) *LVDBIterator {
	options := defaultReadOptions()
	iter := snapshot.NewIterator(nil, options)
	return &LVDBIterator{
		iterator: iter,
	}
}

func (ldi *LVDBIterator) Next() {
	ldi.iterator.Next()
}

func (ldi *LVDBIterator) Prev() {
	ldi.iterator.Prev()
}

func (ldi *LVDBIterator) Valid() bool {
	return ldi.iterator.Valid()
}

func (ldi *LVDBIterator) SeekToFirst() {
	ldi.iterator.First()
}

func (ldi *LVDBIterator) SeekToLast() {
	ldi.iterator.Last()
}

func (ldi *LVDBIterator) Seek(key []byte) {
	ldi.iterator.Seek(key)
}

func (ldi *LVDBIterator) Key() []byte {
	return ldi.iterator.Key()
}

func (ldi *LVDBIterator) Value() []byte {
	return ldi.iterator.Value()
}

func (ldi *LVDBIterator) GetError() error {
	return ldi.iterator.Error()
}

func (ldi *LVDBIterator) Close() error {
	ldi.iterator.Release()
	return nil
}
