package index

type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size uint64 `json:"size"`
}

type IndexInfo struct {
	Fields []FieldInfo `json:"fields"`
	Size   uint64      `json:"size"`
}
