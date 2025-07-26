# Notes

## Setting up Dgraph instance with ratel in docker in makefile

```make
DGRAPH := dgraph/standalone:latest
RATEL  := dgraph/ratel:latest


# dev-docker: pull required docker images

dev-docker:
    docker pull $(DGRAPH)
    docker pull $(RATEL)


#dev-docker: run the containers in docker
dev-dgraph-local:
    docker run --name dgraph-local -d -it \
        -p "9999:8080" -p "9080:9080" -p "8090:8090" \
        -v ~/dgraph $(DGRAPH)
    docker run --name ratel-local \
        --platform linux/amd64 -d -p "8000:8000" $(RATEL)
```



## DGraph Query Syntax

Dgraph ql queries have an optional prarmaterization section followed by a non-zero amount of query blocks.

A DQL query has
- [ ] an optional parameterization, ie a name and a list of parameters
- [ ] an opening curly bracket
- [ ] at least one query block, but can contain many blocks
- [ ] optional var blocks
- [ ] a closing curly bracket

**Sample Query**
```
query myQuery ($name : string, $terms : string = "Jerry")                       // Parameterization block
{                                                                               // Start query block
    queryList(func: allofterms(name@en, $terms), first:2) @filter(lt(age,$age)) // block name + root criteria + ordering + filters
        name@en                                                                 // fetch attributes
        age
        friends_with {                                                          // relationship to fetch
            name@en
        }
    }
}                                                                               // End query block
```


### Paramaterization
Parameters
- [ ] must have a name starting with a $ symbol.
- [ ] must have a type int, float, bool or string.
- [ ] may have a default value. In the example below, $age has a default value of 95
- [ ] may be mandatory by suffixing the type with a !. Mandatory parameters can’t have a default value.

- you can only query where a string, float, int, or bool is needed
- if looking to search by uids assign string variable with value quoted list in square brackets
    - `query title($uidsParam: string = "[0x1, 0x2, 0x3]"){...}`
- assigning a paramter with ! makes the paramater mandatory

### Query Block

A query block

- [ ] must have name
- [ ] must have a node criteria defined by the keyword func:
- [ ] may have ordering and pagination information
- [ ] may have a combination of filters (to apply to the root nodes)
- [ ] must provide the list of attributes and relationships to fetch for each node matching the root node

- if predicate has special character wrap it in < >
    - `<https://myschema.org#name`

#### Root Criteria and Filter Functions

- string attributes
    - match terms
        - allofterms - match strings that have a specified term in any order and case insensitive
            - `allofterms(predicate, "space-seperated term list")`
        - anyofterms - match string that has any of the specified terms in any order and case insensitive
            - `anyofterms(predicate, "space-seperated term list")`
    - regular expression
        - regexp - match string by regular expression (go regex)
            - `regexp(predicate, /regexp/)`
    - fuzzy match
        - match - match string by fuzzy
            - `match(predicate, string, distance)`
    - full-text search
        - alloftext - apply full-text search with stemming and stop words to find strings matching all of the text (bleve full-text search indexing)
            - `alloftext(predicate, "space-seperated text)"`
        - anyoftext - apply full-text search with stemming and stop words to find strings match any of the text
???LINES MISSING
    - the `@if` directive accepts a condition on a variable defined in the query block using `AND, OR, NOT` expressions
        - this is an optional conditional statement that will only be executed if the condition is met
```
{
    upsert {
        query <query block>
        [fragment <fragment block>]
        mutation [@if(<condition>)] <mutation block1>
        [mutation [@if(<condition>)] <mutation block2>]
        ...
    }
}
```


## DQL and RDF QL Format

RDF = Resource Description Framework - Semantic WEb Standard for data interchange
    - expressive statements describing resources in the form of triples

- **triple** - one singular facet about a node (subject)
       `<subject> <predicate> <object>` 
    - **subject** - the node represented by a numeric UID
    - **object** - another node or a literal value
        - another node - represents a relationship between two nodes
            - `<0x01> <knows> <0x02> .`
        - literal value - assigns value to given predicate
            - `<0x01> <name> "Alice" .`
    - **predicate** - smallest piece of information about an object - liternal value or a relation to another entity

Types can be specified using the `^^` operator
`<0x01> <age> "21"^^<xs:int> .`

- string - `<xs:string>`
- dateTime - `<xs:dateTime>`
- int - `<xs:int>` | `<xs:integer>`
- bool - `<xs:bool>`
- float - `<xs:float>` | `<xs:double>`
- geo - `<geo:geojson>`
- password - `<xs:password>`


### Facets

Facets are a way to extend the RDF standard to allow "sub-properties" of an object
```
{
    set {
        _:Julian <name> "Julian" .
        _:Julian <nickname> "Jay-Jay" (kind="first") .
        _:Julian <nickname> "Jules" (kind="official") .
        _:Julian <nickname> "JB" (kind="CS-GO") .
    }
}
{
    q(func: eq(name, "Julian")){
        name
        nickname @facets
    }
}
```

JSON Response:
```json
{
    "data": {
        "q": [
            {
                "name": "Julian",
                "nickname|kind": {
                    "0": "first",
                    "1": "official",
                    "2": "CS-GO"
                },
                "nickname": [
                    "Jay-Jay",
                    "Jules",
                    "JB"
                ]
            }
        ]
    }
}
```

