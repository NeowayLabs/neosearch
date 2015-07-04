# Motivation

To understand the motivation behind NeoSearch's creation, we need a bit of
background about the project [Lucene](http://lucene.apache.org/) and the type
of problems for which it doesn't work.

## Data Join

Lucene and SOLR were used internally at [Neoway](http://www.neoway.com.br)
for over five years, and during this time, it was the only mature tech for
full-text search that we could find. When we had only one main index,
which stored all the search information, the SOLR solved our problem very well.
But as the company grew, and the information captured by our robots became
more structured, the flat characteristic of Lucene / SOLR began to show it's
cracks. In short, Lucene was not designed to JOIN between different indices.
All current solutions to this problem, both in SOLR and ElasticSearch, are
workarounds to solve a problem in an architecture that is not designed to
solve that problem. At first we tried to arrange the information in separate
indices and use the "Join" syntax available in SOLR-4 to search relations
between them. But in this way we completely
[lost the ability to scale horizontally](https://wiki.apache.org/solr/DistributedSearch#line-38).
The actual solution presented for this by SOLR and ElasticSearch is the
parent-child relationship between documents. This technique is a better
approach, but, in the same way, the index doesn't scale correctly across
shards and requires a special way to index documents that have relationships.
Some problems are:

* The child documents are always stored in the same shard of parent documents;
* All map of parent and child IDs are stored in memory;
* Child document is limited to have only one parent;

To explain these limitations, think in the following example of indexing USA
population and companies:

Imagine we have an index called "people" that have 310 million entries and we
have other index called "company" that have 31.6 million entries. The company
index has a relationship with the people index by the "partner" field and
"employ" field. 

* Each company have one or more partners in the people index;
* Each company have zero or more employees in the people index;

Using the solution available in Lucene indices, we have to first index the
parent documents, in this case `company` documents, and then index two others
indices for people. The first for index the partners and the other for index
the employees. For each partner, we will index in the `people_partner` index
specifying the correct company parent. And, for each employee people we will
index in the "people_employee" index specifying the parent company document.

* So, the first problem that arise is that we will end with irregular shards.
ElasticSearch put the parent and children documents in the same shard, then the
size (MEM, CPU, Disk, etc.) needed by the shard machines isn't predictable
because each company have a different number of employees and partners.

* Another problem is that for each relationship with people index, you will
need to replicate the information in another index (like `partner_people`,
`employee_people`, `former_employee_people`, etc.).

* If the information of one entry in the people index change, we will need
to update this information in every `people like` index. **Critical** in our
data model.

* For each parent-child relationship, ElasticSearch will maintain the parent
IDs (in string format) and child IDs (8 bytes per ID) for each relation in
memory. This implementation can be a serious problem if the parent index has
a lot of relations. In the case above is ~ **4.75 GB** only for the memory map
if we consider an average of 3 partners and three employees per company.

The item above is only one example that shows that relationships are a big
problem in the current search solutions. In the business intelligence field,
we need to cross a lot of information to find patterns, trends, frauds, etc.,
and duplicate all of that information on the indices isn't an option. We know
that search engines aren't relational databases, but to manage relationships
in a reverse index is crucial today, and for this reason ElasticSearch and
SOLR support workarounds for this.

Neoway is a Big Data company that scrapes the internet 
NeoSearch was born in this context. We seek for a reverse index solution that
manages relationships and index updates efficiently.
