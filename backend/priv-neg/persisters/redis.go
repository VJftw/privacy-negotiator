package persisters

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/garyburd/redigo/redis"
)

// NewRedisDB - Returns a new redis connection.
func NewRedisDB(logger *log.Logger) *redis.Pool {

	if !utils.WaitForService(fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), "6379"), logger) {
		panic("Could not find Redis.")
	}

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        dialRedis,
	}
}

func dialRedis() (redis.Conn, error) {
	return redis.DialURL(
		fmt.Sprintf("redis://%s:6379/0", os.Getenv("REDIS_HOST")))
}
