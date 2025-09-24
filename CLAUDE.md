# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository contains both Postman collections for API testing and a custom Go-based Postman Collection Tester CLI tool.

## Structure

- `postman/` - Contains original Postman collection files for Korean financial services APIs
- `test-collection.json` - Test collection using JSONPlaceholder API for demonstration
- Go source files for the Postman Collection Tester:
  - `main.go` - CLI interface and main program logic
  - `postman.go` - Postman collection data structures
  - `runner.go` - HTTP request execution engine
  - `reporter.go` - Test result reporting (text, JSON, HTML)

## Development Commands

### Build the tester
```bash
go build -o postman-tester
```

### Cross-platform builds
```bash
# Linux/macOS
./build.sh

# Windows
build.bat
```

### Run the tester
```bash
# Test single file
./postman-tester -file test-collection.json

# Test directory
./postman-tester -dir ./postman

# Generate HTML report
./postman-tester -file test-collection.json -output report.html -format html

# Parallel execution
./postman-tester -dir ./postman -parallel 3 -verbose
```

## Key Features

- **Cross-platform**: Single binary runs on Windows, macOS, Linux
- **Multiple input modes**: Single file (-file) or directory (-dir)
- **Output formats**: Text, JSON, HTML reports
- **Parallel execution**: Run multiple collections concurrently
- **Comprehensive reporting**: Detailed test results with timing and error information

## CLI Options

- `-file`: Single Postman collection file
- `-dir`: Directory containing collection files (default: ./postman)
- `-output`: Save results to file
- `-format`: Output format (text, json, html)
- `-parallel`: Number of concurrent collections (default: 1)
- `-verbose`: Detailed output
- `-timeout`: Request timeout in seconds (default: 30)

## Testing

Use `test-collection.json` for testing the CLI tool - it uses JSONPlaceholder API which is publicly accessible and reliable for testing purposes.