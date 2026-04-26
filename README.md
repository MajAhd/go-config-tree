# go-config-tree

[![Go Report Card](https://goreportcard.com/badge/github.com/MajAhd/go-config-tree)](https://goreportcard.com/report/github.com/MajAhd/go-config-tree)
[![Go Reference](https://pkg.go.dev/badge/github.com/MajAhd/go-config-tree.svg)](https://pkg.go.dev/github.com/MajAhd/go-config-tree)
[![Build Status](https://github.com/MajAhd/go-config-tree/actions/workflows/go.yml/badge.svg)](https://github.com/MajAhd/go-config-tree/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

`go-config-tree` is a lightweight, dependency-free Go package designed to manage environment variables and configuration files for any Go project.

It provides a clean workflow for configuration management:
1. Define your own schema with normal Go structs.
2. Initialize the struct with your **hardcoded default values**.
3. Use `.env` files and `os.Getenv()` environment variables to dynamically override defaults.
4. Strictly enforce required variables in production.

## Features
- **Zero Dependencies**: A highly lightweight package with no external dependencies.
- **`.env` File Support**: Includes a built-in, lightweight `.env` parser.
- **Required Fields**: Enforce required environment variables using the `required` tag option.
- **Hierarchical Structs**: Map values seamlessly to nested Go structs using `env` struct tags.
- **Type Conversions**: Automatically converts environment variables into your struct field types (string, int, bool, float).

## Installation

```bash
go get github.com/MajAhd/go-config-tree
```

## Quick Start

### 1. Define your Configuration Schema

Create your Go structs exactly the way you want your config tree structured. Use `env` tags for environment variable bindings, and append `,required` to fields that must be present.

```go
package main

import (
	"fmt"
	"log"
	"os"

	configtree "github.com/MajAhd/go-config-tree"
)

type Config struct {
	App struct {
		Name  string `env:"APP_NAME,required"`
		Port  int    `env:"APP_PORT,required"`
		Debug bool   `env:"APP_DEBUG"`
	}
	DB struct {
		Host     string `env:"DB_HOST,required"`
		Username string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD,required"`
	}
}

func main() {
	// 1. Set your hardcoded default values
	cfg := Config{}
	cfg.App.Name = "my-default-app"
	cfg.App.Port = 3000
	cfg.App.Debug = false
	cfg.DB.Host = "localhost"
	cfg.DB.Username = "root"

	// (Simulate an OS-level environment variable like you would set in docker/production)
	os.Setenv("DB_PASSWORD", "supersecret123")

	// 2. Load the configuration
	// This will read the `.env` file (if it exists) and merge all environment 
	// variables over the defaults inside `cfg`.
	if err := configtree.Load(&cfg, ".env"); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Loaded Config:\n")
	fmt.Printf("App Name: %s\n", cfg.App.Name)
	fmt.Printf("App Port: %d\n", cfg.App.Port)
	fmt.Printf("DB Host: %s\n", cfg.DB.Host)
	fmt.Printf("DB Password: %s\n", cfg.DB.Password)
}
```

### 2. Create an optional `.env` file

In the root of your project, you can drop a `.env` file (or `.env.local`) to easily switch environments:

```env
# .env
APP_PORT=8080
APP_DEBUG=true
DB_HOST="127.0.0.1"
DB_USER=admin
```

### 3. Run the application

When `Load` is called, `go-config-tree` performs the following merges over your pre-populated defaults:
1. Applies variables found in the `.env` file (if provided).
2. Applies system environment variables (like those from `os.Setenv` or Docker).

Output:
```text
Loaded Config:
App Name: my-default-app
App Port: 8080
DB Host: 127.0.0.1
DB Password: supersecret123
```

---

## Strict Production Validation (`required` tags)

If you omit the hardcoded defaults and pass a blank struct into `Load()`, `go-config-tree` will strictly enforce that all fields marked with `,required` must exist in the environment variables.

If an environment variable is missing, `Load()` will return an error:
`required environment variable DB_PASSWORD is missing and no default value was provided`

This is extremely powerful for ensuring production environments don't start up with missing configurations.

---

## Best Practices Architecture

We highly recommend separating your configurations into a dedicated `config/` directory for scalable applications:

* `config/configSchema.go` - Holds your configuration struct.
* `config/default.go` - Contains your baseline, hardcoded values and returns a populated `*Schema`.
* `config/local.go` - Uses the `Default()` baseline, and loads `.env.local` over it for local development overrides.
* `config/prod.go` - Bypasses `Default()` and initializes an empty `&Schema{}`. It then calls `Load()` without passing an `.env` file, meaning it strictly pulls from `os.Getenv()` (like from Kubernetes Secrets) and strictly validates all `required` tags!

See the [examples/basic](examples/basic) folder for a fully functional implementation of this best-practice architecture.
