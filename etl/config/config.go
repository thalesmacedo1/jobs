package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
}

var Settings Config

func LoadConfig(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}

	Settings = Config{
		Neo4jURI:      os.Getenv("NEO4J_URI"),
		Neo4jUser:     os.Getenv("NEO4J_USER"),
		Neo4jPassword: os.Getenv("NEO4J_PASSWORD"),
	}

	if Settings.Neo4jURI == "" || Settings.Neo4jUser == "" || Settings.Neo4jPassword == "" {
		log.Fatal("Missing required environment variables")
	}

	return nil
}
