package engine

import (
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

func TestCommand(t *testing.T) {
	for _, testTable := range []struct {
		cmd      Command
		expected string
	}{
		{
			cmd: Command{
				Database:  "name.idx",
				Index:     "empresas",
				Command:   "mergeset",
				Key:       []byte("teste"),
				KeyType:   TypeString,
				Value:     utils.Uint64ToBytes(1000),
				ValueType: TypeUint,
			},
			expected: `USING empresas.name.idx MERGESET 'teste' uint(1000);`,
		},
		{
			cmd: Command{
				Database:  "name.idx",
				Index:     "empresas",
				Command:   "batch",
				Key:       nil,
				KeyType:   TypeNil,
				Value:     nil,
				ValueType: TypeNil,
			},
			expected: `USING empresas.name.idx BATCH;`,
		},
		{
			cmd: Command{
				Database:  "name.idx",
				Index:     "empresas",
				Command:   "get",
				Key:       []byte("teste"),
				KeyType:   TypeString,
				Value:     nil,
				ValueType: TypeNil,
			},
			expected: `USING empresas.name.idx GET 'teste';`,
		},
	} {
		cmdRev := testTable.cmd.Reverse()

		if cmdRev != testTable.expected {
			t.Errorf("Differs: '%s' !== '%s'", cmdRev, testTable.expected)
		}
	}
}
