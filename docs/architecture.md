# neosearch architecture

NeoSearch is a document-based reverse index, schemaless (at the moment) that
supports field searches across indices. NeoSearch was initially designed
thinking about use cases related to efficient JOIN between indices and easy
index management.

One of the biggest problems in the current reverse index solutions is the
black-box approach of the index format and the lack of tools to interact
directly with them. In the majority of technology for search solutions,
the only interface for the index is a very limited REST API. NeoSearch
solves that designing from the ground-up with a exposed low-level
interface with the storage engine. 

In the image below, you can see the two user interfaces:

![Architecture](https://raw.githubusercontent.com/NeowayLabs/neosearch/master/docs/img/NeoSearch.png)

The low-level interface is done in the `engine` package by
[engine.Engine](https://github.com/NeowayLabs/neosearch/blob/master/engine/engine.go#L37).
This object is responsible for open the index, cache them and execute the
commands in the storage.

Packages `neosearch` and `index` exposes the high-level interface.

If you want hack into neosearch, you need to know the Engine interface very well.

# Engine

The first thing you need to know is that the Engine doesn't know what
storage engine you are using, it only knows the
[store.KVStore](https://github.com/NeowayLabs/neosearch/blob/master/store/store.go#L6)
interface. With NeoSearch is very easy to create an index, add some
documents, get them, and close. See below:

```go
// Simple example using the Engine interface
// Build with:
//     go build -v -tags leveldb
package main

import (
    "fmt"
    "github.com/NeowayLabs/neosearch/engine"
    "github.com/NeowayLabs/neosearch/store"
)

func main() {
    ng := engine.New(&engine.Config{
	    KVConfig: &store.KVConfig{
		  DataDir: "/tmp/",
		  Debug:   false,
	    },
    })

    defer ng.Close()

    command := engine.Command{}
    command.Index = "document.db"
    command.Command = "set"
    command.Key = []byte("hello")
    command.Value = []byte("world")
    
    _, err := ng.Execute(command)

    if err != nil {
        panic(err)
    }

    command.Command = "get"
    command.Key = []byte("hello")

    ret, err := ng.Execute(command)

    if err != nil {
        panic(err)
    }

    fmt.Println("Got: ", string(ret))
}
```

After run this, you can verify that the directory `/tmp/document.db` was
created with the LSM content for the index.

The `engine.Command` is the unique data structure for communicates with
the index, the available commands are below:

* set
* get
* delete
* mergeset

# Indexing steps

The easiest way of explain the internal working is with some examples.

Let's say that we want to create an index called "operating_systems"
and index the five documents below for future searches:

Document 1
```json
{
    "id": 1,
    "name": "Unix",
    "family": "unix",
    "year": 1971,
    "kernel": "unix",
    "kernelType": "monolithic",
    "authors": [
        "Ken Thompson", 
        "Dennis Ritchie",
        "Brian Kernighan",
        "Douglas McIlroy",
        "Joe Ossanna"     
    ]
}
```

Document 2
```json
{
    "id": 2,
    "name": "Plan9 From Outer Space",
    "family": "unix",
    "kernel": "plan9",
    "kernelType": "Hybrid",
    "year": 1992,
    "authors": [
        "Ken Thompson",
        "Dennis Ritchie",
        "Rob Pike",
        "Russ Cox",
        "Dave Presotto", 
        "Phil Winterbottom"
    ]
}
```
Document 3
```json
{
    "id": 3,
    "name": "ArchLinux",
    "family": "unix",
    "year": 2002,
    "kernel": "Linux",
    "kernelType": "monolithic",
    "authors": [
        "Judd Vinet"
    ]
}
```
Document 4
```json
{
    "id": 4,
    "name": "Slackware",
    "family": "unix",
    "kernel": "Linux",
    "kernelType": "monolithic",
    "year": 1993,
    "authors": [
        "Patrick Volkerding"
    ]
}
```
Document 5
```json
{
    "id": 5,
    "name": "Windows NT",
    "family": "windows",
    "kernel": "Windows NT",
    "kernelType": "hybrid",
    "year": 1993,
    "authors": [
        "Dave Cutler",
        "Others"
    ]
}
```

We can easily simulate the process of indexing this five documents with the
[neosearch-cli](https://github.com/NeowayLabs/neosearch/tree/master/neosearch-cli)
tool and a bit of manual process.

If we want to create the index located at /data directory. Then, the
process/algorithm for create and index the documents is like below:

1. Create a new directory at /data/operating_system;
2. For each document, do the steps below:
    1. Stores the entire document at "/data/operating_system/documents.db" with key = &lt;id of document&gt; and value the entire document; See #1
    2. For each field of document, do the steps below:
        1. Create (if not exists) a database at /data/operating_system/&lt;field name&gt;.idx;
        2. Use an analyzer to create a list of tokens of field value;
        3. For each token, do the steps below:
            1. Verify in the field database if already exists documents indexed for that token: `store.Get("&lt;token name&gt;")`. 
            2. Add the current document ID for the set of ID's for the token. (if this is the first time the token was indexed for this field, then create an array with only the document ID)
            3. Stores this token in the <field>.idx, using the token as key and the array of ID's as value. `store.Set("&lt;token&gt;", [&lt;id1&gt;, &lt;id2&gt;, ...])`

The algorithm above can be summarized in the following shell script:

```bash
mkdir -p /data/operating_system/
neosearch-cli -f ./operating_system.ns
```

When operating_system.ns is the file below:

```
using id.idx mergeset "1" 1;
using name.idx mergeset "unix" 1;
using family.idx mergeset "unix" 1;
using year.idx mergeset 1971 1;
using kernel.idx mergeset "unix" 1;
using kerneltype.idx mergeset "monolithic" 1;
# using a simple analyser/tokenizer
using authors.idx mergeset "ken thompson" 1;
using authors.idx mergeset "ken" 1;
using authors.idx mergeset "thompson" 1;
using authors.idx mergeset "dennis ritchie 1;
using authors.idx mergeset "dennis" 1;
using authors.idx mergeset "ritchie" 1;
using authors.idx mergeset "brian kernighan" 1;
using authors.idx mergeset "brian" 1;
using authors.idx mergeset "kernigham" 1;
using authors.idx mergeset "douglas mcIlroy" 1;
using authors.idx mergeset "douglas" 1;
using authors.idx mergeset "mcIlroy" 1;
using authors.idx mergeset "joe ossanna" 1;
using authors.idx mergeset "joe" 1;
using authors.idx mergeset "ossanna" 1;
using id.idx mergeset "2" 2;
using name.idx mergeset "plan9 from outer space" 2;
using name.idx mergeset "plan9" 2;
using name.idx mergeset "from" 2;
using name.idx mergeset "outer" 2;
using name.idx mergeset "space" 2;
using family.idx mergeset "unix" 2;
using kernel.idx mergeset "plan9" 2;
using kerneltype.idx mergeset "hybrid" 2;
using year.idx mergeset 1992 2;
using authors.idx mergeset "ken thompson" 2;
using authors.idx mergeset "ken" 2;
using authors.idx mergeset "thompson" 2;
using authors.idx mergeset "dennis ritchie 2;
using authors.idx mergeset "dennis" 2;
using authors.idx mergeset "ritchie" 2;
using authors.idx mergeset "rob pike" 2;
using authors.idx mergeset "russ cox" 2;
using authors.idx mergeset "dave presotto" 2;
using authors.idx mergeset "phil winterbottom" 2;

# and so on...
```

# package neosearch

Package neosearch is the main user-interface of the library, with that the user can get/create/update/delete index (index.Index) instances. 
