package index

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

// TODO; test buildIndex* functions
//       use table-driven tests

func TestIndexBuildUintCommands(t *testing.T) {
	var (
		idx      *Index
		err      error
		indexDir = DataDirTmp + "/test-uint"
		commands []engine.Command
		cmd      engine.Command
	)

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

	err = os.MkdirAll(DataDirTmp, 0755)

	if err != nil {
		t.Error("Failed to create directory")
		goto cleanup
	}

	idx, err = New("test-uint", cfg, true)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	commands, err = idx.buildIndexUint64(uint64(1), "teste", uint64(1))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if len(commands) != 1 {
		t.Error("invalid commands returned by buildUint")
		goto cleanup
	}

	cmd = commands[0]

	if cmd.Index != "test-uint" ||
		cmd.Database != "teste_uint.idx" ||
		utils.BytesToUint64(cmd.Key) != uint64(1) ||
		utils.BytesToUint64(cmd.Value) != uint64(1) ||
		strings.ToLower(cmd.Command) != "mergeset" {
		t.Error("commands differs")
		fmt.Println("Key: ", utils.BytesToUint64(cmd.Key))
		fmt.Println("Value: ", utils.BytesToUint64(cmd.Value))
		fmt.Println("Index: ", cmd.Index)
		cmd.Println()
		goto cleanup
	}

cleanup:
	idx.Close()
	os.RemoveAll(indexDir)
}

func TestIndexBuildBoolCommands(t *testing.T) {
	var (
		idx      *Index
		err      error
		indexDir = DataDirTmp + "/test-bool"
		commands []engine.Command
		cmd      engine.Command
	)

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

	err = os.MkdirAll(DataDirTmp, 0755)

	if err != nil {
		t.Error("Failed to create directory")
		goto cleanup
	}

	idx, err = New("test-bool", cfg, true)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	commands, err = idx.buildIndexBool(uint64(1), "teste", true)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if len(commands) != 1 {
		t.Error("invalid commands returned by buildUint")
		goto cleanup
	}

	cmd = commands[0]

	if cmd.Index != "test-bool" ||
		cmd.Database != "teste_bool.idx" ||
		utils.BytesToBool(cmd.Key) != true ||
		utils.BytesToUint64(cmd.Value) != uint64(1) ||
		strings.ToLower(cmd.Command) != "mergeset" {
		t.Error("commands differs")
		fmt.Println("Key: ", utils.BytesToUint64(cmd.Key))
		fmt.Println("Value: ", utils.BytesToUint64(cmd.Value))
		fmt.Println("Index: ", cmd.Index)
		cmd.Println()
		goto cleanup
	}

cleanup:
	idx.Close()
	os.RemoveAll(indexDir)
}
