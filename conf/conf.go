package conf

import (
	"chat/model"
	"context"
	"fmt"
	"strings"

	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
)

var (
	// go get go.mongodb.org/mongo-driver/mongo
	// go get go.mongodb.org/mongo-driver/mongo/options
	MongoDBClient *mongo.Client
	AppMode       string
	HttpPort      string
	DB            string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	//RedisDb     string
	//RedisAddr   string
	//RedisPw     string
	//RedisDbName string

	MongoDBName string
	MongoDBAddr string
	MongoDBPwd  string
	MongoDBPort string
)

func Init() {
	//本地读取环境变量   go get gopkg.in/ini.v1
	file, err := ini.Load("./conf/conf.ini")

	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)

	LoadMongoDB(file)
	fmt.Println(MongoDBAddr)
	MongoDB() // MongoDB连接
	path := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True"}, "")
	model.Database(path) // mysql连接
}
func MongoDB() {
	uri := fmt.Sprintf("mongodb://%s:%s", MongoDBAddr, MongoDBPort)
	fmt.Println("MongoDB URI:", uri) // 直接打印出拼接的 URI 字符串

	// 创建客户端连接选项
	clientOptions := options.Client().ApplyURI(uri)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info(err)

	}
	err = MongoDBClient.Ping(context.TODO(), nil)
	if err != nil {
		logging.Info(err)
	}
	logging.Info("MongoDB connect success")
}
func LoadServer(file *ini.File) {
	//读取server下的AppMode转换为字符串
	AppMode = file.Section("server").Key("AppMode").String()
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8080")

}
func LoadMysql(file *ini.File) {
	DB = file.Section("mysql").Key("DB").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").MustString("3306")
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()

}

//	func LoadRedis(file *ini.File) {
//		RedisDb = file.Section("redis").Key("RedisDb").String()
//		RedisAddr = file.Section("redis").Key("RedisAddr").String()
//		RedisPw = file.Section("redis").Key("RedisPw").String()
//		RedisDbName = file.Section("redis").Key("RedisDbName").String()
//
// }
func LoadMongoDB(file *ini.File) {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()

}
