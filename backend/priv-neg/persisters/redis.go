package persisters

import (
	"fmt"
	"log"
	"os"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/garyburd/redigo/redis"
)

// NewRedisDB - Returns a new redis connection.
func NewRedisDB(logger *log.Logger) redis.Conn {

	if !utils.WaitForService(fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), "6379"), logger) {
		panic("Could not find Redis.")
	}

	c, err := redis.DialURL(fmt.Sprintf("redis://%s:6379/0", os.Getenv("REDIS_HOST")))

	if err != nil {
		panic("Failed to connect to Redis.")
	}

	return c
}
