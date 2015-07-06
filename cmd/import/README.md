# neosearch-import

# Build

```
go get -u github.com/NeowayLabs/neosearch
go build -v -tags leveldb
```

# usage

```
$ ./neosearch-import
[General options]
     --file, -f: Read NeoSearch JSON database from file. (Required)
   --create, -c: Create new index database
     --name, -n: Name of index database
 --data-dir, -d: Data directory
     --help, -h: Display this help
```

Indexing the sample file:
```
$ mkdir /tmp/data
$ ./neosearch-import -f samples/operating_systems.json  -c -d /tmp/data -n operating-systems
```

# How to verify the indexed data?

Use the [neosearch-cli](https://github.com/NeowayLabs/neosearch-cli) tool:

```
$ go get -v -tags leveldb github.com/NeowayLabs/neosearch-cli
$ neosearch-cli -d /tmp/data

neosearch>using document.db get 0
get: Success
Result: {"_id":0,"authors":["Ken Thompson","Dennis Ritchie","Brian Kernighan","Douglas McIlroy","Joe Ossanna"],"family":"unix","id":1,"kernel":"unix","kernelType":"monolithic","name":"Unix","year":1971}
neosearch>using document.db get 1
get: Success
Result: {"_id":1,"authors":["Ken Thompson","Dennis Ritchie","Rob Pike","Russ Cox","Dave Presotto","Phil Winterbottom"],"family":"unix","id":2,"kernel":"plan9","kernelType":"Hybrid","name":"Plan9 From Outer Space","year":1992}
neosearch>using document.db get 2
get: Success
Result: {"_id":2,"authors":["Judd Vinet"],"family":"unix","id":3,"kernel":"Linux","kernelType":"monolithic","name":"ArchLinux","year":2002}
neosearch>using document.db get 3
get: Success
Result: {"_id":3,"authors":["Patrick Volkerding"],"family":"unix","id":4,"kernel":"Linux","kernelType":"monolithic","name":"Slackware","year":1993}
neosearch>using document.db get 4
get: Success
Result: {"_id":4,"authors":["Dave Cutler","Others"],"family":"windows","id":5,"kernel":"Windows NT","kernelType":"hybrid","name":"Windows NT","year":1993}
neosearch>
neosearch>
neosearch>using name.idx get "plan9"
get: Success
Result[idx]: [1]
neosearch>
```
