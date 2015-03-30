package parser

import (
	"strings"
	"testing"

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/utils"
)

func TestFromReader(t *testing.T) {
	commands := []engine.Command{}
	error := FromReader(strings.NewReader("using test mergeset a 1;"), &commands)

	if error != nil {
		t.Error(error)
	}

	compareCommand(commands[0], engine.Command{
		Index:   "test",
		Command: "mergeset",
		Key:     []byte("a"),
		Value:   utils.Uint64ToBytes(1),
	}, t)

	compareArray(`using test.idx mergeset a 2;
             using document.db mergeset 1 "{id: 1, name: \"teste\"}";
             using lalala set hello "world";
             using mimimi get hello;
             using lelele delete "teste";
             using aaa delete "bbb"
        `, []engine.Command{
		engine.Command{
			Index:   "test.idx",
			Command: "mergeset",
			Key:     []byte("a"),
			Value:   utils.Uint64ToBytes(2),
		},
		engine.Command{
			Index:   "document.db",
			Command: "mergeset",
			Key:     utils.Int64ToBytes(1),
			Value:   []byte("{id: 1, name: \"teste\"}"),
		},
		engine.Command{
			Index:   "lalala",
			Command: "set",
			Key:     []byte("hello"),
			Value:   []byte("world"),
		},
		engine.Command{
			Index:   "mimimi",
			Command: "get",
			Key:     []byte("hello"),
		},
		engine.Command{
			Index:   "lelele",
			Command: "delete",
			Key:     []byte("teste"),
		},
		engine.Command{
			Index:   "aaa",
			Command: "delete",
			Key:     []byte("bbb"),
		},
	}, t)

	// underscore in the index name should pass
	commands = []engine.Command{}
	error = FromReader(strings.NewReader(`using user_password set admin "s3cr3t"`),
		&commands)

	if error != nil {
		t.Error(error)
	}

	compareCommand(commands[0], engine.Command{
		Index:   "user_password",
		Command: "set",
		Key:     []byte("admin"),
		Value:   []byte("s3cr3t"),
	}, t)

	// invalid keyword "usinga"
	shouldThrowError(`
             usinga test.idx set "hello" "world";
        `, t)

}

func compareCommand(cmd engine.Command, expected engine.Command, t *testing.T) {
	if cmd.Index != expected.Index ||
		cmd.Command != expected.Command ||
		string(cmd.Key) != string(expected.Key) ||
		string(cmd.Value) != string(expected.Value) {
		t.Errorf("Unexpected parsed command: %v !== %v", cmd, expected)
	}
}

func shouldThrowError(bufferCommands string, t *testing.T) {
	resultCommands := []engine.Command{}

	error := FromReader(strings.NewReader(bufferCommands), &resultCommands)

	if error == nil {
		t.Errorf("Test SHOULD fail: %v", resultCommands)
		return
	}
}

func compareArray(bufferCommands string, expectedCommands []engine.Command, t *testing.T) {
	resultCommands := []engine.Command{}

	error := FromReader(strings.NewReader(bufferCommands), &resultCommands)

	if error != nil {
		t.Error(error)
		return
	}

	if len(resultCommands) != len(expectedCommands) {
		t.Errorf("Failed to parse all of the cmdline tests:\n\t %v !== \n\t %v", resultCommands, expectedCommands)
	}

	for i := 0; i < len(resultCommands); i++ {
		compareCommand(resultCommands[i], expectedCommands[i], t)
	}
}
