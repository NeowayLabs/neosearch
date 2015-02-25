package neosearch_test

import (
	"fmt"
	"os"

	"github.com/neowaylabs/neosearch"
)

func Example() {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir(dataDir))
	cfg.Option(neosearch.Debug(false))

	neo := neosearch.New(cfg)

	index, err := neo.CreateIndex("test")

	if err != nil {
		panic(err)
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		panic(err)
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`))

	if err != nil {
		panic(err)
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`))

	if err != nil {
		panic(err)
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`))

	if err != nil {
		panic(err)
	}

	data, err := index.Get(1)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	// Output:
	// {"id": 1, "name": "Neoway Business Solution"}

	neo.Close()

	os.RemoveAll(dataDir)
}

func ExampleMatchPrefix() {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir(dataDir))
	cfg.Option(neosearch.Debug(false))

	neo := neosearch.New(cfg)

	index, err := neo.CreateIndex("test")

	if err != nil {
		panic(err)
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		panic(err)
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`))

	if err != nil {
		panic(err)
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`))

	if err != nil {
		panic(err)
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`))

	if err != nil {
		panic(err)
	}

	values, err := index.MatchPrefix([]byte("name"), []byte("neoway"))

	if err != nil {
		panic(err)
	}

	for _, value := range values {
		fmt.Println(value)
	}

	// Output:
	// {"id": 1, "name": "Neoway Business Solution"}
	// {"id": 4, "name": "Neoway Teste"}

	neo.Close()

	os.RemoveAll(dataDir)
}
