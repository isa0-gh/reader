package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DBDSN  string
	Debug  bool
	S3     S3Config
	CDN    string
}

type S3Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UsePathStyle    bool
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:  getEnv("PORT", "8080"),
		Debug: getEnv("DEBUG", "false") == "true",
		DBDSN: getEnv("DB_DSN", ""), // We'll build this in a bit or use components
		CDN:   getEnv("CDN_URL_PREFIX", ""),
		S3: S3Config{
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			Region:          getEnv("S3_REGION", "auto"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			Bucket:          getEnv("S3_BUCKET", ""),
			UsePathStyle:    getEnv("S3_USE_PATH_STYLE", "false") == "true",
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
