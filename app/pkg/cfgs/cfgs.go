package cfgs

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Configs struct {
	DB_CONNECTION_URI    string
	REDIS_CONNECTION_URI string
	LLM_URI              string
}

func LoadConfigs() Configs {
	defaultConfig := Configs{}
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	configValue := reflect.ValueOf(&defaultConfig).Elem()
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envKey := "GJ_" + field.Name
		if value := os.Getenv(envKey); value != "" {
			configValue.Field(i).SetString(value)
		} else {
			log.Fatalf("Error: Required environment variable %s not set", envKey)
		}
	}
	return defaultConfig
}
