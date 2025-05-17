package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	fmt.Println("=== Type Safety Example ===")

	// Define types
	type User struct {
		Name  string
		Age   int
		Email string
	}

	type Settings struct {
		Theme  string
		Active bool
	}

	// Create type-safe objects
	users := orderedobject.NewObject[User]().
		Set("user1", User{Name: "Alice", Age: 30, Email: "alice@example.com"}).
		Set("user2", User{Name: "Bob", Age: 25, Email: "bob@example.com"})

	settings := orderedobject.NewObject[Settings]().
		Set("default", Settings{Theme: "light", Active: true}).
		Set("custom", Settings{Theme: "dark", Active: false})

	// Access typed values
	if user, found := users.Get("user1"); found {
		fmt.Printf("\nUser: %s, Age: %d\n", user.Name, user.Age)
	}

	fmt.Println("\nAll settings:")
	settings.ForEach(func(key string, value Settings) {
		fmt.Printf("  %s: Theme=%s, Active=%v\n", key, value.Theme, value.Active)
	})
}
