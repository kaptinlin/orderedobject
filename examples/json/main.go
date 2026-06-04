// Package main demonstrates JSON encoding and decoding with orderedobject.
package main

import (
	"fmt"
	"log"

	jsonlib "github.com/go-json-experiment/json"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== JSON Operations Example ===")

	user := orderedobject.New[any]().
		Set("id", 1001).
		Set("name", "John Doe").
		Set("email", "john@example.com").
		Set("active", true)

	fmt.Println("\n1. Using MarshalJSON:")
	data1, err := user.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data1))

	fmt.Println("\n2. Using json.Marshal:")
	data2, err := jsonlib.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data2))

	fmt.Println("\n3. Via map (order not preserved):")
	data3, err := jsonlib.Marshal(user.ToUnorderedMap())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data3))

	jsonStr := `{"name":"Alice","age":30,"skills":["Go","Python"]}`
	parsed, err := orderedobject.FromJSON[any]([]byte(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nParsed JSON:")
	parsed.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})
}
