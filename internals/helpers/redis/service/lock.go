package service

import "github.com/gomodule/redigo/redis"

func (r RedisService) Acquire(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("SETNX", key, "lock"))
}

func (r RedisService) Release(key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}
