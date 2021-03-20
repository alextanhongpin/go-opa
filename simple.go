package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
)

//go:embed example.rego
var example string

func main() {
	ctx := context.Background()
	rego := rego.New(
		rego.Module("example.rego", example),
		rego.Query(`data.example.admins[_] == "john"`),
	)
	rs, err := rego.Eval(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", rs)
}
