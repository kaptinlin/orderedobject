package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Nested Structures Example ===")

	// Create a nested configuration
	config := orderedobject.NewObject[any]().
		Set("app", orderedobject.NewObject[any]().
			Set("name", "MyApp").
			Set("version", "1.0.0").
			Set("debug", true)).
		Set("server", orderedobject.NewObject[any]().
			Set("host", "localhost").
			Set("port", 8080).
			Set("ssl", orderedobject.NewObject[any]().
				Set("enabled", true).
				Set("cert", "/path/to/cert.pem"))).
		Set("database", orderedobject.NewObject[any]().
			Set("driver", "postgres").
			Set("host", "db.example.com").
			Set("port", 5432).
			Set("credentials", orderedobject.NewObject[any]().
				Set("username", "admin").
				Set("password", "secret")))

	// Access nested values
	if app, found := config.Get("app"); found {
		if appObj, ok := app.(*orderedobject.Object[any]); ok {
			if name, found := appObj.Get("name"); found {
				fmt.Printf("\nApp name: %v\n", name)
			}
		}
	}

	// Print nested structure
	fmt.Println("\nFull configuration:")
	printNestedObject(config, 0)
}

func printNestedObject(obj *orderedobject.Object[any], indent int) {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	obj.ForEach(func(key string, value any) {
		fmt.Printf("%s%s: ", indentStr, key)
		if nested, ok := value.(*orderedobject.Object[any]); ok {
			fmt.Println()
			printNestedObject(nested, indent+1)
		} else {
			fmt.Printf("%v\n", value)
		}
	})
}
