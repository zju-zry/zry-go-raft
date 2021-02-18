/**
 * @note: 节点的客户端，将本节点作为客户机来对其他节点发起请求
 * @author: zhangruiyuan
 * @date:2021/1/17
**/
package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
	pb "zry-raft/proto"
)

/**
 * @Description: peer节点作为客户端的类
 * @author zhangruiyuan
 * @date 2021/1/17 1:29 下午
 */
type Peer struct {
	Name             string `json:"name"`             // Name：peer的名称
	ConnectionString string `json:"connectionString"` // ConnectionString：peer的ip地址，形式为”ip:port”
	PrevLogIndex     int64  // prevLogIndex：这个很关键，记录了该peer的当前日志index，接下来leader将该index之后的日志继续发往该peer

	IfAsk bool // 初为leader，要问一下需要从什么时候开始的log

	sync.RWMutex // 互斥保护当前Peer
}

/**
 * @Description:
 * @Param
 * @return
 * @author zhangruiyuan
 * @date 2021/1/17 1:40 下午
 */
func (p Peer) SendVoteRequest(request *pb.VoteRequest, respChan chan *pb.VoteReply) {
	// 创建一条到服务端的链接
	fmt.Println("连接请求：", p.ConnectionString)
	conn, err := grpc.Dial(p.ConnectionString, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("链接不上服务器: %v", err)
	}
	defer conn.Close()
	c := pb.NewPeerClient(conn)

	// 与服务端通信，并将返回结果进行打印
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	res, err := c.SendVoteRequest(ctx, request)
	if err != nil {
		log.Fatalf("完成投票请求中出现错误: %v", err)
	}
	// 将返回的结果放在这个channel中， 客户端接收到这个请求，并进行一个整合，在收到足够的响应之后，当前节点就可以成为一个leader
	respChan <- res
}

/**
 * @Description: 启动对这个节点的append提交
 * @author zhangruiyuan
 * @date 2021/1/17 5:34 下午
 */
func (p Peer) sentHeartbeat(req *pb.AppendEntriesRequest, respChan chan *pb.AppendEntriesReply) {

	//fmt.Println("本节点发送了一条心脏请求",p.ConnectionString)

	// 创建一条到服务端的链接
	conn, err := grpc.Dial(p.ConnectionString, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("链接不上服务器: %v", err)
	}
	defer conn.Close()
	c := pb.NewPeerClient(conn)

	// 与服务端通信，并将返回结果进行打印
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.AppendEntries(ctx, req)
	if err != nil {
		log.Fatalf("完成投票请求中出现错误: %v", err)
	}
	// 将返回的结果放在这个channel中， leader接收到这个返回值进行一个相应的处理
	respChan <- res
}
