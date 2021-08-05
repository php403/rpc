package main

import (
	"context"
	"github.com/golang/glog"
	"github.com/php403/im/api/logic"
	"github.com/php403/im/internal/logic/conf"
	"github.com/php403/im/internal/logic/server"
	"github.com/php403/im/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"

	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)


func main()  {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	encoderConfig,atomicLevel := log.InitConfig()
	logger := log.NewZapLogger(encoderConfig,atomicLevel)
	app, cleanup, err := initApp(conf.Conf.Server, conf.Conf.Data, logger)
	if err != nil {
		panic(err)
	}
	app.Start()


	defer cleanup()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		glog.Infof("goim-logic get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Infof("goim-logic exit")
			glog.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

type App struct {
	grpc 	*server.GrpcServer
	http	*http.Server
	log		log.Logger
	ctx 	context.Context
	cli		*clientv3.Client
}

func newApp(logger log.Logger, hs *http.Server,grpc *server.GrpcServer) *App {
	return &App{
		http: hs,
		log: logger,
		grpc: grpc,
	}
}

func(app *App) Start()  {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := app.http.ListenAndServe()
		wg.Done()
		panic(err)
	}()
	wg.Add(1)
	go func() {
		listen,err := net.Listen("tcp",app.grpc.C.Grpc.Addr)
		if err != nil {
			wg.Done()
			panic(err)
		}
		err = app.grpc.Grpc.Serve(listen)
		if err != nil {
			wg.Done()
			panic(err)
		}
		logic.RegisterLogicServer(app.grpc.Grpc,app.grpc.User)
		wg.Done()
	}()
	wg.Wait()

}

func(app *App) Stop()  {
	if err := app.http.Shutdown(app.ctx); err != nil {
		panic(err)
	}
	app.grpc.Grpc.Stop()
	if err := app.cli.Close();err != nil{
		panic(err)
	}
}

func RegDiscovery() *clientv3.Client {
	cli, err := clientv3.NewFromURL("http://127.0.0.1:2379")
	if err != nil {
		panic(err)
	}
	em,err := endpoints.NewManager(cli, "etcd:///game/im/logic")
	if err != nil {
		panic(err)
	}
	err = em.AddEndpoint(context.TODO(), "etcd:///game/im/logic/"+"1",endpoints.Endpoint{Addr:"127.0.0.1:9001"})
	if err != nil {
		panic(err)
	}

	return cli
}

