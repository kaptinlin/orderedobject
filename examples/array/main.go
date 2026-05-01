package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Array Operations Example ===")

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

	tags, found := config.Get("tags")
	if tagArray, ok := tags.([]string); found && ok {
		fmt.Printf("\nTags: %v\n", tagArray)
	}

	settings, found := config.Get("settings")
	settingsArray, ok := settings.([]any)
	if !found || !ok {
		fmt.Println("\nAll values:")
		config.ForEach(func(key string, value any) {
			fmt.Printf("  %s: %v\n", key, value)
		})
		return
	}

	fmt.Println("\nSettings:")
	for i, setting := range settingsArray {
		settingObj, ok := setting.(*orderedobject.Object[any])
		if !ok {
			continue
		}
		name, found := settingObj.Get("name")
		if !found {
			continue
		}
		value, found := settingObj.Get("value")
		if !found {
			continue
		}
		fmt.Printf("  %d. %s = %v\n", i+1, name, value)
	}

	fmt.Println("\nAll values:")
	config.ForEach(func(key string, value any) {
		fmt.Printf("  %s: %v\n", key, value)
	})
}
