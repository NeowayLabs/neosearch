package index

import (
	"encoding/json"

	"bytes"

	"bitbucket.org/i4k/neosearch/engine"
	"bitbucket.org/i4k/neosearch/utils"
)

func (i *Index) FilterTerm(field []byte, value []byte) ([]byte, error) {
	cmd := engine.Command{}
	cmd.Index = string(field) + ".idx"
	cmd.Command = "get"
	cmd.Key = value
	ret, err := i.engine.Execute(cmd)
	return ret, err
}

func (i *Index) MatchPrefix(field []byte, value []byte) ([]string, error) {
	var (
		docIDs []uint64
		docs   []string
	)

	store, err := i.engine.GetStore(string(field) + ".idx")

	if err != nil {
		return nil, err
	}

	it := (*store).GetIterator()

	defer it.Close()

	for it.Seek(value); it.Valid(); it.Next() {
		if bytes.HasPrefix(it.Key(), value) {
			var ids []uint64

			err := json.Unmarshal(it.Value(), &ids)

			if err != nil {
				return nil, err
			}

			for _, id := range ids {
				docIDs = utils.UniqueAdd(docIDs, id)
			}
		}
	}

	if err := it.GetError(); err != nil {
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
