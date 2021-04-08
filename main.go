package main

import (
	"fmt"
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
		//fmt.Println("正在执行命令：", command)
		switch command {
		case "currentLeader":
			fmt.Println(s.CurrentLeader)

		case "addLog":
			fmt.Printf("输入你在系统中留下的痕迹:")
			msg := ""
			fmt.Scanln(&msg)
			// 将日志的信息追加到s.messages中
			s.PushMassage(msg)

		case "showLog":
			s.ShowMassage()

		case "showPrevLogIndex":
			s.MutexPeers.Lock()
			for _, p := range s.Peers {
				p.Lock()
				defer p.Unlock()
				fmt.Printf("%s节点的前项日志索引为:%d,是否经过询问:%s\n", p.Name, p.PrevLogIndex, p.IfAsk)
			}
			s.MutexPeers.Unlock()

		default:
			fmt.Println()
			fmt.Println("共识节点的命令如下")
			fmt.Println("currentLeader      展示系统当前的leader")
			fmt.Println("addLog             测试追加日志")
			fmt.Println("showLog            展示系统中的日志")
			fmt.Println("showPrevLogIndex   展示系统当前最大的日志编号")
			fmt.Println()
		}

	}

}
