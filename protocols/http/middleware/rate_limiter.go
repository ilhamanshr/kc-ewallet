package middleware

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"kc-ewallet/internals/errors"
	log_color "kc-ewallet/internals/helpers/color"
	redis_service "kc-ewallet/internals/helpers/redis/service"
	"kc-ewallet/protocols/http/response"

	"github.com/gin-gonic/gin"
)

const (
	keyPrefix          = "rate-limiter:%s:%s"
	violationKeyPrefix = "violation-checker:%s:%s"
)

type RateLimiterInterface interface {
	AllowRequest(handler, velocity string) bool
	CleanRateLimiter(handler, velocity string) bool
}

type RateLimiter struct {
	redis      redis_service.RedisServiceInterface
	rate       float64  // Rate at which the bucket leaks tokens (tokens per second)
	capacity   float64  // Maximum number of tokens the bucket can hold
	timeWindow float64  // since elapsed time is on seconds, enable option to set window time. Default is 1 as 1 second
	allowedIPs []string // whitelisted IPs will not go through IP rate limiter checkinng
}

// This module will return Rate Limiter with 1 request per second limit
func NewRateLimiter(redis redis_service.RedisServiceInterface, allowedIPs []string) *RateLimiter {
	return &RateLimiter{
		redis:      redis,
		rate:       1,
		capacity:   1,
		timeWindow: 1,
		allowedIPs: allowedIPs,
	}
}

func NewRateLimiterWithOption(redis redis_service.RedisServiceInterface, rate, capacity, timeWindow float64, allowedIPs []string) RateLimiterInterface {
	return &RateLimiter{
		redis:      redis,
		rate:       rate,
		capacity:   capacity,
		timeWindow: timeWindow,
		allowedIPs: allowedIPs,
	}
}

// AllowRequest will check based on handler name and velocity key
// velocity key can be used as IP, DeviceID, etc
func (r *RateLimiter) AllowRequest(handler, velocity string) bool {
	var (
		isAllowed = true // rate limiter error should not block incoming API requests
		err       error
	)

	defer func() {
		if err != nil {
			log_color.PrintRedf("rate limiter return error: %s", err)
		}
	}()

	// whitelisted IPs will not go through IP rate limiter checking
	if slices.Contains(r.allowedIPs, velocity) {
		return true
	}

	key := fmt.Sprintf(keyPrefix, handler, velocity)

	// Get Last Leak
	lastLeakString, err := r.redis.Hget(key, "lastLeak")
	if err != nil {
		return isAllowed
	}
	// Parse the string into a time.Time object
	var lastLeak time.Time
	if *lastLeakString == "" {
		lastLeak = time.Now().Add(-time.Hour)
	} else {
		lastLeak, err = time.Parse(time.RFC3339Nano, *lastLeakString)
		if err != nil {
			return isAllowed
		}
	}

	// Get Stored Tokensed
	tokensString, err := r.redis.Hget(key, "tokens")
	if err != nil {
		return isAllowed
	}
	var storedTokens float64
	if *tokensString != "" {
		storedTokens, err = strconv.ParseFloat(*tokensString, 64)
		if err != nil {
			return isAllowed
		}
	}

	currentTime := time.Now()
	elapsedTime := currentTime.Sub(lastLeak).Seconds() / r.timeWindow
	if r.timeWindow == 60 { // use proper minutes for 60 second time window
		elapsedTime = currentTime.Sub(lastLeak).Minutes()
	}
	leakedTokens := elapsedTime * r.rate

	storedTokens -= leakedTokens
	if storedTokens < 0 {
		storedTokens = 0
	}

	if storedTokens+1 <= r.capacity {
		storedTokens++
		err = r.redis.HsetWithExpiry(key, "lastLeak", currentTime.Format(time.RFC3339Nano), int(1*time.Hour))
		if err != nil {
			return isAllowed
		}
		err = r.redis.HsetWithExpiry(key, "tokens", strconv.FormatFloat(storedTokens, 'f', -1, 64), int(1*time.Hour))
		if err != nil {
			return isAllowed
		}

		return true
	}

	return false
}

// IncrViolationCount is a mini bot detection based on velocity frequency anomaly
// if violation have been done 10 times under one minute, mark that entity
func (r *RateLimiter) IncrViolationCount(handler, velocity string) {
	var (
		key        = fmt.Sprintf(violationKeyPrefix, handler, velocity)
		anomalyKey = fmt.Sprintf("%s:marked", key)
	)

	// Get Last Leak
	var violationLimitCount int
	r.redis.Get(key, &violationLimitCount)

	if violationLimitCount >= 10 {
		r.redis.SetWithExpiry(anomalyKey, true, int(24*time.Hour))
		return
	} else {
		r.redis.SetWithExpiry(key, violationLimitCount+1, int(1*time.Minute))
	}
}

// IsViolationMarked is a mini bot detection based on velocity frequency anomaly
func (r *RateLimiter) IsViolationMarked(handler, velocity string) bool {
	var (
		isAllowed  = false
		key        = fmt.Sprintf(violationKeyPrefix, handler, velocity)
		anomalyKey = fmt.Sprintf("%s:marked", key)
		err        error
	)

	var violationMarked bool
	err = r.redis.Get(anomalyKey, &violationMarked)
	if err != nil {
		return isAllowed
	}

	return violationMarked
}

func CheckRateLimit(limiter *RateLimiter, opts ...middlewareOptionFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt := defaultMiddlewareOption()
		for _, o := range opts {
			o(opt)
		}

		if strings.Contains(c.FullPath(), "private") {
			return
		}

		handlerName := getHandlerNameFromGinContext(c.HandlerNames())

		registered := true

		if opt.registerHandlers != nil {
			registered = opt.registerHandlers[handlerName]
		}

		if opt.excludedHandlers != nil {
			if opt.excludedHandlers[handlerName] {
				registered = false
			}
		}

		if !registered {
			return
		}

		if !limiter.AllowRequest(handlerName, c.ClientIP()) {
			limiter.IncrViolationCount(handlerName, c.ClientIP())
			response.RespondError(c, errors.TooManyRequests.New("Aktivitas Anda terdeteksi tidak wajar, hubungi tim Customer Service atau coba lagi nanti"))
			c.Abort()
		}

		if limiter.IsViolationMarked(handlerName, c.ClientIP()) {
			response.RespondError(c, errors.TooManyRequests.New("Aktivitas Anda terdeteksi tidak wajar, mohon coba lagi nanti"))
			c.Abort()
		}
	}
}

func (r *RateLimiter) CleanRateLimiter(handler, velocity string) bool {
	var (
		isAllowed = false
		key       = fmt.Sprintf(keyPrefix, handler, velocity)
		result    bool
		err       error
	)

	result, err = r.redis.Delete(key)
	if err != nil {
		return isAllowed
	}

	return result
}

func MockCheckRateLimit(limiter *RateLimiter, opts ...middlewareOptionFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		log_color.PrintGreen("rate limiter simulated as allowed")
	}
}
