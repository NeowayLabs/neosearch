package index

import "github.com/NeowayLabs/neosearch/version"

type FieldInfo struct {
	Type string `json:"type,omitempty"`
	Size uint64 `json:"size,omitempty"`
}

type IndexInfo struct {
	Version string
	Fields  map[string]FieldInfo `json:"fields"`
	Size    uint64               `json:"size"`
}

func NewIndexInfo() *IndexInfo {
	indFields := &IndexInfo{
		Version: version.Version,
		Size:    0,
		Fields:  make(map[string]FieldInfo),
	}

	return indFields
}
