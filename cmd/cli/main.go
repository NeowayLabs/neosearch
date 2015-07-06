package main

import (
	"os"

	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/jteeuwen/go-pkg-optarg"
)

func main() {
	var fileOpt, dataDirOpt, homeOpt string
	var helpOpt, debugOpt bool

	optarg.Add("f", "from-file", "Read NeoSearch low-level instructions from file", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("t", "trace-debug", "Enable trace for debug", false)
	optarg.Add("h", "help", "Display this help", false)
	optarg.Add("m", "home", "User home for store command history", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "f":
			fileOpt = opt.String()
			break
		case "d":
			dataDirOpt = opt.String()
			break
		case "m":
			homeOpt = opt.String()
			break
		case "t":
			debugOpt = true
			break
		case "h":
			helpOpt = true
			break
		}
	}

	if helpOpt {
		optarg.Usage()
		os.Exit(0)
	}

	if homeOpt == "" {
		if homeEnv := os.Getenv("HOME"); homeEnv != "" {
			homeOpt = homeEnv
		}
	}

	if dataDirOpt == "" {
		dataDirOpt, _ = os.Getwd()
	}

	ng := engine.New(engine.NGConfig{
		KVCfg: &store.KVConfig{
			DataDir: dataDirOpt,
			Debug:   debugOpt,
		},
	})

	defer ng.Close()

	if fileOpt != "" {
		batch(ng, fileOpt)
	} else {
		cli(ng, homeOpt)
	}
}
