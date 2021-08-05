package comet

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/php403/im/api/logic"
	"github.com/php403/im/api/protocol"
	"github.com/php403/im/internal/comet/conf"
	"github.com/php403/im/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"math/rand"
	"time"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
	grpcBackoffMaxDelay       = time.Second * 3
)


type Server struct {
	c         	*conf.Config
	serverName 	string
	buckets   	[]*Bucket // subkey bucket
	bucketIdx 	uint32
	serverID  	string
	log			log.Logger
	rpcClient 	logic.LogicClient
	NsqConsumer *Nsq

	AllRcvNumsMsg 	int64
	AllSendNumsMsg 	int64
	AllUser			int64
}

func NewServer(c *conf.Config,logger log.Logger) *Server {
	s := &Server{
		c:         c,
		log: logger,
		rpcClient: newLogicClient(c.RPCClient),
	}
	// init bucket
	//todo 修改为读取c.conf.env配置
	s.serverName = "comet1"
	s.buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(c.Bucket.Size)
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(c.Bucket)
	}
	s.NsqConsumer = NewNsq(c.Nsq)
	s.serverID = c.Env.Host
	return s
}

func newLogicClient(c *conf.RPCClient) logic.LogicClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	config := clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:10*time.Second,
	}
	cli,err := clientv3.New(config)
	if err!= nil {
		panic(err)
	}
	_, err = cli.Status(ctx, config.Endpoints[0])
	if err!= nil {
		panic(err)
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	if err!= nil {
		panic(err)
	}
	conn, err := grpc.DialContext(ctx, "127.0.0.1:9001",
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithResolvers(etcdResolver),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		}...)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}

// Buckets return all buckets.
func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

// RandServerHearbeat rand server heartbeat.
func (s *Server) RandServerHearbeat() time.Duration {
	return (minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat))))
}

// Close close the server.
func (s *Server) Close() (err error) {
	return
}

func (s *Server) Bucket(uid int64) *Bucket {
	idx := uid % int64(s.bucketIdx)
	return s.buckets[idx]
}

func(s *Server) OperateQueueMsg() error {
	msg := <- s.NsqConsumer.msg
	if len(msg.Body) == 0 {
		return nil
	}
	pushMsg := new(logic.PushMsg)
	err := proto.Unmarshal(msg.Body,pushMsg)
	if err != nil {
		panic(err)
	}
	uid := pushMsg.UserId
	bucket := s.Bucket(uid)
	pb := &protocol.Proto{}
	pb.Ver = uint32(1)
	pb.Body = pushMsg.Msg
	switch pushMsg.Type {
	case logic.PushMsg_BROADCAST:
		pb.Op = protocol.OpSendAppMsgReply
		bucket.apps[pushMsg.AppId].Push(pb)
	case logic.PushMsg_ROOM:
		pb.Op = protocol.OpSendRoomMsgReply
		bucket.rooms[pushMsg.RoomId].Push(pb)
	}
	return nil
}







