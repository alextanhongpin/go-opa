package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage/inmem"
)

func main() {
	// Define a role-based access control (RBAC) policy that decides whether to
	// allow or deny requests. Requests are allowed if the user is bound to a
	// role that grants permission to perform the operation on the resource.
	ctx := context.Background()
	module := `
	package example

	import data.bindings
	import data.roles

	default allow = false

	allow {
		user_has_role[role_name]
		role_has_permission[role_name]
	}

	user_has_role[role_name] {
		b = bindings[_]
		b.role = role_name
		b.user = input.subject.user
	}

	role_has_permission[role_name] {
		r = roles[_]
		r.name = role_name
		match_with_wildcard(r.operations, input.operation)
		match_with_wildcard(r.resources, input.resource)
	}

	match_with_wildcard(allowed, value) {
		allowed[_] = "*"
	}

	match_with_wildcard(allowed, value) {
		allowed[_] = value
	}
	`

	store := inmem.NewFromReader(bytes.NewBufferString(`{
	"roles": [
		{
			"resources": ["documentA", "documentB"],
			"operations": ["read"],
			"name": "analyst"
		},
		{
			"resources": ["*"],
			"operations": ["*"],
			"name": "admin"
		}
	],
	"bindings": [
		{
			"user": "bob",
			"role": "admin"
		},
		{
			"user": "alice",
			"role": "analyst"
		}
	]
}`))

	r := rego.New(
		rego.Query("data.example.allow"),
		rego.Module("example.rego", module),
		rego.Store(store),
	)
	pr, err := r.PartialResult(ctx)
	if err != nil {
		log.Fatal(err)
	}

	examples := []map[string]interface{}{
		{
			"resource":  "documentA",
			"operation": "write",
			"subject": map[string]interface{}{
				"user": "bob",
			},
		},
		{
			"resource":  "documentB",
			"operation": "write",
			"subject": map[string]interface{}{
				"user": "alice",
			},
		},
		{
			"resource":  "documentB",
			"operation": "read",
			"subject": map[string]interface{}{
				"user": "alice",
			},
		},
	}

	for i := range examples {
		r := pr.Rego(
			rego.Input(examples[i]),
		)

		rs, err := r.Eval(ctx)
		if err != nil || len(rs) != 1 || len(rs[0].Expressions) != 1 {
			log.Fatal(err)
		} else {
			fmt.Printf("input %d allowed: %v\n", i+1, rs[0].Expressions[0].Value)
		}
	}
}
