package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main()  {
	ctx := context.Background()
	config := clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:10*time.Second,
	}
	cli,err := clientv3.New(config)
	if err!= nil {
		panic(err)
	}
	_, err = cli.Status(ctx, config.Endpoints[0])
	if err != nil {
		panic(err)
	}
	cli.Put(ctx,"1111","2222")
}