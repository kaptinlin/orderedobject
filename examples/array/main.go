package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Array Operations Example ===")

	// Create an object with array values
	config := orderedobject.NewObject[any]().
		Set("tags", []string{"go", "json", "ordered"}).
		Set("numbers", []int{1, 2, 3, 4, 5}).
		Set("settings", []any{
			orderedobject.NewObject[any]().
				Set("name", "setting1").
				Set("value", 100),
			orderedobject.NewObject[any]().
				Set("name", "setting2").
				Set("value", 200),
		})

	// Access array values
	if tags, found := config.Get("tags"); found {
		if tagArray, ok := tags.([]string); ok {
			fmt.Printf("\nTags: %v\n", tagArray)
		}
	}

	// Access nested array objects
	if settings, found := config.Get("settings"); found {
		if settingsArray, ok := settings.([]any); ok {
			fmt.Println("\nSettings:")
			for i, setting := range settingsArray {
				if settingObj, ok := setting.(*orderedobject.Object[any]); ok {
					if name, found := settingObj.Get("name"); found {
						if value, found := settingObj.Get("value"); found {
							fmt.Printf("  %d. %s = %v\n", i+1, name, value)
						}
					}
				}
			}
		}
	}

	// Print all values
	fmt.Println("\nAll values:")
	config.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})
}
