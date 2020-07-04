package main

import (
	"context"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/util"
)

func main() {
	ctx := context.Background()

	data := `{
    "example": {
        "users": [
            {
                "name": "alice",
                "likes": ["dogs", "clouds"]
            },
            {
                "name": "bob",
                "likes": ["pizza", "cats"]
            }
        ]
    }
}`

	var json map[string]interface{}
	err := util.UnmarshalJSON([]byte(data), &json)
	if err != nil {
		log.Fatal(err)
	}

	// Manually creating the storage layer.
	store := inmem.NewFromObject(json)
	rego := rego.New(
		rego.Store(store),
		rego.Query("data.example.users[0].likes"),
	)

	rs, err := rego.Eval(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", rs)
}
