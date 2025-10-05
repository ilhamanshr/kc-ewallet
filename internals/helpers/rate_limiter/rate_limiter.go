package rate_limit

import (
	"kc-ewallet/configurations"
	"kc-ewallet/internals/helpers/redis"
	"kc-ewallet/internals/helpers/redis/service"
)

func InitCoreRedis(redisConfiguration configurations.IRedisConfiguration) {
	redis.InitRedisPool(&redis.Configuration{
		Host: redisConfiguration.GetHost(),
		Port: redisConfiguration.GetPort(),
		Db:   redisConfiguration.GetDb(),
	})
}

func NewCacheService() service.RedisServiceInterface {
	return service.NewRedisService(redis.RedisPool)
}
