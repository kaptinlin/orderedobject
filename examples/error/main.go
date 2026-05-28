// Package main demonstrates orderedobject error handling.
package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Error Handling Example ===")

	config := printConfig(`{"name": "test", "port": "not_a_number"}`)
	printPort(config)
	printLookup(config, "nonexistent")

	nestedConfig := orderedobject.NewObject[any]().
		Set("server", orderedobject.NewObject[any]().
			Set("port", 8080))
	printServerPort(nestedConfig)
}

func printConfig(data string) *orderedobject.Object[any] {
	config, err := orderedobject.FromJSON[any]([]byte(data))
	if err != nil {
		fmt.Printf("\nParse error: %v\n", err)
		return orderedobject.NewObject[any]()
	}
	fmt.Println("\nSuccessfully parsed JSON")
	return config
}

func printPort(config *orderedobject.Object[any]) {
	port, found := config.Get("port")
	if !found {
		return
	}
	portInt, ok := port.(float64)
	if !ok {
		fmt.Printf("Port type error: %T\n", port)
		return
	}
	fmt.Printf("Port number: %d\n", int(portInt))
}

func printLookup(config *orderedobject.Object[any], key string) {
	value, found := config.Get(key)
	if !found {
		fmt.Printf("Key '%s' does not exist\n", key)
		return
	}
	fmt.Printf("Found value: %v\n", value)
}

func printServerPort(nestedConfig *orderedobject.Object[any]) {
	server, found := nestedConfig.Get("server")
	if !found {
		fmt.Println("server configuration not found")
		return
	}
	serverObj, ok := server.(*orderedobject.Object[any])
	if !ok {
		fmt.Println("server is not an object type")
		return
	}
	if port, found := serverObj.Get("port"); found {
		fmt.Printf("\nServer port: %v\n", port)
	}
}
