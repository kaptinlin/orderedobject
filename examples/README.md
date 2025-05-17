# Example Usage

This directory contains a comprehensive set of examples demonstrating how to use the OrderedObject library. The examples cover both basic and advanced features, with a focus on practical use cases.

## Running the Examples

Each example is a standalone executable program that can be run directly:

```bash
# Run basic operations example
go run basic/main.go

# Run JSON operations example
go run json/main.go

# Run Map operations example
go run map/main.go

# Run type safety example
go run type/main.go

# Run nested structures example
go run nested/main.go

# Run array operations example
go run array/main.go

# Run error handling example
go run error/main.go
```

## Example Structure

The examples are organized into the following functional categories:

### Core Features

1. **Basic Operations** (`basic/main.go`)
   - Creating and populating objects
   - Getting and updating values
   - Deleting keys
   - Iterating through entries
   - Cloning objects
   - Configuration management example

2. **JSON Operations** (`json/main.go`)
   - Converting objects to JSON (three methods)
   - Parsing JSON into objects
   - User information handling
   - JSON serialization best practices

3. **Map Operations** (`map/main.go`)
   - Creating objects from maps
   - Converting objects to maps
   - Application settings management
   - Order preservation demonstration

4. **Type Safety** (`type/main.go`)
   - Using concrete types with generics
   - Type-safe value retrieval
   - Multiple type handling
   - Struct integration

### Advanced Features

1. **Nested Structures** (`nested/main.go`)
   - Complex object hierarchies
   - Application configuration management
   - Nested value access
   - JSON serialization of nested objects

2. **Array Operations** (`array/main.go`)
   - Arrays of objects
   - Order management
   - Array element access
   - Settings item processing example

3. **Error Handling** (`error/main.go`)
   - JSON parsing errors
   - Type assertion
   - Safe value access
   - Error recovery patterns

## Key Concepts Demonstrated

The examples illustrate several important concepts:

1. **Order Preservation**
   - Key insertion order is maintained
   - JSON serialization preserves order
   - Map conversion order handling

2. **Type Safety**
   - Generic type parameters
   - Type assertions
   - Struct integration
   - Type-safe value access

3. **Method Chaining**
   - Fluent interface pattern
   - Builder-style API
   - Immutable operations

4. **Error Handling**
   - JSON parsing
   - Type conversion
   - Safe value access
   - Error recovery

## Best Practices

The examples demonstrate several best practices:

1. **Configuration Management**
   - Using ordered objects for config
   - Environment-specific settings
   - Configuration cloning

2. **Data Processing**
   - Type-safe data handling
   - Nested structure management
   - Array operations
   - Error handling

3. **JSON Handling**
   - Multiple serialization methods
   - Order preservation
   - Nested object handling
   - Error recovery

## Expected Output

Each example includes clear output demonstrating:
- Object structure
- JSON serialization
- Order preservation
- Type safety
- Error handling

The output is formatted for clarity and includes comments explaining the results.

## Notes

- All examples use the latest features of the library
- Error handling is demonstrated where appropriate
- Type safety is emphasized throughout
- Real-world use cases are provided
- Code is thoroughly commented
- Output is clearly formatted

## Contributing Examples

If you'd like to contribute additional examples:
1. Follow the existing code style
2. Include clear comments
3. Demonstrate a specific feature or use case
4. Add appropriate error handling
5. Update this README if necessary