package config

import "os"

type Config struct {
	MongoURI   string
	MongoDB    string
	SiriDBHost string
	SiriDBPort string
	SiriDBUser string
	SiriDBPass string
	MQTTBroker string
	APIPort    string
	WebPort    string
}

func Load() *Config {
	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:    getEnv("MONGO_DB", "home_dashboard"),
		SiriDBHost: getEnv("SIRIDB_HOST", "localhost"),
		SiriDBPort: getEnv("SIRIDB_PORT", "9000"),
		SiriDBUser: getEnv("SIRIDB_USER", "iris"),
		SiriDBPass: getEnv("SIRIDB_PASS", "siri"),
		MQTTBroker: getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		APIPort:    getEnv("API_PORT", "8080"),
		WebPort:    getEnv("WEB_PORT", "8090"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
