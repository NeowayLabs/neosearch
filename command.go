package neosearch

// Command defines a NeoSearch command
type Command struct {
	Index   string
	Command string
	Key     []byte
	Value   []byte
}
