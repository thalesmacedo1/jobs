package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Neo4j Configuration
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string

	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

var Settings Config

func LoadConfig(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}

	// Carregar as configurações do Neo4j
	Settings.Neo4jURI = strings.TrimSpace(os.Getenv("NEO4J_URI"))
	Settings.Neo4jUser = strings.TrimSpace(os.Getenv("NEO4J_USER"))
	Settings.Neo4jPassword = strings.TrimSpace(os.Getenv("NEO4J_PASSWORD"))

	// Validar configurações do Neo4j
	if Settings.Neo4jURI == "" || Settings.Neo4jUser == "" || Settings.Neo4jPassword == "" {
		log.Fatal("Missing required Neo4j environment variables")
	}

	// Carregar as configurações do Redis
	Settings.RedisHost = strings.TrimSpace(os.Getenv("REDIS_HOST"))
	Settings.RedisPort = strings.TrimSpace(os.Getenv("REDIS_PORT"))
	Settings.RedisPassword = strings.TrimSpace(os.Getenv("REDIS_PASSWORD"))
	redisDBStr := strings.TrimSpace(os.Getenv("REDIS_DB"))

	// Validar configurações do Redis
	if Settings.RedisHost == "" || Settings.RedisPort == "" {
		log.Fatal("Missing required Redis environment variables (REDIS_HOST, REDIS_PORT)")
	}

	// Converter RedisDB para int
	if redisDBStr == "" {
		Settings.RedisDB = 0
	} else {
		db, err := strconv.Atoi(redisDBStr)
		if err != nil {
			log.Fatalf("Invalid Redis DB value: %v", err)
		}
		Settings.RedisDB = db
	}

	return nil
}
