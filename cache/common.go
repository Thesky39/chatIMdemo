package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
)

// go get github.com/go-redis/redis
var (
	RedisClient   *redis.Client
	RedisDb       string
	RedisAddr     string
	RedisPassword string
	RedisDbName   string
)

func Init() {
	file, err := ini.Load("./conf/conf.ini") //加载配置信息
	if err != nil {
		fmt.Println("redis ini load failed", err)
	}
	LoadRedis(file) //读取配置信息
	Redis()         //redis连接
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPassword = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()

}
func Redis() {
	db, _ := strconv.ParseInt(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})
	_, err := client.Ping().Result()
	//  go get github.com/sirupsen/logrus  日志包
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}
