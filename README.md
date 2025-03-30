[![Go Report Card](https://goreportcard.com/badge/github.com/KarpelesLab/ini?style=flat-square)](https://goreportcard.com/report/github.com/KarpelesLab/ini)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/KarpelesLab/ini)](https://pkg.go.dev/github.com/KarpelesLab/ini)
[![Tags](https://img.shields.io/github/tag/KarpelesLab/ini.svg?style=flat-square)](https://github.com/KarpelesLab/ini/tags)

# INI File Parser for Go

A robust INI file handler in Go that supports common INI file features.

## Features

- Simple and intuitive API
- Implements standard Go interfaces (`io.ReaderFrom` and `io.WriterTo`)
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

	// Load from a file (two equivalent ways)
	file, _ := os.Open("config.ini")
	defer file.Close()
	
	// Using io.ReaderFrom interface (preferred)
	config.ReadFrom(file)
	
	// Or using the legacy method
	// config.Load(file)

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

	// Write to a file (two equivalent ways)
	outFile, _ := os.Create("newconfig.ini")
	defer outFile.Close()
	
	// Using io.WriterTo interface (preferred)
	config.WriteTo(outFile)
	
	// Or using the legacy method
	// config.Write(outFile)
}
```

## Standard Go Interfaces

The library implements the standard Go interfaces:

- `io.ReaderFrom`: Parse INI files with `ReadFrom(r io.Reader) (n int64, err error)`
- `io.WriterTo`: Write INI files with `WriteTo(w io.Writer) (n int64, err error)`

These interfaces make the library more idiomatic and composable with other Go libraries.

## Thread-Safety

For thread-safe operations, use the `IniSafe` wrapper:

```go
// Create a thread-safe INI structure
config := ini.NewThreadSafe()

// Use the same API methods
config.Set("section", "key", "value")
value, ok := config.Get("section", "key")

// Use standard interfaces
file, _ := os.Open("config.ini")
config.ReadFrom(file)

outFile, _ := os.Create("output.ini")
config.WriteTo(outFile)
```

## Section Management

```go
// Check if section exists
if config.HasSection("mysection") {
    // Get all section names
    sections := config.Sections()
    
    // Get all keys in a section
    keys := config.Keys("mysection")
}
```

## License

This library is available under the [LICENSE](./LICENSE) file in the repository.
