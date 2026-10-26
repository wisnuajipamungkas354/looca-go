package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
	AdminName     string
	UploadDir     string
	SeedAssetsDir string
	CORSOrigin    string
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func Load() *Config {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPass:        os.Getenv("DB_PASS"),
		DBName:        getEnv("DB_NAME", "looca"),
		JWTSecret:     jwtSecret,
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@looca.id"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
		AdminName:     getEnv("ADMIN_NAME", "Looca Admin"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		SeedAssetsDir: getEnv("SEED_ASSETS_DIR", "./database/seed_assets"),
		CORSOrigin:    getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}
