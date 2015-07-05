# Proposal for dump/restore

We need a way of dump and restore of indices database. Today we have a low level interface using `engine.Command` to communicate with storage. This is the interface used by [neosearch-cli](https://github.com/NeowayLabs/neosearch-cli) to access the indices. With `neosearch-cli` we can process commands stored in a text file with neosearch command syntax like below:

```
using title.idx set "hello" 1;
using id.idx set 1 "{test: \"1\"}";
using title.idx mergeset "hello" 2;
using title.idx mergeset "hello" 10;
using document.db set 1 "{\"title\": \"hello\", \"id\": 1}";
```

As said [here](https://github.com/NeowayLabs/neosearch/wiki/Internal-Concepts#indexing-steps) we can simulate the internal process of indexing a document in neosearch with commands stored in a file and processed with `neosearch-cli`. We can generate a huge `index.ns` file with all of the commands needed to re-index the database.

To implement the `dump` feature, we only need to get a [Iterator](https://github.com/NeowayLabs/neosearch/blob/master/store/store.go#L15) in each index database and create `engine.Command` entries in a file with the `neosearch-cli` syntax (like [this](https://github.com/NeowayLabs/neosearch-cli/blob/master/index_data.ns)).

To implement the `restore` feature we only need to add parallelism to `neosearch-cli` tool to process the dumped file. As both the neosearch library and neosearch-cli tool uses the `engine.Command` to interact with storage we can guarantee that this will works as expected.
