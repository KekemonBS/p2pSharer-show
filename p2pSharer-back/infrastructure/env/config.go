package env

import (
	"fmt"
	"os"
)

// Config struct holds data parsed from environment
type Config struct {
	//Add some field if i need to
	RedisURI string
}

// NewConfig parses env to struct
func NewConfig() (*Config, error) {
	// postgresql://[userspec@][hostspec][/dbname][?paramspec]
	redisURI, ok := os.LookupEnv("REDISURI")
	if !ok {
		return nil, fmt.Errorf("no REDISURI env variable")
	}

	return &Config{
		RedisURI: redisURI,
	}, nil
}
