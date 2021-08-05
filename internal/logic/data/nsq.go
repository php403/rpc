package data

func(ur * UserRepo) SendMsg(topic string,body []byte) (err error) {
	err = ur.data.nsq.Publish(topic,body)
	return
}


