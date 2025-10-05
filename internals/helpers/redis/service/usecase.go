package service

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

func (r RedisService) SetWithLock(key string, data interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	// Try to aquire lock
	lockKey := fmt.Sprintf("%s_lock", key)
	reply, err := conn.Do("SETNX", lockKey, uuid.New().String())
	if err != nil {
		return err
	}
	if _, err = conn.Do("EXPIRE", lockKey, DefaultLockTimeout); err != nil {
		return err
	}

	// Return if failed to get lock
	replyValue, ok := reply.(int64)
	if !ok {
		return ErrInvalidReply
	}
	if replyValue == 0 {
		return ErrFailedToGetLock
	}

	// Set value
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := conn.Do("SET", key, value); err != nil {
		return err
	}

	// Release the lock
	if _, err := conn.Do("DEL", lockKey); err != nil {
		return nil
	}

	return nil
}

// Set a key/value
func (r RedisService) Set(key string, data interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

// SetWithExpiry a key/value
func (r RedisService) SetWithExpiry(key string, data interface{}, time int) error {
	conn := r.Pool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// SetnxWithExpiry a key/value if not exists
func (r RedisService) SetnxWithExpiry(key string, data interface{}, time int) error {
	conn := r.Pool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETNX", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func (r RedisService) Hset(key, field, value string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return err
	}

	return nil
}

func (r RedisService) HsetWithExpiry(key, field, value string, time int) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func (r RedisService) Exists(key string) bool {
	conn := r.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func (r RedisService) Get(key string, data interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	err = json.Unmarshal(reply, &data)
	if err != nil {
		return err
	}

	return nil
}

func (r RedisService) Hget(key, field string) (*string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("HGET", key, field))
	if err != nil && err.Error() != "redigo: nil returned" {
		return nil, err
	}
	return &value, nil
}

// Delete delete a key
func (r RedisService) Delete(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// Incr increment a key value
func (r RedisService) Incr(key string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return reply, nil

}
