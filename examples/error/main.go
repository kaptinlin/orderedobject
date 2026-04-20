package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Error Handling Example ===")

	invalidJSON := `{"name": "test", "port": "not_a_number"}`
	config, err := orderedobject.FromJSON[any]([]byte(invalidJSON))
	if err != nil {
		fmt.Printf("\nParse error: %v\n", err)
		config = orderedobject.NewObject[any]()
	} else {
		fmt.Println("\nSuccessfully parsed JSON")
	}

	if port, found := config.Get("port"); found {
		if portInt, ok := port.(float64); ok {
			fmt.Printf("Port number: %d\n", int(portInt))
		} else {
			fmt.Printf("Port type error: %T\n", port)
		}
	}

	if value, found := config.Get("nonexistent"); found {
		fmt.Printf("Found value: %v\n", value)
	} else {
		fmt.Println("Key 'nonexistent' does not exist")
	}

	nestedConfig := orderedobject.NewObject[any]().
		Set("server", orderedobject.NewObject[any]().
			Set("port", 8080))

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
