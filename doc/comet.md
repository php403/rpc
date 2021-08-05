# 1协议设计
| 协议字段 | 长度 | 协议解析 |
|:------ | :------ | :------ |
|package len | 4 byte | 包长度 |
|header len | 2 byte | 头长度 |
|protocol version | 2 byte | 协议版本 1开始(不改动协议不修改此字段)|
|operation | 4 byte | 业务区分 1 心跳检测 2游戏聊天 3客服消息|
|sequence id | 4 byte | 序列号 |
|body | package len -header len | 消息长度 |

#2 结构设计
1 bucket->app->room->conn bucket 减少锁粒度 

    //bucket
    struct{
        c     *conf.Bucket
        cLock sync.RWMutex        // protect the channels for chs
        chs   map[string]*Channel // map sub key to a channel
        // room
        rooms       map[string]*Room // bucket room channels
        routines    []chan *pb.BroadcastRoomReq
        routinesNum uint64
        ipCnts map[string]int32
    }