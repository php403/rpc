package comet

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/php403/im/api/protocol"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	maxInt = 1<<31 - 1
)

func InitTCP(server *Server, bind string, accept int) (err error) {
	var (
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {

		server.log.Error("net.ResolveTCPAddr(tcp, %s) error(%v)", bind, err)
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		server.log.Error("net.ListenTCP(tcp, %s) error(%v)", bind, err)
		return
	}
	server.log.Infof("start tcp listen: %s", bind)
	// split N core accept
	for i := 0; i < accept; i++ {
		go acceptTCP(server, listener)
	}
	return
}

func acceptTCP(server *Server, listener *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			server.log.Errorf("listener.Accept(\"%s\") error(%v)", listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(server.c.TCP.KeepAlive); err != nil {
			server.log.Errorf("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(server.c.TCP.Rcvbuf); err != nil {
			server.log.Errorf("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(server.c.TCP.Sndbuf); err != nil {
			server.log.Errorf("conn.SetWriteBuffer() error(%v)", err)
			return
		}
		go server.ServeTCP(conn, r)
		if r++; r == maxInt {
			r = 0
		}
	}
}

func (s *Server) ServeTCP(conn *net.TCPConn, r int) {
	var(
		user = NewUser(s.c.Protocol.CliProto, s.c.Protocol.SvrProto)
		err error
		b       *Bucket
		lastHb  = time.Now()
	)
	user.Reader = bufio.NewReaderSize(conn,s.c.TCP.Reader)
	user.Writer = bufio.NewWriterSize(conn,s.c.TCP.Writer)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//todo 改环形buf存储 用完回收 或者sync.pool实现 减少gc
	p := &protocol.Proto{}
	p.Ver = 1
	p.Op = protocol.OpClosetReply
	userInfo,err := s.authTCP(ctx, user.Reader, user.Writer, p)
	if err != nil {
		//todo 添加im全局错误码
		p.Body = []byte("auth error!")
		s.log.Errorf("Auth err! %V",err.Error())
		_ = p.WriteTCP(user.Writer)
		_ = conn.Close()
		return
	}
	atomic.AddInt64(&s.AllUser,1)
	p.Body = []byte("success")
	err = p.WriteTCP(user.Writer)
	if err != nil {
		s.log.Error(err.Error())
	}
	user.Uid = userInfo.Uid
	user.AppId = strconv.FormatInt(userInfo.AppId,10)
	user.RoomId = strconv.FormatInt(userInfo.TeamId,10)
	token,err := json.Marshal(userInfo)
	if err != nil {
		//todo 解析协议失败error
		s.log.Errorf("token marshal err uid:%v!",user.Uid)
	}
	user.token = string(token)
	b = s.Bucket(user.Uid)
	err = b.Put(user.AppId,user.RoomId,user)
	go s.dispatchTCP(ctx,conn,user)

	serverHeartbeat := s.RandServerHearbeat()
	for {
		p = &protocol.Proto{}
		if err = p.ReadTCP(user.Reader); err != nil {
			break
		}
		if p.Op == protocol.OpHeartbeat {
			//todo 心跳需要完善超时处理
			if now := time.Now(); now.Sub(lastHb) > serverHeartbeat {
				lastHb = now
			}
		} else {
			if err = s.Operate(ctx, p, user, b); err != nil {
				break
			}
		}

	}

}

type UserToken struct {
	Uid 	int64			`json:"uid"`
	AppId 	int64			`json:"app_id"`
	AreaId 	int64			`json:"area_id"`
	TeamId 	int64			`json:"team_id"`
}


func (s *Server) authTCP(ctx context.Context, rr *bufio.Reader, wr *bufio.Writer, p *protocol.Proto) (token *UserToken,err error) {
	for {
		if err = p.ReadTCP(rr); err != nil {
			return
		}
		if p.Op == protocol.OpAuth {
			break
		} else {
			s.log.Errorf("tcp request operation(%d) not auth", p.Op)
		}
	}
	err = json.Unmarshal(p.Body,&token)
	if err != nil {
		s.log.Error(fmt.Sprintf("authTCP.Connect(token:%v).err(%v)",string(p.Body), err))
		return
	}
	//todo 连接验证
	/*if server, appId, areaId, roomId, uid, err = s.Connect(ctx, string(p.Body)); err != nil {
		s.log.Error(fmt.Sprintf("authTCP.Connect(key:%v).err(%v)", uid, err))
		return
	}
	*/
	p.Op = protocol.OpAuthReply
	p.Body = nil
	if err = p.WriteTCP(wr); err != nil {
		s.log.Errorf("authTCP.WriteTCP(key:%v).err(%v)", token.Uid, err)
		return
	}
	err = wr.Flush()
	return
}

func (s *Server) dispatchTCP(ctx context.Context,conn *net.TCPConn, user *User) {
	var err error
	for {
		select {
		case <- ctx.Done():
			return
		case msg := <-user.GetMsgs():
			if err = msg.WriteTCP(user.Writer); err != nil {
				s.log.Errorf("key: %s dispatch tcp error(%v)", user.Uid, err)
				_ = conn.Close()
			}
		}
	}
}

func (s *Server) Operate(ctx context.Context, p *protocol.Proto, user *User, b *Bucket) error {
	if err := s.Receive(ctx, user.token, p); err != nil {
		s.log.Errorf("s.Report %v",  err.Error())
	}
	p.Body = nil

	//todo need add b.changeRoom
	return nil
}




