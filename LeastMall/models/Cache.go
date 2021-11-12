package models

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"time"

	"github.com/beego/beego/v2/adapter/cache"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/redis"
)

var redisClient cache.Cache
var enableRedis, _ = beego.AppConfig.Bool("enableRedis")
var redisTime, _ = beego.AppConfig.Int("redisTime")

func init() {
	if enableRedis {
		redisKey, _ := beego.AppConfig.String("redisKey")
		redisConn, _ := beego.AppConfig.String("redisConn")
		redisDBNum, _ := beego.AppConfig.String("redisDbNum")
		redisPwd, _ := beego.AppConfig.String("redisPwd")
		config := map[string]string{
			"key":      redisKey,
			"conn":     redisConn,
			"dbNum":    redisDBNum,
			"password": redisPwd,
		}
		bytes, _ := json.Marshal(config)

		redisClient, err = cache.NewCache("redis", string(bytes))
		_, _ = cache.NewCache("redis", string(bytes))
		if err != nil {
			log.Error("连接redis数据库失败")
		} else {
			log.Info("连接redis数据库成功")
		}

	}
}

type cacheDb struct{}

// Set 写入数据的方法
func (c cacheDb) Set(key string, value interface{}) {
	if enableRedis {
		bytes, _ := json.Marshal(value)
		_ = redisClient.Put(key, string(bytes), time.Second*time.Duration(redisTime))
	}
}

// Get 获取数据的方法
func (c cacheDb) Get(key string, obj interface{}) bool {
	if enableRedis {
		if redisStr := redisClient.Get(key); redisStr != nil {
			fmt.Println("在redis里面读取数据...")
			redisValue, ok := redisStr.([]uint8)
			if !ok {
				fmt.Println("获取redis数据失败")
				return false
			}
			_ = json.Unmarshal([]byte(redisValue), obj)
			return true
		}
		return false
	}
	return false
}

var CacheDb = &cacheDb{}
