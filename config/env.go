package config

import (
	"os"

	"github.com/joho/godotenv"
)

func SetupEnvs() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}
