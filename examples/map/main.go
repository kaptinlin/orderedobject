// Package main demonstrates map conversion with orderedobject.
package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Map Operations Example ===")

	settings := map[string]any{
		"theme":         "dark",
		"font_size":     14,
		"notifications": true,
		"language":      "en",
	}

	obj := orderedobject.FromMap(settings)
	fmt.Println("\nFrom map (order follows Go map iteration):")
	obj.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})

	fmt.Println("\nBack to map (order not preserved):")
	for k, v := range obj.ToMap() {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
