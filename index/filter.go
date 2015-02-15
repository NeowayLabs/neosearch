package index

import "bitbucket.org/i4k/neosearch/engine"

func (i *Index) FilterTerm(field []byte, value []byte) ([]byte, error) {
	cmd := engine.Command{}
	cmd.Index = string(field) + ".idx"
	cmd.Command = "get"
	cmd.Key = value
	ret, err := i.engine.Execute(cmd)
	return ret, err
}
