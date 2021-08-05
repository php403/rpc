package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/php403/im/internal/comet"
	"github.com/php403/im/internal/comet/conf"
	"github.com/php403/im/pkg/log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

func main()  {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	encoderConfig,atomicLevel := log.InitConfig()
	logger := log.NewZapLogger(encoderConfig,atomicLevel)

	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	srv := comet.NewServer(conf.Conf,logger)
	if err := comet.InitTCP(srv, conf.Conf.TCP.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	ctx,cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
				err := srv.OperateQueueMsg()
				if err != nil {
					logger.Warnf("queue msg error")
				}
			}
		}
	}()

	go func() {
		for  {
			fmt.Println(fmt.Sprintf("allUser %d",atomic.LoadInt64(&srv.AllUser)))
			time.Sleep(time.Second * time.Duration(10))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		logger.Info("comet","comet get a signal"+ s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}

			/*rpcSrv.GracefulStop()
			srv.Close()*/
			logger.Info("comet exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}













