package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func main() {
	ctx := context.Background()

	// Define a simple policy for example purposes.
	module := `
	package petclinic.authz

	default allow = false

	allow {
		input.method = "GET"
		input.path = ["pets", name]
		allowed[pet]
		pet.name = name
	}

	allowed[pet] {
		pet = data.pets[_]
		pet.owner = input.subject.user
	}

	allowed[pet] {
		pet = data.pets[_]
		pet.veterinarian = input.subject.user
		pet.clinic = input.subject.location
	}
	`

	r := rego.New(
		rego.Query("data.petclinic.authz.allow == true"),
		rego.Module("petclinic.authz.rego", module),
		rego.Input(map[string]interface{}{
			"method": "GET",
			"path":   []string{"pets", "fluffy"},
			"subject": map[string]interface{}{
				"user":     "alice",
				"location": "SOMA",
			},
		}),
		// The values to treat as unknown during evaluation.
		rego.Unknowns([]string{"data.pets"}),
	)
	// https://blog.openpolicyagent.org/write-policy-in-opa-enforce-policy-in-sql-d9d24db93bf4
	// Perform partial evaluation on the current
	// query, and then use the results to generate
	// SQL queries "where" condition.
	pq, err := r.Partial(ctx)
	if err != nil {
		log.Fatal(err)
	}

	conditions := make([]string, len(pq.Queries))
	for i := range pq.Queries {
		condition := toSQL(pq.Queries[i])
		conditions[i] = fmt.Sprintf("(%s)", condition)
	}
	stmt := strings.Join(conditions, " OR ")
	fmt.Println(stmt)
}

func toSQL(in ast.Body) string {
	result := make([]string, len(in))
	for i := range in {
		// Convert to string.
		expr := fmt.Sprint(in[i])
		q := strings.Replace(expr, "data.pets[_].", "", -1)
		values := strings.Split(q, " = ")
		sort.Sort(sort.Reverse(sort.StringSlice(values)))
		q = strings.Join(values, " = ")
		result[i] = q
	}
	return strings.Join(result, " AND ")
}
