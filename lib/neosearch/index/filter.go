package index

import (
	"bytes"

	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

func (i *Index) FilterTermID(field, value []byte, limit uint64) ([]uint64, uint64, error) {
	cmd := engine.Command{}
	cmd.Index = i.Name
	cmd.Database = string(field) + ".idx"
	cmd.Command = "get"
	cmd.Key = value
	data, err := i.engine.Execute(cmd)

	if err != nil {
		return nil, 0, err
	}

	dataLimit := uint64(len(data) / 8)
	total := dataLimit

	if limit > 0 && limit < dataLimit {
		dataLimit = limit
	}

	docIDs := make([]uint64, dataLimit)

	if len(data) > 0 {
		for i, j := uint64(0), uint64(0); i < dataLimit*8; i, j = i+8, j+1 {
			v := utils.BytesToUint64(data[i : i+8])
			docIDs[j] = v
		}

	}

	return docIDs, total, nil
}

// FilterTerm filter the index for all documents that have `value` in the
// field `field` and returns upto `limit` documents. A limit of 0 (zero) is
// the same as no limit (all of the records will return)..
func (i *Index) FilterTerm(field []byte, value []byte, limit uint64) ([]string, uint64, error) {
	docIDs, total, err := i.FilterTermID(field, value, limit)

	if err != nil {
		return nil, 0, err
	}

	docs := make([]string, len(docIDs))

	for idx, docID := range docIDs {
		if byteDoc, err := i.Get(docID); err == nil {
			docs[idx] = string(byteDoc)
		} else {
			return nil, 0, err
		}
	}

	return docs, total, nil
}

func (i *Index) matchPrefix(field []byte, value []byte) ([]uint64, error) {
	var (
		docIDs []uint64
	)

	storekv, err := i.engine.GetStore(i.Name, string(field)+".idx")

	if err != nil {
		return nil, err
	}

	it := storekv.GetIterator()

	defer it.Close()

	for it.Seek(value); it.Valid(); it.Next() {
		if bytes.HasPrefix(it.Key(), value) {
			var ids []uint64
			dataBytes := it.Value()

			if len(dataBytes) == 0 {
				continue
			}

			for i := 0; i < len(dataBytes); i += 8 {
				v := utils.BytesToUint64(dataBytes[i : i+8])
				ids = append(ids, v)
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
