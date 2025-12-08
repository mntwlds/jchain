# jchain

`jchain` is a lightweight, zero-dependency Go library for parsing and traversing JSON data using a method chaining syntax. It provides a simple and expressive way to access nested JSON values without defining structs.

## Installation

```bash
go get github.com/mntwlds/json
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/mntwlds/json"
)

func main() {
	jsonStr := `
	{
		"users": [
			{
				"id": 1,
				"name": "Alice",
				"details": {
					"age": 30,
					"active": true
				}
			}
		]
	}`

	// Parse the JSON string
	val := jchain.Parse(jsonStr)

	// Check for parse errors
	if err := val.Error(); err != nil {
		log.Fatal(err)
	}

	// Access nested values using method chaining
	name, err := val.Get("users").Index(0).Get("name").String()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name:", name) // Output: Name: Alice

	// Access another nested value
	age, err := val.Get("users").Index(0).Get("details").Get("age").Int()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Age:", age) // Output: Age: 30
}
```

## Features

- **Method Chaining**: Traverse JSON structures fluently (e.g., `.Get().Index().Get()`).
- **Safety**: Safe access to nested values; errors propagate down the chain and can be checked at the end or at any step.
- **Zero Dependencies**: Uses only the Go standard library.
- **Type Conversion**: Easy conversion to native Go types (`String()`, `Int()`, `Bool()`, etc.).

## License

MIT License
