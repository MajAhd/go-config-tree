package config

import (
	"log"

	"github.com/MajAhd/go-config-tree"
)

// Prod loads the production configuration overrides.
// It bypasses the hardcoded defaults and relies PURELY on the system environment variables
// (like those injected by Kubernetes Secrets/ConfigMaps).
func Prod() *Schema {
	// Instead of Default(), we use a blank struct.
	// This forces all fields marked with "required" to come from the environment!
	cfg := &Schema{}

	if err := goconfigtree.Load(cfg); err != nil {
		// Using log.Fatalf here will crash the app if prod variables are missing!
		log.Fatalf("CRITICAL: Failed to load prod config: %v\n", err)
	}

	return cfg
}
