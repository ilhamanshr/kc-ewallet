package configurations

import (
	"os"
	"strconv"
	"strings"
)

type redisConfiguration struct {
	addr string
	db   int
}

type IRedisConfiguration interface {
	GetAddr() string
	GetDb() int
	GetHost() string
	GetPort() int
	SetRedisConfiguration(addr string, db int)
}

func NewRedisConfiguration() *redisConfiguration {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 1
	}
	return &redisConfiguration{
		addr: os.Getenv("REDIS_URL"),
		db:   db,
	}
}

func (c *redisConfiguration) SetRedisConfiguration(addr string, db int) {
	c.addr = addr
	c.db = db
}

func (c *redisConfiguration) GetAddr() string {
	return c.addr
}

func (c *redisConfiguration) GetDb() int {
	return c.db
}

func (c *redisConfiguration) GetHost() string {
	host := strings.Split(c.GetAddr(), ":")[0]
	if host == "" {
		return "localhost"
	}
	return host
}

func (c *redisConfiguration) GetPort() int {
	portString := strings.Split(c.GetAddr(), ":")[0]
	if portString == "" {
		return 6379
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		return 6379
	}
	return port
}
