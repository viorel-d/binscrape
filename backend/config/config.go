package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	ServerAddress            string
	MongoDBURI               string
	MongoUsername            string
	MongoPassword            string
	MongoAuthSource          string
	MongoDBName              string
	MongoItemsCollectionName string
)

func init() {
	_ = godotenv.Load()
	ServerAddress = getEnvVar("SERVER_ADDRESS", "0.0.0.0:8000")
	MongoDBURI = getEnvVar("MONGODB_URI", "mongodb://mongo:27017")
	MongoUsername = getEnvVar("MONGODB_USERNAME", "root")
	MongoPassword = getEnvVar("MONGODB_PASSWORD", "pass1337")
	MongoAuthSource = getEnvVar("MONGODB_AUTH_SOURCE", "admin")
	MongoDBName = getEnvVar("MONGODB_NAME", "pasteBinItems")
	MongoItemsCollectionName = getEnvVar("MONGODB_NAME", "raw")
}

func getEnvVar(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = defaultValue
	}
	return val
}
