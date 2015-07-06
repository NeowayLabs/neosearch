package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"launchpad.net/gommap"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/index"
	"github.com/jteeuwen/go-pkg-optarg"
)

func main() {
	var (
		fileOpt,
		dataDirOpt,
		databaseName,
		profileFile,
		metadataStr string
		metadata                    = index.Metadata{}
		helpOpt, newIndex, debugOpt bool
		err                         error
		index                       *index.Index
		batchSize                   int
	)

	optarg.Header("General options")
	optarg.Add("f", "file", "Read NeoSearch JSON database from file. (Required)", "")
	optarg.Add("c", "create", "Create new index database", false)
	optarg.Add("b", "batch-size", "Batch size", 1000)
	optarg.Add("n", "name", "Name of index database", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("t", "trace-debug", "Enable trace for debug", false)
	optarg.Add("h", "help", "Display this help", false)
	optarg.Add("p", "cpuprofile", "write cpu profile to file", "")
	optarg.Add("m", "metadata", "metadata of documents", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "f":
			fileOpt = opt.String()
		case "b":
			batchSize = opt.Int()
		case "d":
			dataDirOpt = opt.String()
		case "n":
			databaseName = opt.String()
		case "c":
			newIndex = true
		case "t":
			debugOpt = true
		case "p":
			profileFile = opt.String()
		case "m":
			metadataStr = opt.String()
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

	if fileOpt == "" {
		optarg.Usage()
		os.Exit(1)
	}

	if profileFile != "" {
		f, err := os.Create(profileFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Profiling to file: ", profileFile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if metadataStr != "" {
		err = json.Unmarshal([]byte(metadataStr), &metadata)

		if err != nil {
			log.Fatal(err)
		}
	}

	cfg := neosearch.NewConfig()

	cfg.Option(neosearch.DataDir(dataDirOpt))
	cfg.Option(neosearch.Debug(debugOpt))
	cfg.Option(neosearch.KVCacheSize(1 << 15))

	neo := neosearch.New(cfg)

	if newIndex {
		log.Printf("Creating index %s\n", databaseName)
		index, err = neo.CreateIndex(databaseName)
	} else {
		log.Printf("Opening index %s ...\n", databaseName)
		index, err = neo.OpenIndex(databaseName)
	}

	if err != nil {
		log.Fatalf("Failed to open database '%s': %v", err)
		return
	}

	file, err := os.OpenFile(fileOpt, os.O_RDONLY, 0)

	if err != nil {
		log.Fatalf("Unable to open file: %s", fileOpt)
		return
	}

	jsonBytes, err := gommap.Map(file.Fd(), gommap.PROT_READ,
		gommap.MAP_PRIVATE)

	if err != nil {
		panic(err)
	}

	data := make([]map[string]interface{}, 0)

	err = json.Unmarshal(jsonBytes, &data)

	if err != nil {
		panic(err)
	}

	jsonBytes = nil

	startTime := time.Now()

	index.Batch()
	var count int
	totalResults := len(data)

	runtime.GC()

	cleanup := func() {
		neo.Close()
		file.Close()
		if profileFile != "" {
			fmt.Println("stopping profile: ", profileFile)
			pprof.StopCPUProfile()
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic", r)
			cleanup()
			os.Exit(1)
		}

		cleanup()
	}()

	fmt.Println("Importing ", len(data), " records")

	for idx := range data {
		dataEntry := data[idx]

		if dataEntry["_id"] == nil {
			dataEntry["_id"] = idx
		}

		entryJSON, err := json.Marshal(&dataEntry)
		if err != nil {
			log.Println(err)
			return
		}

		err = index.Add(uint64(idx), entryJSON, metadata)
		if err != nil {
			panic(err)
		}

		if count == batchSize {
			count = 0

			fmt.Println("Flushing batch: ", idx, " from ", totalResults)
			index.FlushBatch()
			if idx != (totalResults - 1) {
				index.Batch()
			}

			runtime.GC()
		} else {
			count = count + 1
		}

		data[idx] = nil
	}

	index.FlushBatch()
	index.Close()
	neo.Close()

	elapsed := time.Since(startTime)

	log.Printf("Database indexed in %v\n", elapsed)
}
