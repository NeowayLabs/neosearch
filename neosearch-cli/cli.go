package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/peterh/liner"

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/neosearch-cli/parser"
	"github.com/NeowayLabs/neosearch/utils"
)

var (
	historyFile = "cli.history.txt"
	keywords    = []string{"using", "set", "get", "mergeset", "delete"}
)

func setupNeosearchDir(homePath string) error {
	return os.Mkdir(homePath+"/.neosearch/", 0755)
}

func cli(ng *engine.Engine, homePath string) error {
	var cmdline string
	var err error
	var enableHistory bool

	if homePath != "" {
		enableHistory = true
		setupNeosearchDir(homePath)
	} else {
		fmt.Printf("No user home provided... \n")
		fmt.Printf("Provide --home to enable command history.\n")
	}

	line := liner.NewLiner()
	defer line.Close()

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range keywords {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if enableHistory {
		if f, err := os.Open(homePath + "/.neosearch/" + historyFile); err == nil {
			line.ReadHistory(f)
			f.Close()
		}
	}

	// command-line here
	for {
		if cmdline, err = line.Prompt("neosearch>"); err != nil {
			if err.Error() == "EOF" {
				break
			}

			continue
		}

		line.AppendHistory(cmdline)

		command := []engine.Command{}

		if strings.ToLower(cmdline) == "quit" ||
			strings.ToLower(cmdline) == "quit;" {
			break
		}

		err = parser.FromString(cmdline, &command)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, cmd := range command {
				data, err := ng.Execute(cmd)
				if err != nil {
					fmt.Println("ERROR: ", err)
				} else {
					fmt.Printf("%s: Success\n", cmd.Command)

					if data != nil {
						ext := cmd.Index[len(cmd.Index)-3 : len(cmd.Index)]
						if ext == "idx" {
							uints := utils.GetUint64Array(data)
							fmt.Printf("Result[%s]: %v\n", ext, uints)
						} else {
							fmt.Printf("Result: %s\n", string(data))
						}
					}
				}

			}
		}
	}

	if enableHistory {
		if f, err := os.Create(homePath + "/.neosearch/" + historyFile); err != nil {
			log.Print("Error writing history file: ", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}

	fmt.Println("Exiting...")
	return nil
}
