package main

import (
	"fmt"
	"os"

	"github.com/NeowayLabs/neosearch/cmd/cli/parser"
	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
)

func batch(ng *engine.Engine, filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	commands := []engine.Command{}

	err = parser.FromReader(file, &commands)

	for _, cmd := range commands {
		d, err := ng.Execute(cmd)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Data: ", d)
		}
	}

	return nil
}
