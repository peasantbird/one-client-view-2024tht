package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

// Load the .env file and return DSN
func GetDSN() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    
    return "host=" + os.Getenv("DB_HOST") +
           " user=" + os.Getenv("DB_USER") +
           " dbname=" + os.Getenv("DB_NAME") +
           " password=" + os.Getenv("DB_PASSWORD") +
           " port=" + os.Getenv("DB_PORT") +
           " sslmode=" + os.Getenv("DB_SSLMODE") +
           " TimeZone=" + os.Getenv("DB_TIMEZONE")
}
