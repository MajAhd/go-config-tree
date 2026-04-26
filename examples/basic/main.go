package main

import (
	"fmt"
	"os"

	"github.com/MajAhd/go-config-tree/examples/basic/config"
)

func main() {
	// Let's pretend Kubernetes sets APP_STAGE for us.
	// Normally you'd read os.Getenv("APP_STAGE"), but if it's empty, default to "local"
	stage := os.Getenv("APP_STAGE")
	if stage == "" {
		stage = "local" // fallback to local development
	}

	var cfg *config.Schema

	// The application can now easily decide what stage needed to load!
	if stage == "prod" {
		fmt.Println("--> Running in PRODUCTION stage")

		// In a real prod environment, these would be injected by Kubernetes Secrets!
		// For this example to work, we'll manually set them in the OS environment here:
		os.Setenv("APP_NAME", "my-prod-app")
		os.Setenv("APP_PORT", "80")
		os.Setenv("GOOGLE_SEC", "prod-google-k8s-secret")
		os.Setenv("MICROSOFT_SEC", "prod-microsoft-k8s-secret")
		os.Setenv("DB_HOST", "postgres-prod.cluster.local")

		cfg = config.Prod()

	} else {
		fmt.Println("--> Running in LOCAL stage")

		// Local stage reads from .env.local!
		cfg = config.Local()
	}

	fmt.Printf("\n--- Loaded Configuration ---\n")
	fmt.Printf("App Name:  %s\n", cfg.App.Name)
	fmt.Printf("App Port:  %d\n", cfg.App.Port)
	fmt.Printf("App Host:  %s\n", cfg.App.Host)
	fmt.Printf("DB Host:   %s\n", cfg.DB.Host)
	fmt.Printf("DB User:   %s\n", cfg.DB.User)
	fmt.Printf("Google:    %s\n", cfg.Secret.GoogleSec)
	fmt.Printf("Microsoft: %s\n", cfg.Secret.MicrosoftSec)
}
