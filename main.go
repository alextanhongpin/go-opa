package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/util"
)

func main() {
	{
		input := `{"method": "GET", "path": ["salary", "bob"], "user_id": "bob"}`
		run(input)
	}
	{

		input := `{"method": "GET", "path": ["salary", "bob"], "user_id": "alice"}`
		run(input)
	}
}

func run(rawInput string) {
	ctx := context.Background()

	module := `
package example

default allow = false
allow {
	input.method = "GET"
	input.path = ["salary", id]
	input.user_id = id
}

allow {
	input.method = "GET"
	input.path = ["salary", id]
	managers = data.management_chain[id]
	input.user_id = managers[_]
}
`

	d := json.NewDecoder(bytes.NewBufferString(rawInput))
	// Numeric values must be represented using json.Number.
	d.UseNumber()

	var input interface{}
	if err := d.Decode(&input); err != nil {
		panic(err)
	}

	r := rego.New(
		rego.Module("example.rego", module),
		rego.Store(newStore()),
		rego.Query("data.example.allow == true"),
		rego.Input(input),
	)
	rs, err := r.Eval(ctx)
	if err != nil {
		log.Fatal("evalErr:", err)
	}
	fmt.Println(len(rs))
	fmt.Printf("%+v\n", rs)
	fmt.Println(rs[0].Expressions[0])
}

// Creates a new inmem store for the data required for rules.
func newStore() storage.Store {
	data := `{
  "management_chain": {
    "bob": [
      "ken",
      "janet"
    ],
    "alice": [
      "janet"
    ]
  }
}`

	var json map[string]interface{}

	err := util.UnmarshalJSON([]byte(data), &json)
	if err != nil {
		// Handle error.
	}

	// Manually create the storage layer. inmem.NewFromObject returns an
	// in-memory store containing the supplied data.
	store := inmem.NewFromObject(json)
	return store
}
