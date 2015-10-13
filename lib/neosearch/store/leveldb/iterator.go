// +build leveldb

package leveldb

import "github.com/jmhodges/levigo"

// Simple wrapper around levigo.Iterator to proper implement the KVIterator interface
type LVDBIterator struct {
	*levigo.Iterator
}

// Close the iterator. It's only a wrapper for levigo.Iterator, that does not returns error in Close
// method
func (i LVDBIterator) Close() error {
	i.Iterator.Close()
	return nil
}
