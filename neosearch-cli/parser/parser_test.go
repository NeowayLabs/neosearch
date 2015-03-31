package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/utils"
)

func TestCliParserFromReader(t *testing.T) {
	commands := []engine.Command{}
	error := FromReader(strings.NewReader("using test mergeset a 1;"), &commands)

	if error != nil {
		t.Error(error)
	}

	compareCommand(commands[0], engine.Command{
		Index:     "test",
		Command:   "mergeset",
		Key:       []byte("a"),
		KeyType:   engine.TypeString,
		Value:     utils.Uint64ToBytes(1),
		ValueType: engine.TypeUint,
	}, t)

	compareArray(`using test.idx mergeset a 2;
             using document.db set 1 "{id: 1, name: \"teste\"}";
             using lalala set hello "world";
             using mimimi get hello;
             using lelele delete "teste";
             using bleh.idx get uint(1);
             using aaaa.bbb set uint(10000) int(10);
             using bbbb.ccc mergeset "hellooooooooooooooooo" uint(102999299112211223);
             using aaa delete "bbb"
        `, []engine.Command{
		engine.Command{
			Index:     "test.idx",
			Command:   "mergeset",
			Key:       []byte("a"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		engine.Command{
			Index:     "document.db",
			Command:   "set",
			Key:       utils.Int64ToBytes(1),
			KeyType:   engine.TypeInt,
			Value:     []byte("{id: 1, name: \"teste\"}"),
			ValueType: engine.TypeString,
		},
		engine.Command{
			Index:     "lalala",
			Command:   "set",
			Key:       []byte("hello"),
			KeyType:   engine.TypeString,
			Value:     []byte("world"),
			ValueType: engine.TypeString,
		},
		engine.Command{
			Index:   "mimimi",
			Command: "get",
			Key:     []byte("hello"),
			KeyType: engine.TypeString,
		},
		engine.Command{
			Index:   "lelele",
			Command: "delete",
			Key:     []byte("teste"),
			KeyType: engine.TypeString,
		},
		engine.Command{
			Index:   "bleh.idx",
			Command: "get",
			Key:     utils.Uint64ToBytes(1),
			KeyType: engine.TypeUint,
		},
		engine.Command{
			Index:     "aaaa.bbb",
			Command:   "set",
			Key:       utils.Uint64ToBytes(10000),
			KeyType:   engine.TypeUint,
			Value:     utils.Int64ToBytes(10),
			ValueType: engine.TypeInt,
		},
		engine.Command{
			Index:     "bbbb.ccc",
			Command:   "mergeset",
			Key:       []byte("hellooooooooooooooooo"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(102999299112211223),
			ValueType: engine.TypeUint,
		},
		engine.Command{
			Index:   "aaa",
			Command: "delete",
			Key:     []byte("bbb"),
			KeyType: engine.TypeString,
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
		Index:     "user_password",
		Command:   "set",
		Key:       []byte("admin"),
		KeyType:   engine.TypeString,
		Value:     []byte("s3cr3t"),
		ValueType: engine.TypeString,
	}, t)

	// invalid keyword "usinga"
	shouldThrowError(`
             usinga test.idx set "hello" "world";
        `, t)

}

func compareCommand(cmd engine.Command, expected engine.Command, t *testing.T) {
	if !reflect.DeepEqual(cmd, expected) {
		t.Errorf("Unexpected parsed command: %v !== %v", cmd.Reverse(), expected.Reverse())
		fmt.Printf("%v !== %v\n", cmd, expected)
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
