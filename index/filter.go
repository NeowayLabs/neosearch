package index

import (
	"bytes"

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/utils"
)

func (i *Index) filterTerm(field, value []byte) ([]uint64, error) {
	var (
		docIDs []uint64
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
		for i := 0; i < len(data); i += 8 {
			v := utils.BytesToUint64(data[i : i+8])
			docIDs = append(docIDs, v)
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
		docIDs []uint64
	)

	storekv, err := i.engine.GetStore(string(field) + ".idx")

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
