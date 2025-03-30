[![GoDoc](https://godoc.org/github.com/KarpelesLab/ini?status.svg)](https://godoc.org/github.com/KarpelesLab/ini)

# INI File Parser for Go

A robust INI file handler in Go that supports common INI file features.

## Features

- Simple and intuitive API
- Support for different comment styles (`;` and `#`)
- Handling of quoted values (both single and double quotes)
- Thread-safe option via `IniSafe` wrapper
- Section and key management (add, get, set, delete)
- Case-insensitive keys and sections
- Error reporting with line numbers
- Auto-quoting of values with spaces

## Basic Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/KarpelesLab/ini"
)

func main() {
	// Create a new INI structure
	config := ini.New()

	// Load from a file
	file, _ := os.Open("config.ini")
	defer file.Close()
	config.Load(file)

	// Get values
	if value, ok := config.Get("section", "key"); ok {
		fmt.Printf("Value: %s\n", value)
	}

	// Get with default fallback
	value := config.GetDefault("section", "missing_key", "default_value")
	fmt.Printf("Value with default: %s\n", value)

	// Set values
	config.Set("newsection", "newkey", "newvalue")

	// Check if a section exists
	if config.HasSection("section") {
		fmt.Println("Section exists!")
	}

	// Write to a file
	outFile, _ := os.Create("newconfig.ini")
	defer outFile.Close()
	config.Write(outFile)
}
```

## Thread-Safety

For thread-safe operations, use the `IniSafe` wrapper:

```go
// Create a thread-safe INI structure
config := ini.NewThreadSafe()

// Use the same API methods
config.Set("section", "key", "value")
value, ok := config.Get("section", "key")
```

## License

This library is available under the [LICENSE](./LICENSE) file in the repository.