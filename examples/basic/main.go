// Package main demonstrates basic orderedobject operations.
package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Basic Operations Example ===")

	config := orderedobject.New[any]().
		Set("app_name", "MyApp").
		Set("version", "1.0.0").
		Set("debug", true).
		Set("max_connections", 100)

	if version, found := config.Get("version"); found {
		fmt.Printf("Version: %v\n", version)
	}

	config.Set("version", "1.0.1")
	config.Delete("debug")

	fmt.Println("\nConfiguration:")
	config.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})

	fmt.Println("\nAll entries:")
	for _, entry := range config.Entries() {
		fmt.Printf("  %s: %v\n", entry.Key, entry.Value)
	}

	devConfig := config.Clone().
		Set("debug", true).
		Set("environment", "development")

	fmt.Println("\nDevelopment config:")
	devConfig.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})
}
