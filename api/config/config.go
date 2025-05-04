package config

import (
	"os"

	"github.com/joho/godotenv"

	"belchi/logger"
)

type Config struct {
	HTTPPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBName    string
	DBPwd     string
	DBSSLMode string

	UploadPath string
}

var AppConfig Config

func LoadAppConfig() error {

	err := godotenv.Load()
	if err != nil {
		logger.Log("DEBUG", "Error on read .env file", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
	}

	AppConfig = Config{
		HTTPPort: getEnvStr("HTTPPort", "8080"),

		DBHost:    getEnvStr("DBHost", "localhost"),
		DBPort:    getEnvStr("DBPort", "5432"),
		DBUser:    getEnvStr("DBUser", "postgres"),
		DBName:    getEnvStr("DBName", "belchi"),
		DBPwd:     getEnvStr("DBPwd", "postgres"),
		DBSSLMode: getEnvStr("DBSSLMode", "disable"),

		UploadPath: getEnvStr("UploadPath", "uploads/"),
	}

	logger.Log("ENV", "Defined environment variables", []logger.LogDetail{
		{Key: "HTTPPort", Value: AppConfig.HTTPPort},
		{Key: "DBHost", Value: AppConfig.DBHost},
		{Key: "DBPort", Value: AppConfig.DBPort},
		{Key: "DBUser", Value: AppConfig.DBUser},
		{Key: "DBName", Value: AppConfig.DBName},
		{Key: "DBPwd", Value: AppConfig.DBPwd},
		{Key: "DBSSLMode", Value: AppConfig.DBSSLMode},
		{Key: "UploadPath", Value: AppConfig.UploadPath},
	})
	return nil
}

func getEnvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
