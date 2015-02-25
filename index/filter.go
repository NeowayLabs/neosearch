package index

import (
	"encoding/gob"

	"bytes"

	"github.com/neowaylabs/neosearch/engine"
	"github.com/neowaylabs/neosearch/utils"
)

func (i *Index) filterTerm(field, value []byte) ([]uint64, error) {
	var (
		valueBytes bytes.Buffer
		decoder    = gob.NewDecoder(&valueBytes)
		docIDs     []uint64
	)

	cmd := engine.Command{}
	cmd.Index = string(field) + ".idx"
	cmd.Command = "get"
	cmd.Key = value
	data, err := i.engine.Execute(cmd)

	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		valueBytes.Write(data)
		err = decoder.Decode(&docIDs)

		if err != nil {
			return nil, err
		}
	}

	return docIDs, nil
}

// FilterTerm filters index for all documents that have `value` in the
// field `field`.
func (i *Index) FilterTerm(field []byte, value []byte) ([]string, error) {
	var docs []string

	docIDs, err := i.filterTerm(field, value)

	if err != nil {
		return nil, err
	}

	for _, docID := range docIDs {
		if byteDoc, err := i.Get(docID); err == nil {
			docs = append(docs, string(byteDoc))
		} else {
			return nil, err
		}
	}

	return docs, nil
}

func (i *Index) matchPrefix(field []byte, value []byte) ([]uint64, error) {
	var (
		docIDs     []uint64
		valueBytes bytes.Buffer
	)

	store, err := i.engine.GetStore(string(field) + ".idx")

	if err != nil {
		return nil, err
	}

	it := store.GetIterator()

	defer it.Close()

	for it.Seek(value); it.Valid(); it.Next() {
		if bytes.HasPrefix(it.Key(), value) {
			var ids []uint64
			dataBytes := it.Value()

			if len(dataBytes) == 0 {
				continue
			}

			valueBytes.Write(dataBytes)

			// We have some problem with encoding/gob here.
			// NewDecoder is a expansive call, but I'm having the
			// error "extra data in buffer". Aparently, the decoder
			// isn't cleaning something internal after successful
			// decoding the integer array
			decoder := gob.NewDecoder(&valueBytes)

			// ids will be a quick sorted array
			err := decoder.Decode(&ids)

			if err != nil {
				return nil, err
			}

			if len(docIDs) == 0 {
				docIDs = ids
				continue
			}

			for _, id := range ids {
				docIDs = utils.UniqueUint64Add(docIDs, id)
			}
		}
	}

	if err := it.GetError(); err != nil {
		return nil, err
	}

	return docIDs, nil
}

// MatchPrefix search documents where field `field` starts with `value`.
func (i *Index) MatchPrefix(field []byte, value []byte) ([]string, error) {
	var docs []string

	docIDs, err := i.matchPrefix(field, value)

	if err != nil {
		return nil, err
	}

	for _, docID := range docIDs {
		d, err := i.Get(docID)

		if err != nil {
			return nil, err
		}

		docs = append(docs, string(d))
	}

	return docs, nil
}
