package neosearch_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
)

func OnErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func Example() {
	dataDir, err := ioutil.TempDir("", "neosearchExample")
	defer os.RemoveAll(dataDir)

	OnErrorPanic(err)

	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir(dataDir))
	cfg.Option(neosearch.Debug(false))

	neo := neosearch.New(cfg)
	defer neo.Close()

	index, err := neo.CreateIndex("test")
	OnErrorPanic(err)

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`), nil)
	OnErrorPanic(err)

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`), nil)
	OnErrorPanic(err)

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`), nil)
	OnErrorPanic(err)

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`), nil)
	OnErrorPanic(err)

	data, err := index.Get(1)
	OnErrorPanic(err)

	fmt.Println(string(data))
	// Output:
	// {"id": 1, "name": "Neoway Business Solution"}
}

func ExampleMatchPrefix() {
	dataDir, err := ioutil.TempDir("", "neosearchExample")
	defer os.RemoveAll(dataDir)

	OnErrorPanic(err)

	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir(dataDir))
	cfg.Option(neosearch.Debug(false))

	neo := neosearch.New(cfg)
	defer neo.Close()

	index, err := neo.CreateIndex("test")
	OnErrorPanic(err)

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`), nil)
	OnErrorPanic(err)

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`), nil)
	OnErrorPanic(err)

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`), nil)
	OnErrorPanic(err)

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`), nil)
	OnErrorPanic(err)

	values, err := index.MatchPrefix([]byte("name"), []byte("neoway"))
	OnErrorPanic(err)

	for _, value := range values {
		fmt.Println(value)
	}

	// Output:
	// {"id": 1, "name": "Neoway Business Solution"}
	// {"id": 4, "name": "Neoway Teste"}
}
