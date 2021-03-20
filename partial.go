package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

//go:embed partial.rego
var module string

func main() {
	ctx := context.Background()

	// Define a simple policy for example purposes.
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
		fmt.Println(toSQLWhere(pq.Queries[i]))
	}
	stmt := strings.Join(conditions, " OR ")
	fmt.Println(stmt)
}

func toSQL(in ast.Body) string {
	var result []string
	for i := range in {
		expr := in[i]
		if !expr.IsCall() {
			continue
		}
		var op string
		switch v := expr.Operator(); v.String() {
		case "eq":
			op = " = "
		default:
			log.Fatalf("unsupported operator: %s", v)
		}
		// Unfortunately the order is not guaranteed.
		l, r := expr.Operand(0).String(), expr.Operand(1).String()
		if strings.Contains(l, "data.pets[_]") {
			l, r = strings.ReplaceAll(l, "data.pets[_].", ""), r
		} else {
			l, r = strings.ReplaceAll(r, "data.pets[_].", ""), l
		}
		q := strings.Join([]string{l, r}, op)
		result = append(result, q)
	}
	// Produces the following:
	//(owner = "alice" AND name = "fluffy") OR (veterinarian = "alice" AND clinic = "SOMA" AND name = "fluffy")
	return strings.Join(result, " AND ")
}

func toSQLWhere(in ast.Body) map[string]string {
	result := make(map[string]string)
	for i := range in {
		expr := in[i]
		if !expr.IsCall() {
			continue
		}

		// Unfortunately the order is not guaranteed.
		l, r := expr.Operand(0).String(), expr.Operand(1).String()
		if strings.Contains(l, "data.pets[_]") {
			l, r = strings.ReplaceAll(l, "data.pets[_].", ""), r
		} else {
			l, r = strings.ReplaceAll(r, "data.pets[_].", ""), l
		}
		result[l] = r
	}
	return result
}
