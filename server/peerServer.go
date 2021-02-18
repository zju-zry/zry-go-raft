/**
 * @note: 节点的服务端，将本节点作为服务响应请求
 * @author: zhangruiyuan
 * @date:2021/1/15
**/
package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"zry-raft/enum"
	"zry-raft/proto"
	"zry-raft/util"

	"google.golang.org/grpc"
)

/**
 *	静态变量
 */
const (
	// 心脏超时
	HeartbeatInterval = 50 * time.Millisecond
	// 选举超时
	ElectionTimeout = 10000 * time.Millisecond

	// 当前节点的三种状态
	StateFollower  = "follower"
	StateCandidate = "candidate"
	StateLeader    = "leader"
)

/**
 * @Description: 消息对象
 * @author zhangruiyuan
 * @date 2021/2/18 8:03 下午
 */
type message struct {
	// 交易的编号
	Index int64
	// 提交该交易的用户名字
	CommitName string
	// 提交该交易的leader名字
	LeaderName string
	// 提交该笔交易的任期
	Term int64
	// 该笔交易的内容
	Data string
	// 提交该笔交易的时间
	CommitTime time.Time
}

/**
 * @Description: server用来实现节点的服务：PeerServer.
 * @author zhangruiyuan
 * @date 2021/1/16 2:57 下午
 */
type server struct {
	MutexPeers sync.Mutex
	Peers      []Peer
	// 本节点的名称
	Name string
	// 每个节点总是处于以下状态的一种：follower、candidate、leader
	MutexState sync.Mutex
	State      string
	Ip         string
	Port       string

	// Raft协议关键概念: 当前任期，每个term内都会产生一个新的leader
	// 任期改变的条件
	// 1. 收到新的一轮投票（投票的任期高于自己记录的任期，即要求更新的任期需要大于我记录的任期我才会更新）
	// 2. 收到心脏数据

	MutexCurrentTerm sync.Mutex
	CurrentTerm      int64
	CurrentLeader    string

	// 超时状态
	FTimeoutChan <-chan time.Time

	// 系统操作的信息列表
	MutexMessages sync.Mutex
	Messages      []message
}

/**
 * @Description: 重置跟随者状态
 * @date 2021/2/18 7:11 下午
 */
func (s *server) ResetFollowerState() {
	s.FTimeoutChan = util.AfterBetween(ElectionTimeout, ElectionTimeout*2)
	s.SetState(StateFollower)
}

/**
 * @Description: 处理客户端发来的请求投票的请求
 * @Param 上下文信息、投票的请求
 * @return 投票回复的信息
 * @author zhangruiyuan
 * @date 2021/1/16 2:55 下午
 */
func (s *server) SendVoteRequest(ctx context.Context, in *proto.VoteRequest) (*proto.VoteReply, error) {
	if in.GetTerm() > s.GetCurrentTerm() {
		fmt.Println("收到一条请求投票信息, 本节点转为follow节点", in.CandidateName, "任期为", in.Term)
		s.SetCurrentTerm(in.Term)
		s.ResetFollowerState()
		return &proto.VoteReply{
			Term:        s.GetCurrentTerm(), // 当前最高的任期
			VoteGranted: enum.Success,       // 同意该请求
		}, nil
	} else {
		fmt.Println("收到一条非法的请求投票信息, 本节点转为follow节点", in.CandidateName, "任期为", in.Term, "已经过时")
		return &proto.VoteReply{
			Term:        s.GetCurrentTerm(), // 当前最高的任期
			VoteGranted: enum.Fail,          // 不同意该请求
		}, nil
	}
}

/**
 * @Description: 处理客户端发送过来的日志信息
 * @Param 上下文信息、添加日志信息请求
 * @return 添加日志信息的回复
 * @author zhangruiyuan
 * @date 2021/1/16 3:40 下午
 */
func (s *server) AppendEntries(ctx context.Context, in *proto.AppendEntriesRequest) (*proto.AppendEntriesReply, error) {
	if in.GetTerm() >= s.GetCurrentTerm() {
		s.SetCurrentTerm(in.GetTerm())
		s.CurrentLeader = in.LeaderName
		s.ResetFollowerState()
		oldSize := s.GetMassagesHeight()
		s.InsertMassagesJsonData(in.Massages) //接收发来的数据
		newSize := s.GetMassagesHeight()
		if newSize != oldSize {
			//s.ShowMassage()
			fmt.Printf("收到leader节点%s发来的追加日志信息，本节点的状态立即转为follower。\n", in.GetLeaderName())
		}
	}
	s.MutexMessages.Lock()
	defer s.MutexMessages.Unlock()
	return &proto.AppendEntriesReply{
		Term:         s.GetCurrentTerm(),
		PrevLogIndex: int64(len(s.Messages)),
		PeerName:     s.Name,
	}, nil
}

/**
 * @Description: 向leader节点追加消息
 * @author zhangruiyuan
 * @date 2021/2/19 12:00 上午
 */
func (s *server) PushMessage(ctx context.Context, in *proto.PushMessageRequest) (*proto.PushMessageReply, error) {
	if s.GetState() == StateLeader {
		s.MutexMessages.Lock()
		defer s.MutexMessages.Unlock()
		nt := time.Now()
		res := &proto.PushMessageReply{
			Index:      int64(len(s.Messages)),
			LeaderName: s.Name,
			Term:       s.CurrentTerm,
			CommitTime: nt.String(),
			Status:     enum.Success,
		}
		// 追加数据
		s.Messages = append(s.Messages, message{
			Index:      res.Index,
			CommitName: in.CommitName,
			LeaderName: res.LeaderName,
			Term:       res.Term,
			Data:       in.Data,
			CommitTime: nt,
		})
		return res, nil
	} else {
		return &proto.PushMessageReply{Status: enum.Fail}, nil
	}
}

/**
 * @Description: 完成对server对象的初始化
 * @author zhangruiyuan
 * @date 2021/1/16 8:25 下午
 */
func NewServer() *server {
	// 0. 解析配置文件
	var name, ip, port string
	flag.StringVar(&ip, "ip", "127.0.0.1", "指定启动时的ip地址，默认为127.0.0.1")
	flag.StringVar(&name, "name", "org0.node0.order0", "本节点的名称，默认为org0.node0.order0")
	flag.StringVar(&port, "port", "7060", "端口号，默认为7060")
	flag.Parse()

	// 1. 初始化server文件
	s := &server{
		Name:        name,
		State:       StateFollower,
		Ip:          ip,
		Port:        port,
		CurrentTerm: 0,
		// 配置其他节点的信息（其中包含有本机的信息）
		Peers: []Peer{
			{
				Name:             "org1.node1.order1",
				ConnectionString: "127.0.0.1:7061",
				PrevLogIndex:     0,
				IfAsk:            false,
			},
			{
				Name:             "org1.node2.order1",
				ConnectionString: "127.0.0.1:7062",
				PrevLogIndex:     0,
				IfAsk:            false,
			},
			{
				Name:             "org2.node1.order1",
				ConnectionString: "127.0.0.1:7063",
				PrevLogIndex:     0,
				IfAsk:            false,
			},
			{
				Name:             "org2.node2.order1",
				ConnectionString: "127.0.0.1:7064",
				PrevLogIndex:     0,
				IfAsk:            false,
			},
		},
		Messages: make([]message, 0), // 日志切片的初始化
	}
	// 2. 根据配置的字符串和ip信息，移除一下本节点的peer
	for i, p := range s.Peers {
		if p.ConnectionString == s.Ip+":"+s.Port {
			s.Peers = append(s.Peers[:i], s.Peers[i+1:]...)
		}
	}
	fmt.Println("已知的其他节点的信息为：", s.Peers)
	// 3. 返还字符串的信息
	fmt.Println("初始化server中变量信息完成")
	return s
}

/**
 * @Description: 启动对其他peer节点的服务
 * @author zhangruiyuan
 * @date 2021/1/16 3:41 下午
 */
func (s *server) Start() {
	defer fmt.Println("节点服务关闭成功")
	// 1. 启动监听grpc服务的信息
	lis, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	sr := grpc.NewServer()
	proto.RegisterPeerServer(sr, s)
	if err := sr.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/**
 * @Description: 返回成为leader需要的选票的数量
 * @author zhangruiyuan
 * @date 2021/1/17 2:36 下午
 */
func (s *server) QuorumSize() int {
	return (len(s.Peers) / 2) + 1
}

/**
 * @Description: 状态的获取函数
 * @author zhangruiyuan
 * @date 2021/2/18 7:10 下午
 */
func (s *server) GetState() string {
	s.MutexState.Lock()
	defer s.MutexState.Unlock()
	return s.State
}

/**
 * @Description: 状态的设置函数
 * @author zhangruiyuan
 * @date 2021/2/18 7:10 下午
 */
func (s *server) SetState(state string) {
	s.MutexState.Lock()
	defer s.MutexState.Unlock()
	s.State = state
}

/**
 * @Description: 设置当前任期
 * @author zhangruiyuan
 * @date 2021/2/18 7:19 下午
 */
func (s *server) SetCurrentTerm(t int64) {
	s.MutexCurrentTerm.Lock()
	defer s.MutexCurrentTerm.Unlock()
	s.CurrentTerm = t
}

/**
 * @Description: 获取当前任期
 * @author zhangruiyuan
 * @date 2021/2/18 7:19 下午
 */
func (s *server) GetCurrentTerm() int64 {
	s.MutexCurrentTerm.Lock()
	defer s.MutexCurrentTerm.Unlock()
	return s.CurrentTerm
}
