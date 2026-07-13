package orderedobject_test

import (
	"encoding/json"
	"fmt"

	"github.com/kaptinlin/orderedobject"
)

func ExampleNew() {
	person := orderedobject.New[any]().
		Set("name", "Alice").
		Set("age", 30).
		Set("city", "New York")

	data, err := json.Marshal(person)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
	// Output:
	// {"name":"Alice","age":30,"city":"New York"}
}

func ExampleObject_All() {
	obj := orderedobject.New[int]().Set("first", 1).Set("second", 2)
	for key, value := range obj.All() {
		fmt.Printf("%s=%d\n", key, value)
		break
	}
	// Output:
	// first=1
}

func ExampleFromJSON() {
	obj, err := orderedobject.FromJSON[any]([]byte(`{"second":2,"first":1}`))
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
	// Output:
	// {"second":2,"first":1}
}

func ExampleObject_ToUnorderedMap() {
	obj := orderedobject.New[int]().Set("a", 1).Set("b", 2)
	plain := obj.ToUnorderedMap()
	fmt.Println(plain["a"], len(plain))
	// Output:
	// 1 2
}
