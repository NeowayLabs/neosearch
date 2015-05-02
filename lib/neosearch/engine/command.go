package engine

import (
	"fmt"
	"strconv"

	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

// Command defines a NeoSearch internal command.
// This command describes a single operation in the index storage and is
// decomposed in the following parts:
//   - Index
//   - Database
//   - Key
//   - KeyType
//   - Value
//   - ValueType
//   - Batch
type Command struct {
	Index     string
	Database  string
	Command   string
	Key       []byte
	KeyType   uint8
	Value     []byte
	ValueType uint8

	Batch bool
}

func (c Command) Println() {
	line := c.Reverse()
	fmt.Println(line)
}

func (c Command) Reverse() string {
	var (
		keyStr string
		valStr string
		line   string
	)

	if c.Key != nil {
		if c.KeyType == TypeString {
			keyStr = `'` + string(c.Key) + `'`
		} else if c.KeyType == TypeUint {
			keyStr = `uint(` + strconv.Itoa(int(utils.BytesToUint64(c.Key))) + `)`
		} else if c.KeyType == TypeInt {
			keyStr = `int(` + strconv.Itoa(int(utils.BytesToInt64(c.Key))) + `)`
		} else {
			panic(fmt.Errorf("Invalid command value type: %d", c.ValueType))
		}
	}

	if c.Value != nil {
		if c.ValueType == TypeString {
			valStr = `'` + string(c.Value) + `'`
		} else if c.ValueType == TypeUint {
			valStr = `uint(` + strconv.Itoa(int(utils.BytesToUint64(c.Value))) + `)`
		} else if c.ValueType == TypeInt {
			valStr = `int(` + strconv.Itoa(int(utils.BytesToInt64(c.Value))) + `)`
		} else {
			panic(fmt.Errorf("Invalid command key type: %d", c.KeyType))
		}
	}

	switch c.Command {
	case "set", "mergeset":
		line = fmt.Sprintf("USING %s.%s %s %s %s;", c.Index, c.Database, c.Command, keyStr, valStr)
	case "batch", "flushbatch":
		line = fmt.Sprintf("USING %s.%s %s;", c.Index, c.Database, c.Command)
	case "get", "delete":
		line = fmt.Sprintf("USING %s.%s %s %s;", c.Index, c.Database, c.Command, keyStr)
	default:
		panic(fmt.Errorf("Invalid command: %s: %v", c.Command, c))
	}

	return line
}
