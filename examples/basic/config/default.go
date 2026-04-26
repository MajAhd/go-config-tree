package config

// Default returns the baseline configuration for the application.
// These are the hardcoded values that will be used if no environment variable overrides them.
func Default() *Schema {
	return &Schema{
		App: AppConfig{
			Name: "mycollapp", // Hardcoded default as requested
			Port: 8080,
			Host: "localhost",
		},
		DB: DBConfig{
			Host: "127.0.0.1", // Default to local postgres
			User: "postgres",
		},
		// Secrets are left empty here intentionally because they shouldn't have default values
		Secret: SecretConfig{},
	}
}
