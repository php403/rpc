package conf

import (
	"flag"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
	xtime "time"
)

var (
	confPath  string
	deployEnv string
	host      string
	addrs     string
	weight    int64
	offline   bool
	debug     bool

	// Conf config
	Conf *Config
)

func init() {
	var (
		defHost, _    = os.Hostname()
		defAddrs      = os.Getenv("ADDRS")
		defWeight, _  = strconv.ParseInt(os.Getenv("WEIGHT"), 10, 32)
		defOffline, _ = strconv.ParseBool(os.Getenv("OFFLINE"))
		defDebug, _   = strconv.ParseBool(os.Getenv("DEBUG"))
	)
	flag.StringVar(&confPath, "conf", "comet-example.toml", "default config path.")
	flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	flag.StringVar(&host, "host", defHost, "machine hostname. or use default machine hostname.")
	flag.StringVar(&addrs, "addrs", defAddrs, "server public ip addrs. or use ADDRS env variable, value: 127.0.0.1 etc.")
	flag.Int64Var(&weight, "weight", defWeight, "load balancing weight, or use WEIGHT env variable, value: 10 etc.")
	flag.BoolVar(&offline, "offline", defOffline, "server offline. or use OFFLINE env variable, value: true/false etc.")
	flag.BoolVar(&debug, "debug", defDebug, "server debug. or use DEBUG env variable, value: true/false etc.")

	viper.AddConfigPath(confPath)
	viper.SetConfigType("yaml")


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
		Env:       &Env{DeployEnv: deployEnv, Host: host, Weight: weight, Addrs: strings.Split(addrs, ","), Offline: offline},
		Debug:     debug,
		RPCClient: &RPCClient{
			Dial:    xtime.Duration(time.Second),
			Timeout: xtime.Duration(time.Second),
		},
		RPCServer: &RPCServer{
			Network:           "tcp",
			Addr:              ":3109",
			Timeout:           xtime.Duration(time.Second),
			IdleTimeout:       xtime.Duration(time.Second * 60),
			MaxLifeTime:       xtime.Duration(time.Hour * 2),
			ForceCloseWait:    xtime.Duration(time.Second * 20),
			KeepAliveInterval: xtime.Duration(time.Second * 60),
			KeepAliveTimeout:  xtime.Duration(time.Second * 20),
		},
		TCP: &TCP{
			Bind:         ":3101",
			Sndbuf:       4096,
			Rcvbuf:       4096,
			KeepAlive:    false,
			Reader:       32,
			ReadBuf:      1024,
			ReadBufSize:  8192,
			Writer:       32,
			WriteBuf:     1024,
			WriteBufSize: 8192,
		},
		Bucket: &Bucket{
			Size:          32,
			Channel:       1024,
			App:           8,
			Room:          4096,
			RoutineAmount: 32,
			RoutineSize:   1024,
		},
		Protocol: &Protocol{
			Timer:            32,
			TimerSize:        2048,
			CliProto:         5,
			SvrProto:         10,
			HandshakeTimeout: xtime.Duration(time.Second * 5),
		},
		Nsq: &Nsq{
			Topic: "gameImcomet1",
			Addr:"172.18.180.81:4150",
		},
	}
}

type Nsq struct {
	Topic 	string
	Addr 	string
}

// Config is comet config.
type Config struct {
	Debug     bool
	TCP       *TCP
	Env       *Env
	Protocol  *Protocol
	Bucket    *Bucket
	RPCClient *RPCClient
	RPCServer *RPCServer
	LogConfig *Config
	Nsq		  *Nsq
}

type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
	Weight    int64
	Offline   bool
	Addrs     []string
}



// RPCClient is RPC client config.
type RPCClient struct {
	Dial    xtime.Duration
	Timeout xtime.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           xtime.Duration
	IdleTimeout       xtime.Duration
	MaxLifeTime       xtime.Duration
	ForceCloseWait    xtime.Duration
	KeepAliveInterval xtime.Duration
	KeepAliveTimeout  xtime.Duration
}

// TCP is tcp config.
type TCP struct {
	Bind         string
	Sndbuf       int
	Rcvbuf       int
	KeepAlive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Websocket is websocket config.
type Websocket struct {
	Bind        []string
	TLSOpen     bool
	TLSBind     []string
	CertFile    string
	PrivateFile string
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int
	TimerSize        int
	SvrProto         int
	CliProto         int
	HandshakeTimeout xtime.Duration
}

// Bucket is bucket config.
type Bucket struct {
	Size          int
	Channel       int
	App           int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}


