package goconfigtree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Schema simulates the user's config schema
type Schema struct {
	App struct {
		Name string `env:"APP_NAME,required"`
		Port int    `env:"APP_PORT,required"`
	}
	DB struct {
		Host     string `env:"DB_HOST,required"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD,required"`
	}
}

// Default simulates the default config fallback
func Default() *Schema {
	return &Schema{
		App: struct {
			Name string `env:"APP_NAME,required"`
			Port int    `env:"APP_PORT,required"`
		}{
			Name: "mycollapp",
			Port: 8080,
		},
		DB: struct {
			Host     string `env:"DB_HOST,required"`
			User     string `env:"DB_USER"`
			Password string `env:"DB_PASSWORD,required"`
		}{
			Host:     "127.0.0.1",
			User:     "postgres",
			Password: "default-local-password",
		},
	}
}

func TestScenario1_EverythingWorksLocalAndProd(t *testing.T) {
	// --- Local Setup ---
	tmpDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	envContent := []byte(`
APP_PORT=3000
DB_PASSWORD="overridden_local_password"
`)
	envPath := filepath.Join(tmpDir, ".env.local")
	if err := os.WriteFile(envPath, envContent, 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Local Load (Base Default() + .env.local)
	localCfg := Default()
	if err := Load(localCfg, envPath); err != nil {
		t.Fatalf("Failed to load local config: %v", err)
	}

	if localCfg.App.Name != "mycollapp" {
		t.Errorf("Expected mycollapp (from default), got %v", localCfg.App.Name)
	}
	if localCfg.App.Port != 3000 {
		t.Errorf("Expected 3000 (from .env.local), got %v", localCfg.App.Port)
	}
	if localCfg.DB.Password != "overridden_local_password" {
		t.Errorf("Expected overridden_local_password, got %v", localCfg.DB.Password)
	}

	// --- Prod Setup ---
	os.Setenv("APP_NAME", "my-prod-app")
	os.Setenv("APP_PORT", "80")
	os.Setenv("DB_HOST", "postgres-prod.cluster.local")
	os.Setenv("DB_USER", "prod-user")
	os.Setenv("DB_PASSWORD", "super-secret-prod-password")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_PORT")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")

	// 2. Prod Load (Blank &Schema{} + os.Getenv overrides)
	prodCfg := &Schema{}
	if err := Load(prodCfg); err != nil {
		t.Fatalf("Failed to load prod config: %v", err)
	}

	if prodCfg.App.Name != "my-prod-app" {
		t.Errorf("Expected my-prod-app, got %v", prodCfg.App.Name)
	}
	if prodCfg.DB.Password != "super-secret-prod-password" {
		t.Errorf("Expected prod password, got %v", prodCfg.DB.Password)
	}
}

func TestScenario2_ProdMissingRequiredPassword(t *testing.T) {
	localCfg := Default()
	if err := Load(localCfg); err != nil {
		t.Fatalf("Local failed unexpectedly! It should have used defaults: %v", err)
	}

	// --- Prod Setup ---
	os.Setenv("APP_NAME", "my-prod-app")
	os.Setenv("APP_PORT", "80")
	os.Setenv("DB_HOST", "postgres-prod.cluster.local")
	// DB_USER and DB_PASSWORD intentionally missing
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_PORT")
	defer os.Unsetenv("DB_HOST")

	// 2. Prod Load (Blank &Schema{} + os.Getenv overrides)
	prodCfg := &Schema{}
	err := Load(prodCfg)

	if err == nil {
		t.Fatalf("Expected prod config to fail because DB_PASSWORD is not set and no default exists, but it succeeded!")
	}

	// Verify it failed specifically because of the required DB_PASSWORD
	if !strings.Contains(err.Error(), "required environment variable DB_PASSWORD is missing") {
		t.Errorf("Expected missing DB_PASSWORD error, got: %v", err)
	}
}

func Test_wrongSchema(t *testing.T) {
	os.Setenv("APP_NAME", "my-prod-app")
	os.Setenv("APP_PORT", "NOT_NUMBER")
	os.Setenv("DB_HOST", "postgres-prod.cluster.local")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")

	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_PORT")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")

	prodCfg := &Schema{}
	err := Load(prodCfg)
	if err == nil {
		t.Fatalf("Expected an error due to invalid APP_PORT, but got none")
	}
	if !strings.Contains(err.Error(), "invalid syntax") || !strings.Contains(err.Error(), "APP_PORT") {
		t.Errorf("Expected invalid syntax error for APP_PORT, got: %v", err)
	}
}
