package test

import (
	"reflect"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

func CommonTestStoreSetGet(t *testing.T, kv store.KVStore) {
	var (
		err  error
		data []byte
	)

	type kvTest struct {
		key   []byte
		value []byte
	}

	shouldPass := []kvTest{
		{
			key:   []byte{'t', 'e', 's', 't', 'e'},
			value: []byte{'i', '4', 'k'},
		},
		{
			key: []byte{'p', 'l', 'a', 'n', '9'},
			value: []byte{'f', 'r', 'o', 'm',
				'o', 'u', 't', 'e', 'r', 's',
				's', 'p', 'a', 'c', 'e', '!'},
		},
		{
			key:   []byte{'t', 'h', 'e', 'm', 'a', 't', 'r', 'i', 'x'},
			value: []byte{'h', 'a', 's', 'y', 'o', 'u'},
		},
	}

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	for _, kv := range shouldPass {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for _, kv := range shouldPass {
		if data, err = reader.Get(kv.key); err != nil {
			t.Error(err)
			continue
		} else if data == nil || len(data) != len(kv.value) {
			t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(kv.key), string(data))
			continue
		}

		if !reflect.DeepEqual(data, kv.value) {
			t.Errorf("Data retrieved '%s' != '%s'", string(data), string(kv.value))
		}
	}

	data, err = reader.Get([]byte("do not exists"))
	if err != nil {
		t.Error(err)
	}

	// key does not exists, data should be nil
	if data != nil {
		t.Error("key 'does not exists' returning wrong value")
	}
}

func CommonTestBatchWrite(t *testing.T, kv store.KVStore) {
	var (
		err   error
		key   = []byte{'a'}
		value = []byte{'b'}
		data  []byte
	)

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	writer.StartBatch()

	if writer.IsBatch() == false {
		t.Error("StartBatch not setting isBatch = true")
		return
	}

	if err = writer.Set(key, value); err != nil {
		t.Error(err)
		return
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// should returns nil, nil because the key is in the batch cache
	if data, err = reader.Get(key); err != nil || data != nil {
		t.Error("Key set before wasn't in the write batch cache." +
			" Batch mode isnt working")
	}

	if err = writer.FlushBatch(); err != nil {
		t.Error(err)
	}

	if writer.IsBatch() == true {
		t.Error("FlushBatch doesnt reset the isBatch")
	}

	newReader := kv.Reader()
	if newReader == nil {
		t.Error("newReader not created!")
		return
	}
	defer func() {
		err := newReader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if data, err = newReader.Get(key); err != nil {
		t.Error(err)
	} else if data == nil || len(data) != len(value) {
		t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(key), string(data))
	}

	if !reflect.DeepEqual(data, value) {
		t.Errorf("Data retrieved '%s' != '%s'", string(data), string(value))
	}
}

func CommonTestBatchMultiWrite(t *testing.T, kv store.KVStore) {
	var (
		err  error
		data []byte
	)

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	writer.StartBatch()

	type kvTest struct {
		key   []byte
		value []byte
	}

	shouldPass := []kvTest{
		{
			key:   []byte{'t', 'e', 's', 't', 'e'},
			value: []byte{'i', '4', 'k'},
		},
		{
			key: []byte{'p', 'l', 'a', 'n', '9'},
			value: []byte{'f', 'r', 'o', 'm',
				'o', 'u', 't', 'e', 'r', 's',
				's', 'p', 'a', 'c', 'e', '!'},
		},
		{
			key:   []byte{'t', 'h', 'e', 'm', 'a', 't', 'r', 'i', 'x'},
			value: []byte{'h', 'a', 's', 'y', 'o', 'u'},
		},
	}

	for _, kv := range shouldPass {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for _, kv := range shouldPass {
		if data, err := reader.Get(kv.key); err != nil || data != nil {
			t.Error("Key set before wasn't in the write batch cache." +
				" Batch mode isnt working")
		}
	}

	if err := writer.FlushBatch(); err != nil {
		t.Error(err)
	}

	newReader := kv.Reader()
	if newReader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := newReader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for _, kv := range shouldPass {
		if data, err = newReader.Get(kv.key); err != nil {
			t.Error(err)
			continue
		} else if data == nil || len(data) != len(kv.value) {
			t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(kv.key), string(data))
			continue
		}

		if !reflect.DeepEqual(data, kv.value) {
			t.Errorf("Data retrieved '%s' != '%s'", string(data), string(kv.value))
		}
	}
}

func CommonTestStoreMergeSet(t *testing.T, kv store.KVStore) {
	var (
		err  error
		data []byte
	)

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	key := []byte{'t', 'e', 's', 't', 'e'}
	values := []uint64{0, 2, 1}

	result := append(utils.Uint64ToBytes(values[0]), utils.Uint64ToBytes(values[2])...)
	result = append(result, utils.Uint64ToBytes(values[1])...)

	for _, value := range values {
		if err = writer.MergeSet(key, value); err != nil {
			t.Error(err)
			return
		}
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if data, err = reader.Get(key); err != nil {
		t.Error(err)
	} else if data == nil || len(data) != len(result) {
		t.Errorf("Failed to retrieve key '%s'. Retuns: %+v", string(key), data)
	}

	if !reflect.DeepEqual(data, result) {
		t.Errorf("Data retrieved '%+v' != '%+v'", data, result)
	}
}

func CommonTestStoreIterator(t *testing.T, kv store.KVStore) {
	var err error

	writer := kv.Writer()

	type kvTest struct {
		key   []byte
		value []byte
	}

	data := []kvTest{
		{
			key: []byte{'p', 'l', 'a', 'n', '9'},
			value: []byte{'f', 'r', 'o', 'm',
				'o', 'u', 't', 'e', 'r', 's',
				's', 'p', 'a', 'c', 'e', '!'},
		},
		{
			key:   []byte{'t', 'e', 's', 't', 'e'},
			value: []byte{'i', '4', 'k'},
		},
		{
			key:   []byte{'t', 'h', 'e', 'm', 'a', 't', 'r', 'i', 'x'},
			value: []byte{'h', 'a', 's', 'y', 'o', 'u'},
		},
	}

	writer.StartBatch()
	for _, kv := range data {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}
	}

	err = writer.FlushBatch()
	if err != nil {
		t.Fatal(err)
	}

	reader := kv.Reader()
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	it := reader.GetIterator()
	it.Seek([]byte{0})
	keys := make([][]byte, 0, len(data))
	for it.Valid() {
		k := it.Key()
		key := make([]byte, len(k))
		copy(key, k)
		keys = append(keys, key)
		it.Next()
	}

	if len(keys) != len(data) {
		t.Errorf("expected same number of keys, got %d != %d", len(keys), len(data))
	}
	for i, dk := range data {
		if !reflect.DeepEqual(dk.key, keys[i]) {
			t.Errorf("expected key %s got %s", dk.key, keys[i])
		}
	}

	err = it.Close()
	if err != nil {
		t.Fatal(err)
	}
}
