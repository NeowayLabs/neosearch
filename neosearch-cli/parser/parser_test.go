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
	error := FromReader(strings.NewReader("using sample.TEST mergeset a 1;"), &commands)

	if error != nil {
		t.Error(error)
	}

	compareCommand(commands[0], engine.Command{
		Index:     "sample",
		Database:  "TEST",
		Command:   "mergeset",
		Key:       []byte("a"),
		KeyType:   engine.TypeString,
		Value:     utils.Uint64ToBytes(1),
		ValueType: engine.TypeUint,
	}, t)

	compareArray(`using sample.test.idx mergeset a 2;
             using sample.document.db set 1 "{id: 1, name: \"teste\"}";
             using sample.lalala set hello "world";
             using sample.mimimi get hello;
             using sample.lelele delete "teste";
             using sample.bleh.idx get uint(1);
             using sample.aaaa.bbb set uint(10000) int(10);
             using sample.bbbb.ccc mergeset "hellooooooooooooooooo" uint(102999299112211223);
             using sample.aaa delete "bbb"
        `, []engine.Command{
		engine.Command{
			Index:     "sample",
			Database:  "test.idx",
			Command:   "mergeset",
			Key:       []byte("a"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		engine.Command{
			Index:     "sample",
			Database:  "document.db",
			Command:   "set",
			Key:       utils.Int64ToBytes(1),
			KeyType:   engine.TypeInt,
			Value:     []byte("{id: 1, name: \"teste\"}"),
			ValueType: engine.TypeString,
		},
		engine.Command{
			Index:     "sample",
			Database:  "lalala",
			Command:   "set",
			Key:       []byte("hello"),
			KeyType:   engine.TypeString,
			Value:     []byte("world"),
			ValueType: engine.TypeString,
		},
		engine.Command{
			Index:    "sample",
			Database: "mimimi",
			Command:  "get",
			Key:      []byte("hello"),
			KeyType:  engine.TypeString,
		},
		engine.Command{
			Index:    "sample",
			Database: "lelele",
			Command:  "delete",
			Key:      []byte("teste"),
			KeyType:  engine.TypeString,
		},
		engine.Command{
			Index:    "sample",
			Database: "bleh.idx",
			Command:  "get",
			Key:      utils.Uint64ToBytes(1),
			KeyType:  engine.TypeUint,
		},
		engine.Command{
			Index:     "sample",
			Database:  "aaaa.bbb",
			Command:   "set",
			Key:       utils.Uint64ToBytes(10000),
			KeyType:   engine.TypeUint,
			Value:     utils.Int64ToBytes(10),
			ValueType: engine.TypeInt,
		},
		engine.Command{
			Index:     "sample",
			Database:  "bbbb.ccc",
			Command:   "mergeset",
			Key:       []byte("hellooooooooooooooooo"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(102999299112211223),
			ValueType: engine.TypeUint,
		},
		engine.Command{
			Index:    "sample",
			Database: "aaa",
			Command:  "delete",
			Key:      []byte("bbb"),
			KeyType:  engine.TypeString,
		},
	}, t)

	// underscore in the index name should pass
	commands = []engine.Command{}
	error = FromReader(strings.NewReader(`using sample.user_password set admin "s3cr3t"`),
		&commands)

	if error != nil {
		t.Error(error)
	}

	compareCommand(commands[0], engine.Command{
		Index:     "sample",
		Database:  "user_password",
		Command:   "set",
		Key:       []byte("admin"),
		KeyType:   engine.TypeString,
		Value:     []byte("s3cr3t"),
		ValueType: engine.TypeString,
	}, t)

	// invalid keyword "usinga"
	shouldThrowError(`
             usinga sample.test.idx set "hello" "world";
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
