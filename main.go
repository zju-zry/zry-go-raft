package main

// 1 解析、加载配置文件
import (
	"fmt"
	_ "zry-raft/config"
	"zry-raft/server"
)

func main() {

	// 2 进入到follow的状态
	// 3 启动后台peer server，不断地接收来自外部的请求
	s := server.NewServer()
	go s.Start()
	// 4 启动控制循环
	go s.Loop()
	// 5 启动控制台的输入命令
	command := ""
	for {
		fmt.Scanln(&command)
		fmt.Println("正在执行命令：", command)
	}

}
