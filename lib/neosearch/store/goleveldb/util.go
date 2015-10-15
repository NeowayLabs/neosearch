package goleveldb

import (
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func defaultWriteOptions() *opt.WriteOptions {
	wo := &opt.WriteOptions{}
	// request fsync on write for safety
	wo.Sync = true
	return wo
}

func defaultReadOptions() *opt.ReadOptions {
	ro := &opt.ReadOptions{}
	return ro
}
