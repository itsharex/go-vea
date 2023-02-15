package configs

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/spf13/viper"
	"go-vea/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"
)

var (
	AppConfig         Config
	globalDB          *gorm.DB
	globalRedisClient *redis.Client
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	Database DBConfig     `yaml:"database"`
	Logger   LoggerConfig `yaml:"logger"`
	JWT      JWTConfig    `yaml:"jwt"`
	Redis    RedisConfig  `yaml:"redis"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DBConfig struct {
	Platform string
	Host     string
	Port     int16
	Dbname   string
	Username string
	Password string
	Arg      string
}

type LoggerConfig struct {
	Path   string
	Level  uint32
	Stdout bool
}

type JWTConfig struct {
	Header     string
	Secret     string
	ExpireTime string
}

type RedisConfig struct {
	Addr     string
	Password string
	RDB      int
}

func InitConfig() {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		//panic(err)
		global.Logger.Error("获取项目路径失败", err)
		return
	}
	cfg := viper.New()
	// 设置读取的文件路径
	cfg.AddConfigPath(path + "/configs")
	// 设置读取的文件名
	cfg.SetConfigName("config")
	// 设置文件的类型
	cfg.SetConfigType("yaml")
	// 尝试进行配置读取
	if err := cfg.ReadInConfig(); err != nil {
		//panic(err)
		global.Logger.Error("配置读取失败", err)
		return
	}

	// 把配置文件读取到结构体上
	//var config Config
	// 将配置文件绑定到 config 上
	err = cfg.Unmarshal(&AppConfig)

	//fmt.Println(AppConfig.Server.Mode)
	//fmt.Println(cfg.GetString("database.host"))

	initDb()
	InitRedis()
	InitIp2region()
}

func initDb() {
	sqlDbCfg := AppConfig.Database
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		sqlDbCfg.Username,
		sqlDbCfg.Password,
		sqlDbCfg.Host,
		sqlDbCfg.Port,
		sqlDbCfg.Dbname,
		sqlDbCfg.Arg,
	)
	config := &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			// 使用单数表名，结构体 User 对应的表名为 user
			SingularTable: true,
		},
	}
	db, err := gorm.Open(mysql.Open(dns), config)

	sqlDB, _ := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(20)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	if err != nil {
		//panic("连接数据库失败 " + err.Error())
		global.Logger.Error("数据库连接失败", err)
		return
	}

	// 赋值给私有全局变量
	globalDB = db
}

func GetDB(ctx context.Context) *gorm.DB {
	// db.Session(&Session{Context: ctx}) 每次创建新Session 各db操作互不影响
	return globalDB.WithContext(ctx)
}

func InitRedis() {
	// https://juejin.cn/post/7027347979065360392
	// https://juejin.cn/post/7034322568014364680
	rdbCfg := AppConfig.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     rdbCfg.Addr,
		Password: rdbCfg.Password,
		DB:       rdbCfg.RDB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		//panic(err)
		global.Logger.Error("redis初始化失败", err)
		return
	}
	globalRedisClient = client
	global.Redis = client
}

func GetRedisClient(ctx context.Context) *redis.Client {
	return globalRedisClient.WithContext(ctx)
}

func InitIp2region() {
	path, err := os.Getwd()
	dbPath := path + "/configs/ip2region.xdb"
	if err != nil {
		global.Logger.Error("获取项目路径失败", err)
		return
	}
	// 1、从 dbPath 加载 VectorIndex 缓存，把下述 vIndex 变量全局到内存里面。
	vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
	if err != nil {
		fmt.Printf("failed to load vector index from `%s`: %s\n", dbPath, err)
		return
	}

	// 2、用全局的 vIndex 创建带 VectorIndex 缓存的查询对象。
	searcher, err := xdb.NewWithVectorIndex(dbPath, vIndex)
	if err != nil {
		fmt.Printf("failed to create searcher with vector index: %s\n", err)
		return
	}
	global.IpSearcher = searcher
	// 备注：并发使用，全部 goroutine 共享全局的只读 vIndex 缓存，每个 goroutine 创建一个独立的 searcher 对象
}
