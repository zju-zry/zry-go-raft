/**
 * @note: grpc测试
 * @author: zhangruiyuan
 * @date:2021/1/15
**/
package main

import (
	"context"
	"log"
	"net"
	"zry-raft/experiment/grpc/proto"

	"google.golang.org/grpc"
)

const port = "127.0.0.1:7061"

// server用来实现节点的服务：PeerServer.
type server struct {
	proto.UnimplementedPeerServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SendVoteRequest(ctx context.Context, in *proto.VoteRequest) (*proto.VoteReply, error) {
	log.Printf("Received: %v", in.GetCandidateName())
	return &proto.VoteReply{ Term: 1 , VoteGranted: 200}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterPeerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}