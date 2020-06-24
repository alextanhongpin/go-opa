package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
)

func main() {
	ctx := context.Background()

	module := `
package example

management_chain = {
  "management_chain": {
    "bob": [
      "ken",
      "janet"
    ],
    "alice": [
      "janet"
    ]
  }
}

default allow = false
allow {
	input.method = "GET"
	input.path = ["salary", id]
	input.user_id = id
}

allow {
	input.method = "GET"
	input.path = ["salary", id]
	managers = management_chain[id]
	input.user_id = managers[_]
}
`

	//raw := `{"method": "GET", "path": ["salary", "bob"], "user_id": "bob"}`
	raw := `{"method": "GET", "path": ["salary", "bob"], "user_id": "alice"}`
	d := json.NewDecoder(bytes.NewBufferString(raw))
	// Numeric values must be represented using json.Number.
	d.UseNumber()

	var input interface{}
	if err := d.Decode(&input); err != nil {
		panic(err)
	}

	rego := rego.New(
		rego.Module("example.rego", module),
		rego.Query("data.example.allow == true"),
		rego.Input(input),
	)
	rs, err := rego.Eval(ctx)
	if err != nil {
		log.Fatal("evalErr:", err)
	}
	fmt.Println(len(rs))
	fmt.Printf("%+v\n", rs)
	fmt.Println(rs[0].Expressions[0])
}
