package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/NeowayLabs/neosearch/service/neosearch/server"
	"github.com/jteeuwen/go-pkg-optarg"
)

const (
	DefaultPort = uint16(9500)
	DefaultHost = "0.0.0.0"
)

func main() {
	var (
		configOpt, dataDirOpt string
		kvstoreOpt, hostOpt   string
		goProcsOpt            uint64
		portOpt               uint16
		helpOpt, debugOpt     bool
		err                   error
		cfg                   *config.Config
		cfgServer             *server.ServerConfig
	)

	cfg = config.NewConfig()
	cfgServer = server.NewConfig()

	optarg.Header("General options")
	optarg.Add("c", "config", "Configurations file", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("k", "default-kvstore", "Default kvstore", "")
	optarg.Add("g", "maximum-concurrence", "Set the maximum number of concurrent go routines", 0)
	optarg.Add("t", "trace-debug", "Enable debug traces", false)
	optarg.Add("s", "server-address", "Server host and port", "0.0.0.0:9500")
	optarg.Add("h", "help", "Display this help", false)

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "c":
			configOpt = opt.String()
		case "d":
			dataDirOpt = opt.String()
		case "k":
			kvstoreOpt = opt.String()
		case "s":
			address := opt.String()
			addrParts := strings.Split(address, ":")

			if len(addrParts) > 1 {
				hostOpt = addrParts[0]
				port := addrParts[1]

				portInt, err := strconv.Atoi(port)

				if err == nil {
					portOpt = uint16(portInt)
				} else {
					log.Fatalf("Invalid port number: %s (%s)", port, err)
					return
				}
			} else {
				hostOpt = addrParts[0]
				portOpt = DefaultPort
			}
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
	} else {
		cfg, err = config.ConfigFromFile(configOpt)
		if err != nil {
			log.Fatalf("Failed to read configuration file: %s", err.Error())
			return
		}
	}

	if hostOpt == "" {
		hostOpt = DefaultHost
	}

	if portOpt == 0 {
		portOpt = DefaultPort
	}

	// override config options by argument options
	cfg.Option(config.DataDir(dataDirOpt))
	cfg.Option(config.KVStore(kvstoreOpt))
	cfg.Option(config.Debug(debugOpt))

	cfgServer.Host = hostOpt
	cfgServer.Port = portOpt

	search := neosearch.New(cfg)
	defer func() {
		search.Close()
	}()

	httpServer, err := server.New(search, cfgServer)

	_ = goProcsOpt

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, closing indexes...\n")
			search.Close()
			os.Exit(0)
		}
	}()

	err = httpServer.Start()

	if err != nil {
		log.Fatalf("Failed to start http server: %s", err.Error())
	}
}
