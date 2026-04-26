package config

// Schema defines the structure of the entire application configuration.
// We use `env` tags so that go-config-tree knows which environment variables map to which fields.
// You can now add ",required" to strictly enforce presence in production!
type Schema struct {
	App    AppConfig
	DB     DBConfig
	Secret SecretConfig
}

type AppConfig struct {
	Name string `env:"APP_NAME,required"`
	Port int    `env:"APP_PORT,required"`
	Host string `env:"APP_HOST"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST,required"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD,required"`
}

type SecretConfig struct {
	GoogleSec    string `env:"GOOGLE_SEC,required"`
	MicrosoftSec string `env:"MICROSOFT_SEC,required"`
}
