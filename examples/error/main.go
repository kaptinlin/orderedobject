package main

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Error Handling Example ===")

	// 1. Handle invalid JSON input
	invalidJSON := `{"name": "test", "port": "not_a_number"}`
	config := orderedobject.NewObject[any]()
	err := json.Unmarshal([]byte(invalidJSON), config)
	if err != nil {
		fmt.Printf("\nParse error: %v\n", err)
	} else {
		fmt.Println("\nSuccessfully parsed JSON")
	}

	// 2. Type assertion and error handling
	if config != nil {
		if port, found := config.Get("port"); found {
			// Try to convert port to integer
			if portInt, ok := port.(float64); ok {
				fmt.Printf("Port number: %d\n", int(portInt))
			} else {
				fmt.Printf("Port type error: %T\n", port)
			}
		}

		// 3. Safely access potentially non-existent values
		if value, found := config.Get("nonexistent"); found {
			fmt.Printf("Found value: %v\n", value)
		} else {
			fmt.Println("Key 'nonexistent' does not exist")
		}
	}

	// 4. Handle nested object errors
	nestedConfig := orderedobject.NewObject[any]().
		Set("server", orderedobject.NewObject[any]().
			Set("port", 8080))

	if server, found := nestedConfig.Get("server"); found {
		if serverObj, ok := server.(*orderedobject.Object[any]); ok {
			if port, found := serverObj.Get("port"); found {
				fmt.Printf("\nServer port: %v\n", port)
			}
		} else {
			fmt.Println("server is not an object type")
		}
	} else {
		fmt.Println("server configuration not found")
	}
}
