// Package goconfigtree is a lightweight, dependency-free Go package designed to manage
// environment variables and configuration files for any Go project.
//
// It provides a clean workflow for configuration management by allowing you to define
// your own schema using normal Go structs, initialize them with hardcoded default values,
// and use .env files or system environment variables to dynamically override those defaults.
//
// It also provides strict validation through `required` tags, ensuring production
// environments don't start up with missing configurations.
package goconfigtree

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Load reads optional .env files, loads them into the OS environment (if not already set),
// and then maps all available environment variables to the provided struct pointer v
// based on the `env` struct tags.
// Any existing values in v act as default values.
//
// If a field has a `,required` option in its `env` tag and the environment variable
// is missing while the struct field is at its zero-value, Load will return an error.
func Load(v interface{}, envFiles ...string) error {
	for _, filename := range envFiles {
		if err := loadEnvFile(filename); err != nil {
			// If file doesn't exist, we skip it.
			if !os.IsNotExist(err) {
				return fmt.Errorf("error reading %s: %w", filename, err)
			}
		}
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a non-nil pointer to a struct")
	}

	if err := loadEnvVars(val.Elem()); err != nil {
		return fmt.Errorf("error parsing environment variables: %w", err)
	}

	return nil
}

// loadEnvFile is a lightweight .env file parser.
// It reads KEY=VALUE pairs and sets them in the OS environment
// ONLY if they are not already set.
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}

	return scanner.Err()
}

func loadEnvVars(val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		if fieldVal.Kind() == reflect.Struct {
			if err := loadEnvVars(fieldVal); err != nil {
				return err
			}
			continue
		}

		tagValue := field.Tag.Get("env")
		if tagValue == "" {
			continue
		}

		parts := strings.Split(tagValue, ",")
		envKey := strings.TrimSpace(parts[0])
		isRequired := false
		for _, p := range parts[1:] {
			if strings.TrimSpace(p) == "required" {
				isRequired = true
			}
		}

		if envVal, exists := os.LookupEnv(envKey); exists {
			if err := setField(fieldVal, envVal); err != nil {
				return fmt.Errorf("failed to set field %s from env %s: %w", field.Name, envKey, err)
			}
		} else if isRequired && fieldVal.IsZero() {
			return fmt.Errorf("required environment variable %s is missing and no default value was provided", envKey)
		}
	}
	return nil
}

func setField(field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}
