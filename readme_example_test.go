package orderedobject_test

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func ExampleObject_MarshalJSON() {
	person := orderedobject.New[any]().
		Set("name", "Alice").
		Set("age", 30).
		Set("city", "New York")

	data, err := person.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
	// Output:
	// {"name":"Alice","age":30,"city":"New York"}
}
