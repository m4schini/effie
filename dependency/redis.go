package dependency

import (
	"github.com/go-redis/redis/v8"
	"github.com/m4schini/exstate"
	er "github.com/m4schini/exstate/redis"
	"os"
	"sync"
)

// ====================
// REDIS CONNECTION
// ====================

var redisClient *redis.Client
var eCache exstate.Cache
var eSource exstate.Source
var redisAddr string
var redisPass string
var redisDb int
var redisClientInitiliased = false
var redisInitLock sync.Mutex

func readyRedis() {
	redisInitLock.Lock()
	defer redisInitLock.Unlock()
	if !redisClientInitiliased {
		initRedisClient()
	}
}

func initRedisClient() {
	redisDb = 0

	redisAddr = os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Warn("REDIS_ADDR is missing")
	}
	redisPass = os.Getenv("REDIS_PASS")
	if redisPass == "" {
		log.Warn("REDIS_PASS is missing")
	}

	r, err := er.New(redisAddr, redisPass, redisDb)
	if err != nil {
		log.Fatal("redis connection failed")
	}
	log.Infow("initialised new redis connection")

	eCache = r
	eSource = r

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDb,
	})

	redisClientInitiliased = true
}

func RedisClient() *redis.Client {
	readyRedis()

	return redisClient
}

func ECache() exstate.Cache {
	readyRedis()

	return eCache
}

func ESource() exstate.Source {
	readyRedis()

	return eSource
}
