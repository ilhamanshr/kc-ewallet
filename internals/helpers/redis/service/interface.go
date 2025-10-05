package service

//go:generate mockgen -destination=mocks/mock.go -source=interface.go RedisServiceInterface

type RedisServiceInterface interface {
	SetWithLock(key string, data interface{}) error
	Set(key string, data interface{}) error
	SetWithExpiry(key string, data interface{}, time int) error
	SetnxWithExpiry(key string, data interface{}, time int) error
	Hset(key, field, value string) error
	HsetWithExpiry(key, field, value string, time int) error
	Exists(key string) bool
	Get(key string, data interface{}) error
	Hget(key, field string) (*string, error)
	Delete(key string) (bool, error)
	Acquire(key string) (bool, error)
	Release(key string) error
}
