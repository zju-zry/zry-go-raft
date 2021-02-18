package main

import (
	"fmt"
	"zry-raft/proto"
	"zry-raft/server"
)

func main() {
	// 1 解析、加载配置文件
	// 2 进入到follow的状态
	// 3 启动后台peer server，不断地接收来自外部的请求
	s := server.NewServer()
	go s.Start()
	// 4 启动循环
	go s.Loop()
	// 5 启动控制台的输入命令
	command := ""
	for {
		fmt.Scanln(&command)
		fmt.Println("正在执行命令：", command)
		switch command {
		case "goLeader":
			fmt.Println("成为leader ing ... ")
			respChan := make(chan *proto.VoteReply, len(s.Peers))
			s.State = server.StateCandidate
			for _, peer := range s.Peers {
				go peer.SendVoteRequest(&proto.VoteRequest{
					Term:          s.CurrentTerm + 1,
					LastLogIndex:  0,
					LastLogTerm:   0,
					CandidateName: s.Name,
				}, respChan)
				fmt.Println("发送请求投票信息到：", peer.ConnectionString)
			}
		case "CurrentLeader":
			fmt.Println(s.CurrentLeader)

		}
	}

}
