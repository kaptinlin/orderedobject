package main

import (
	"fmt"

	jsonlib "github.com/go-json-experiment/json"
	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== JSON Operations Example ===")

	// Create object
	user := orderedobject.NewObject[any]().
		Set("id", 1001).
		Set("name", "John Doe").
		Set("email", "john@example.com").
		Set("active", true)

	// Convert to JSON (three ways)
	fmt.Println("\n1. Using ToJSON:")
	data1, _ := user.ToJSON()
	fmt.Println(string(data1))

	fmt.Println("\n2. Using json.Marshal:")
	data2, _ := jsonlib.Marshal(user)
	fmt.Println(string(data2))

	fmt.Println("\n3. Via map (order not preserved):")
	m := user.ToMap()
	data3, _ := jsonlib.Marshal(m)
	fmt.Println(string(data3))

	// Parse JSON
	jsonStr := `{"name":"Alice","age":30,"skills":["Go","Python"]}`
	parsed, _ := orderedobject.FromJSON[any]([]byte(jsonStr))
	fmt.Println("\nParsed JSON:")
	parsed.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})
}
