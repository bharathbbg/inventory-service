package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr  string
	GRPCAddr  string
	MongoDB   MongoDBConfig
	Redis     RedisConfig
	Services  ServicesConfig
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	Host string
	Port int
}

type ServicesConfig struct {
	OrderService ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port int
}

func Load() (*Config, error) {
	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	orderPort, _ := strconv.Atoi(getEnv("ORDER_SERVICE_PORT", "50051"))

	return &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8082"),
		GRPCAddr: getEnv("GRPC_ADDR", ":50053"),
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DB", "inventory_db"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: redisPort,
		},
		Services: ServicesConfig{
			OrderService: ServiceConfig{
				Host: getEnv("ORDER_SERVICE_HOST", "localhost"),
				Port: orderPort,
			},
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}