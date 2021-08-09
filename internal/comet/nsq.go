package comet

import (
	"github.com/nsqio/go-nsq"
	"github.com/php403/im/internal/comet/conf"
)



type Nsq struct {
	topic 		string
	consumer 	*nsq.Consumer
	msg 		chan *nsq.Message
}


func (Nsq *Nsq) HandleMessage(msg *nsq.Message) error {
	 Nsq.msg <- msg
	 return nil
}


func NewNsq(c *conf.Nsq) *Nsq {
	config:=nsq.NewConfig()
	//todo 改成server name
	consumer, err := nsq.NewConsumer(c.Topic, "comet1", config)  //topic， channel， config
	if nil != err {
		panic(err)
	}
	msgChan := make(chan *nsq.Message,1024)
	nsq1 := &Nsq{c.Topic,consumer,msgChan}
	consumer.AddHandler(nsq1)
	err = consumer.ConnectToNSQD(c.Addr)
	if err != nil {
		panic(err)
	}
	return nsq1
}






