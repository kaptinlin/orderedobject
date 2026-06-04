// Package main demonstrates nested ordered objects.
package main

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Nested Structures Example ===")

	config := orderedobject.New[any]().
		Set("app", orderedobject.New[any]().
			Set("name", "MyApp").
			Set("version", "1.0.0").
			Set("debug", true)).
		Set("server", orderedobject.New[any]().
			Set("host", "localhost").
			Set("port", 8080).
			Set("ssl", orderedobject.New[any]().
				Set("enabled", true).
				Set("cert", "/path/to/cert.pem"))).
		Set("database", orderedobject.New[any]().
			Set("driver", "postgres").
			Set("host", "db.example.com").
			Set("port", 5432).
			Set("credentials", orderedobject.New[any]().
				Set("username", "admin").
				Set("password", "secret")))

	app, found := config.Get("app")
	appObj, ok := app.(*orderedobject.Object[any])
	if found && ok {
		printAppName(appObj)
	}

	fmt.Println("\nFull configuration:")
	printNestedObject(config, 0)
}

func printAppName(appObj *orderedobject.Object[any]) {
	name, found := appObj.Get("name")
	if !found {
		return
	}
	fmt.Printf("\nApp name: %v\n", name)
}

func printNestedObject(obj *orderedobject.Object[any], indent int) {
	prefix := strings.Repeat("  ", indent)

	obj.ForEach(func(key string, value any) {
		fmt.Printf("%s%s: ", prefix, key)
		if nested, ok := value.(*orderedobject.Object[any]); ok {
			fmt.Println()
			printNestedObject(nested, indent+1)
		} else {
			fmt.Printf("%v\n", value)
		}
	})
}
