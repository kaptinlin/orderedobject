package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Basic Operations Example ===")

	// Create and populate object
	config := orderedobject.NewObject[any]().
		Set("app_name", "MyApp").
		Set("version", "1.0.0").
		Set("debug", true).
		Set("max_connections", 100)

	// Get and check values
	if version, found := config.Get("version"); found {
		fmt.Printf("Version: %v\n", version)
	}

	// Update existing value
	config.Set("version", "1.0.1")

	// Delete key
	config.Delete("debug")

	// Iterate through entries
	fmt.Println("\nConfiguration:")
	config.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})

	// Get all entries
	fmt.Println("\nAll entries:")
	entries := config.Entries()
	for _, entry := range entries {
		fmt.Printf("  %s: %v\n", entry.Key, entry.Value)
	}

	// Clone and modify
	devConfig := config.Clone().
		Set("debug", true).
		Set("environment", "development")

	fmt.Println("\nDevelopment config:")
	devConfig.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})
}
