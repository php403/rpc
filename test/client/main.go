package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	opHeartbeat      = uint32(2)
	opHeartbeatReply = uint32(3)
	opSendMsg 		= uint32(5)
	OpSendAppMsg 		= uint32(7)
	OpSendRoomMsg 		= uint32(9)
	opAuth           = uint32(3)
	opAuthReply      = uint32(8)
)

const (
	rawHeaderLen = uint16(16)
	heart        = 240 * time.Second
)

// Proto proto.
type Proto struct {
	PackLen   uint32  // package length
	HeaderLen uint16  // header length
	Ver       uint16  // protocol version
	Operation uint32  // operation for request
	Seq       uint32  // sequence number chosen by client
	Body      []byte // body
}

var (
	countDown  int64
	aliveCount int64
	sendCount int64
)


type UserToken struct {
	Uid 	int64			`json:"uid"`
	AppId 	int64			`json:"app_id"`
	AreaId 	int64			`json:"area_id"`
	TeamId 	int64			`json:"team_id"`
}

func main() {

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
	go result()

	ctx,cancel := context.WithCancel(context.Background())
	for i := begin; i < begin+num; i++ {
		go startClient(ctx,int64(begin))
		time.Sleep(time.Microsecond*5)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancel()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func startClient(ctx context.Context,uid int64) {
	atomic.AddInt64(&aliveCount, 1)
	defer func() {
		glog.Infof("error exit")
		atomic.AddInt64(&aliveCount, -1)
	}()
	// connnect to server
	conn, err := net.Dial("tcp", os.Args[3])
	if err != nil {
		glog.Errorf("net.Dial(%s) error(%v)", os.Args[3], err)
		return
	}

	seq := uint32(3)
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)
	proto := new(Proto)
	//暂定10个区
	//1个app
	userToken := &UserToken{
		Uid: uid,
		AppId: 1,
		AreaId: uid/2,
		TeamId:uid/10,
	}
	proto.Ver = 1
	proto.Operation = opAuth
	proto.Seq = seq
	proto.Body,_ = json.Marshal(userToken)
	if err = tcpWriteProto(wr, proto); err != nil {
		_ = fmt.Errorf("tcpWriteProto() error(%v)", err)
	}
	for err = tcpReadProto(rd, proto); err != nil;{
		if proto.Operation == opAuthReply {
			fmt.Println(string(proto.Body))
			break
		}
	}

	go func() {
		for {
			select {
			case <- ctx.Done():
				break
			default:
				// heartbeat
				proto := new(Proto)
				proto.Ver = 1
				proto.Operation = OpSendAppMsg
				proto.Seq = seq
				proto.Body = []byte("hello world")
				if err = tcpWriteProto(wr, proto); err != nil {
					glog.Errorf("key:%d tcpWriteProto() error(%v)", err.Error())
					return
				}
				atomic.AddInt64(&sendCount,1)
				time.Sleep(time.Millisecond*200)
			}
		}
	}()
	time.Sleep(time.Second*1)
	go func() {
		// heartbeat
		for {
			select {
			case <- ctx.Done():
				break
			default:
				proto := new(Proto)
				if err = tcpReadProto(rd, proto); err != nil {
					glog.Errorf("key:%d tcpReadProto() error(%v)", uid, err)
					return
				}
				atomic.AddInt64(&countDown, 1)
			}
		}
	}()
}

func result() {
	var (
		lastTimes int64
		interval  = int64(5)
	)
	for {
		nowCount := atomic.LoadInt64(&countDown)
		nowAlive := atomic.LoadInt64(&aliveCount)
		newSend := atomic.LoadInt64(&sendCount)
		diff := nowCount - lastTimes
		lastTimes = nowCount
		fmt.Println(fmt.Sprintf("%s alive:%d down:%d send:%d down/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowCount,newSend, diff/interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func tcpWriteProto(wr *bufio.Writer, proto *Proto) (err error) {
	// write
	if err = binary.Write(wr, binary.BigEndian, uint32(rawHeaderLen)+uint32(len(proto.Body))); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, rawHeaderLen); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Ver); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Operation); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Seq); err != nil {
		return
	}
	if proto.Body != nil {
		if err = binary.Write(wr, binary.BigEndian, proto.Body); err != nil {
			return
		}
	}
	err = wr.Flush()
	return
}

func tcpReadProto(rd *bufio.Reader, proto *Proto) (err error) {
	var (
		packLen   uint32
		headerLen uint16
	)
	// read
	if err = binary.Read(rd, binary.BigEndian, &packLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &headerLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Ver); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Operation); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Seq); err != nil {
		return
	}
	var (
		n, t    int
		bodyLen = int(packLen - uint32(headerLen))
	)
	if bodyLen > 0 {
		proto.Body = make([]byte, bodyLen)
		for {
			if t, err = rd.Read(proto.Body[n:]); err != nil {
				return
			}
			if n += t; n == bodyLen {
				break
			}
		}
	} else {
		proto.Body = nil
	}
	return
}

