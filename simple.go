package main

import (
	"context"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
)

func main() {
	ctx := context.Background()
	rego := rego.New(
		rego.Module("example.rego", `
package example

admins = ["john"]`),
		rego.Query(`data.example.admins[_] == "john"`),
	)
	rs, err := rego.Eval(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", rs)
}
