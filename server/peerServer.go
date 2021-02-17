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

	return nil, nil
}

/**
 * @Description: 处理客户端发送过来的日志信息
 * @Param 上下文信息、添加日志信息请求
 * @return 添加日志信息的回复
 * @author zhangruiyuan
 * @date 2021/1/16 3:40 下午
 */
func (s *server) AppendEntries(ctx context.Context, in *proto.AppendEntriesRequest) (*proto.AppendEntriesReply, error) {

	return nil, nil
}

/**
 * @Description: 完成对server对象的初始化
 * @author zhangruiyuan
 * @date 2021/1/16 8:25 下午
 */
func NewServer() *server {
	// 1. 初始化server文件
	s := &server{
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
