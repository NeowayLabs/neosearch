package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/neosearch/server"
	"github.com/jteeuwen/go-pkg-optarg"
)

const (
	DefaultPort = uint16(9500)
	DefaultHost = "127.0.0.1"
)

func main() {
	var (
		configOpt, dataDirOpt, hostOpt string
		goProcsOpt                     uint64
		portOpt                        uint16
		helpOpt, debugOpt              bool
		err                            error
		cfg                            *neosearch.Config
		cfgServer                      *server.ServerConfig
	)

	cfg = neosearch.NewConfig()
	cfgServer = server.NewConfig()

	optarg.Header("General options")
	optarg.Add("c", "config", "Configurations file", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("g", "maximum-concurrence", "Set the maximum number of concurrent go routines", 0)
	optarg.Add("t", "trace-debug", "Enable debug traces", false)
	optarg.Add("s", "server-address", "Server host and port", "127.0.0.1:9500")
	optarg.Add("h", "help", "Display this help", false)

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "c":
			configOpt = opt.String()
		case "d":
			dataDirOpt = opt.String()
		case "s":
			address := opt.String()
			addrParts := strings.Split(address, ":")

			if len(addrParts) > 1 {
				hostOpt = addrParts[0]
				port := addrParts[1]

				portInt, err := strconv.Atoi(port)

				if err != nil {
					portOpt = uint16(portInt)
				} else {
					log.Fatalf("Invalid port number: %s", port)
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
		cfg.EnableCache = true
	} else {
		cfg, err = neosearch.ConfigFromFile(configOpt)

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
	cfg.Option(neosearch.DataDir(dataDirOpt))
	cfg.Option(neosearch.Debug(debugOpt))

	cfgServer.Host = hostOpt
	cfgServer.Port = portOpt

	search := neosearch.New(cfg)

	httpServer, err := server.New(search, cfgServer)

	_ = goProcsOpt

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = httpServer.Start()

	if err != nil {
		log.Fatalf("Failed to start http server: %s", err.Error())
	}
}
