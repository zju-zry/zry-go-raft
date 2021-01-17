/**
 * @note: 节点的服务端，将本节点作为服务响应请求
 * @author: zhangruiyuan
 * @date:2021/1/15
**/
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"zry-raft/common"
	"zry-raft/config"
	"zry-raft/proto"

	"google.golang.org/grpc"
)

/**
 * @Description: server用来实现节点的服务：PeerServer.
 * @author zhangruiyuan
 * @date 2021/1/16 2:57 下午
 */
type server struct {
	Peers []Peer
	// server需要继承以下服务
	proto.UnimplementedPeerServer
	// 事件通道，之后收到的服务的信息在循环中进行处理
	c chan *common.Ev
	// 保护当前server节点信息？
	// sync包中的WaitGroup实现了一个类似任务队列的结构，你可以向队列中加入任务，任务完成后就把任务从队列中移除，如果队列中的任务没有全部完成，队列就会触发阻塞以阻止程序继续运
	// 这个东西的用处我当前还不是很清楚
	//routineGroup sync.WaitGroup
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
 * @Description: 处理客户端发来的请求投票的请求
 * @Param 上下文信息、投票的请求
 * @return 投票回复的信息
 * @author zhangruiyuan
 * @date 2021/1/16 2:55 下午
 */
func (s *server) SendVoteRequest(ctx context.Context, in *proto.VoteRequest) (*proto.VoteReply, error) {
	/**
	处理投票请求的流程。
	丢给主循环去做这个事情？ 这样可以对系统不同状态下的信息进行统一的管理
	那怎么在丢给主循环做这个事情之后获得返回值呢？
	可以使用chan去获取返回值，但是这就要求在chan中必须返回值
	经过一天的卡点，我决定使用最新的方式实现这么个功能。
	*/
	log.Printf("Received: %v", in.GetCandidateName())
	// 1. 创建一个新的通道
	ev := &common.Ev{Target: in, ReturnC: make(chan bool)}
	fmt.Println(ev)
	// 2. 将这个通道传入到server端
	s.c <- ev
	fmt.Println("已经将请求事件告知了循环")
	// 3. 等待循环的处理结果
	<-ev.ReturnC
	fmt.Println("循环体已经执行结束啦，开始执行server的返回")
	// 4. 将处理好的结果进行一个返还
	return ev.ReturnValue.(*proto.VoteReply), nil

}

/**
 * @Description: 处理客户端发送过来的日志信息
 * @Param 上下文信息、添加日志信息请求
 * @return 添加日志信息的回复
 * @author zhangruiyuan
 * @date 2021/1/16 3:40 下午
 */
func (s *server) AppendEntries(ctx context.Context, in *proto.AppendEntriesRequest) (*proto.AppendEntriesReply, error) {
	log.Printf("在AppendEntries中收到的leader节点的任期信息: %v", in.GetTerm())
	// 1. 创建一个新的通道
	ev := &common.Ev{Target: in, ReturnC: make(chan bool)}
	// 2. 将这个通道传入到server端
	s.c <- ev
	// 3. 等待循环的处理结果
	<-ev.ReturnC
	// 4. 将处理好的结果进行一个返还
	return ev.ReturnValue.(*proto.AppendEntriesReply), nil
}

/**
 * @Description: 完成对server对象的初始化
 * @author zhangruiyuan
 * @date 2021/1/16 8:25 下午
 */
func NewServer() *server {
	// 1. 初始化server文件
	s := &server{
		c: make(chan *common.Ev, 256),
		// 配置其他节点的信息（其中包含有本机的信息）
		Peers: []Peer{
			{
				Name:             "org1.node1.order1",
				ConnectionString: "127.0.0.1:7061",
			},
			{
				Name:             "org1.node2.order1",
				ConnectionString: "127.0.0.1:7062",
			},
			{
				Name:             "org2.node1.order1",
				ConnectionString: "127.0.0.1:7063",
			},
			{
				Name:             "org2.node2.order1",
				ConnectionString: "127.0.0.1:7064",
			},
		},
	}
	// 2. 根据配置的字符串和ip信息，移除一下本节点的peer
	for i, p := range s.Peers {
		if p.ConnectionString == config.Myconfig.Ip+":"+config.Myconfig.Port {
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
	lis, err := net.Listen("tcp", ":"+config.Myconfig.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	sr := grpc.NewServer()
	proto.RegisterPeerServer(sr, s)
	if err := sr.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
