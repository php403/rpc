package main

import (
	"flag"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var sendCount int64
var recvCount int64

type myMessageHandler struct {}

func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	return nil
}
func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		panic(err)
	}


	consumer.AddHandler(&myMessageHandler{})

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQD("172.17.0.16:4150")
	if err != nil {
		panic(err)
	}
	config = nsq.NewConfig()
	producer, err := nsq.NewProducer("172.17.0.16:4150", config)
	if err != nil {
		panic(err)
	}
	messageBody := []byte("hello")
	topicName := "topic"

	for {
		for i := begin; i < begin+num; i++ {
			go func() {
				err = producer.Publish(topicName, messageBody)
				if err != nil {
					panic(err)
				}
			}()
		}
		time.Sleep(time.Second*1)
	}
	time.Sleep(time.Second*1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

