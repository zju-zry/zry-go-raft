/**
 * @note: 测试peerServer的投票请求连通情况
 * @author: zhangruiyuan
 * @date:2021/1/15
**/
package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
	pb "zry-raft/proto"
)

const (
	address     ="localhost:7061"
)

func main() {
	// 创建一条到服务端的链接
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("链接不上服务器: %v", err)
	}
	defer conn.Close()
	c := pb.NewPeerClient(conn)

	// 与服务端通信，并将返回结果进行打印
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	r, err := c.SendVoteRequest(ctx,&pb.VoteRequest{
		Term: 1,
		LastLogIndex: -1,
		LastLogTerm: -1,
		CandidateName: "zhangruiyuan",
	})
	if err != nil {
		log.Fatalf("完成投票请求中出现错误: %v", err)
	}
	log.Printf("Term: %d，VoteGranted: %d ", r.GetTerm(),r.GetVoteGranted())
}
