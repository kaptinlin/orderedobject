package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Map Operations Example ===")

	// Create from map
	settings := map[string]any{
		"theme":         "dark",
		"font_size":     14,
		"notifications": true,
		"language":      "en",
	}

	obj := orderedobject.FromMap(settings)
	fmt.Println("\nFrom map (order preserved):")
	obj.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})

	// Convert back to map
	m := obj.ToMap()
	fmt.Println("\nBack to map (order not preserved):")
	for k, v := range m {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
