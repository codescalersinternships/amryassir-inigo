# INI Parser Package

A Go package for parsing, manipulating, and serializing INI files.

## Overview

This package provides functionalities to parse INI files into a structured format, manipulate sections and keys, and serialize them back into INI format. It supports reading from files or strings, modifying section data, and saving to new INI file.

## Features

- **Section Management**: ’anage sections and their keys.
- **Serialization**: Convert INI data to a string format.
- **File Operations**: Load from and save to INI files.

## Installation

To install the package, use `go get`:

```bash
go get github.com/codescalersinternships/amryassir-inigo
```
## Usage

``` go
# Load data from a string
ini, err := pkg.LoadFromString(sample String)
if err != nil {
     log.Fatalf("Failed to load sample", err)
}

# Load data from a file
ini, err := pkg.LoadFromFile(fileName String)
if err != nil {
     log.Fatalf("Failed to load from file", err)
}

# Get section names
sections := ini.GetSectionNames()

# Get sections map
var val map[string]map[string]string = ini.GetSections()

# Get a value
value, err := ini.Get("sectionName", "key")
if err != nil {
	log.Fatalf("Failed to get value: %v", err)
}

# Set a new value
ini.Set("sectionName", "key", "value")

# Convert INI data to string
ini := ini.ToString()

# Save INI data to a file
err = ini.SaveToFile("NewFile.ini")
if err != nil {
	log.Fatalf("Failed to save INI data to file: %v", err)
}
```

# Testing
To run the tests for this package, use the following command:

```bash
go test ./...
```
