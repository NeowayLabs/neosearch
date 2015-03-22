package main

import (
	"log"
	"os"
	"strconv"

	"github.com/NeowayLabs/neosearch"
	"github.com/jteeuwen/go-pkg-optarg"
)

func main() {
	var (
		configOpt, dataDirOpt string
		goProcsOpt            uint64
		helpOpt, debugOpt     bool
		err                   error
		cfg                   *neosearch.Config
	)

	cfg = neosearch.NewConfig()

	optarg.Header("General options")
	optarg.Add("c", "config", "Configurations file", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("g", "maximum-concurrence", "Set the maximum number of concurrent go routines", 0)
	optarg.Add("t", "trace-debug", "Enable debug traces", false)
	optarg.Add("h", "help", "Display this help", false)

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "c":
			configOpt = opt.String()
		case "d":
			dataDirOpt = opt.String()
		case "t":
			debugOpt = opt.Bool()
		case "g":
			goprocsStr := opt.String()
			goProcsInt, err := strconv.Atoi(goprocsStr)

			if err != nil || goProcsInt <= 0 {
				log.Fatal("Invalid -g option. Should be a unsigned integer value greater than zero.")
				return
			}

			goProcsOpt = uint64(goProcsInt)
		case "h":
			helpOpt = true
		}
	}

	if helpOpt {
		optarg.Usage()
		os.Exit(0)
	}

	if dataDirOpt == "" {
		dataDirOpt, _ = os.Getwd()
	}

	if configOpt == "" {
		log.Println("No configuration file supplied. Using defaults...")
		cfg.Debug = false
		cfg.DataDir = "/data"
		cfg.EnableCache = true
	} else {
		cfg, err := neosearch.ConfigFromFile(fileOpt)

		if err != nil {
			log.Fatalf("Failed to read configuration file: %s", err.Error())
			return
		}
	}

	// override config options by argument options
	cfg.Option(neosearch.DataDir(dataDirOpt))
	cfg.Option(neosearch.Debug(debugOpt))

	neo := neosearch.New(cfg)
}
