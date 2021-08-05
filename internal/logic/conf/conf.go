package conf

import (
	"flag"
	"github.com/spf13/viper"
	"time"
)

var (
	confPath  string
	// Conf config
	Conf *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "./conf.yaml", "default config path")
	flag.Parse()
	viper.AddConfigPath(confPath)
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	_ = viper.Unmarshal(&Conf)
	return
}

// Default new a config with specified defualt value.
func Default() *Config {
	return &Config{
		Server: &Server{
			&Http{
					Addr: "127.0.0.1:9000",
					ReadTimeout : time.Second,
					WriteTimeout: time.Second,
					MaxHeaderBytes: 1 << 20,
				},
			&Grpc{Addr: ":9001"},
		},
		Data: &Data{
			Database:&Database{
				Dsn: "",
				Drive:"mysql",
			},
			Redis:&Redis{
				Addr:"127.0.0.1:6379",
				Password:"",
				Db:0,
				DialTimeout: time.Second,
				WriteTimeout:time.Second,
				ReadTimeout:time.Second,
			},
			Nsq:&Nsq{
				Addr:"127.0.0.1:4150",
			},
		},
	}
}

// Config config.
type Config struct {
	Server 	*Server
	Data	*Data
}

type Server struct {
	*Http
	*Grpc
}

type Data struct {
	*Database
	*Redis
	*Nsq
}

type Http struct {
	Addr		string
	ReadTimeout		time.Duration
	WriteTimeout	time.Duration
	MaxHeaderBytes	int
}

type Grpc struct {
	Addr		string
}

type Database struct {
	Dsn		string
	Drive	string
}

type Redis struct {
	Addr 			string
	Password 		string
	Db				int
	DialTimeout		time.Duration
	WriteTimeout	time.Duration
	ReadTimeout		time.Duration
}

type Nsq struct {
	Addr 			string
}



