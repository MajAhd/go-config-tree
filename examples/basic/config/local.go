package config

import (
	"log"

	"github.com/MajAhd/go-config-tree"
)

// Local loads the local configuration overrides.
// It starts with Default(), and then uses go-config-tree to parse `.env.local`.
func Local() *Schema {
	cfg := Default()

	// By passing ".env.local", go-config-tree will load those values into the OS environment
	// and map them over our Default `cfg`.
	// The secrets will be read from the .env.local file.
	if err := goconfigtree.Load(cfg, ".env.local"); err != nil {
		log.Printf("Warning: Failed to load local config overrides: %v\n", err)
	}

	return cfg
}
