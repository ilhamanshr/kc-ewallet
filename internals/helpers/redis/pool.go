package redis

import (
	"fmt"
	log_color "kc-ewallet/internals/helpers/color"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

type Configuration struct {
	Password string
	Host     string
	Db       int
	Port     int
}

func PanicOnError(err error, msg string) {
	if err != nil {
		log_color.PrintRedf("%v: %v", msg, err)
		panic(fmt.Sprintf("%v: %v", msg, err))
	}
}

func newRedisPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func InitRedisPool(conf *Configuration) {
	serverUrl := conf.Host + ":" + strconv.Itoa(conf.Port)
	RedisPool = newRedisPool(serverUrl, conf.Password, conf.Db)
	c := RedisPool.Get()
	defer c.Close()

	pong, err := redis.String(c.Do("PING"))
	PanicOnError(err, "Cannot ping Redis")
	log_color.PrintGreenf("Redis Ping: %s", pong)
}
