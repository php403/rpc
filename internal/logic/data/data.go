package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/nsqio/go-nsq"
	"github.com/php403/im/internal/logic/conf"
	"github.com/php403/im/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	//"github.com/Shopify/sarama"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	db   *gorm.DB
	rdb *redis.Client
	nsq *nsq.Producer
	log  log.Logger
}

func NewData(conf *conf.Data,logger log.Logger) (*Data, func(), error) {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(conf.Database.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "im_",
			SingularTable: true,
	},})
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           conf.Redis.Db,
		DialTimeout:  conf.Redis.DialTimeout,
		WriteTimeout: conf.Redis.WriteTimeout,
		ReadTimeout:  conf.Redis.ReadTimeout,
	})

	// Instantiate a producer.
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(conf.Nsq.Addr, config)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	return &Data{db, rdb,producer,logger}, func() {
		sqlDb,err := db.DB()
		if err != nil {
			logger.Log(log.LevelError,"logic db err!",err)
		}
		if err := sqlDb.Close(); err != nil {
			logger.Log(log.LevelError,"logic sqldb err!",err)
		}
		if err := rdb.Close(); err != nil {
			logger.Log(log.LevelError,"logic redis client err!",err)
		}
		producer.Stop()
		}, nil
}
