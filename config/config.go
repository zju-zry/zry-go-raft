package config

import (
	"flag"
	"fmt"
	"time"
)

/**
 *	静态变量
 */
const (
	// 超时时，时间设立的单位1s
	TimeoutDurationUnit = 1000000000
	HeartbeatInterval   = 50 * time.Millisecond
	ElectionTimeout     = 10000 * time.Millisecond

	// 当前节点的三种状态
	StateFollower  = "follower"
	StateCandidate = "candidate"
	StateLeader    = "leader"
)

/**
 * @Description: 配置文件类，保存有当前的配置信息
 * @author zhangruiyuan
 * @date 2021/1/15 2:17 下午
 */
type Config struct {
	// 本节点的名称
	Name string
	// 每个节点总是处于以下状态的一种：follower、candidate、leader
	State string
	Ip    string
	Port  string

	//Raft协议关键概念: 当前任期，每个term内都会产生一个新的leader
	CurrentTerm int64
}

/**
 * @Description: 启动一个初始化配置，以follow节点的身份启动
 * @author zhangruiyuan
 * @date 2021/1/15:1:57
 */
func NewConfig(name, ip, port string) *Config {
	return &Config{
		Name:        name,
		State:       StateFollower,
		Ip:          ip,
		Port:        port,
		CurrentTerm: 0,
	}
}

/**
 *	创建一个全局的配置对象，供当前系统使用
 */
var Myconfig *Config

/**
 * @Description: 初始化配置文件
 * @Param
 * @return
 * @author zhangruiyuan
 * @date 2021/1/15 2:35 下午
 */
func init() {
	// 1. 解析调用的时候传入的配置信息
	fmt.Println("开始初始化配置文件信息 --->")
	fmt.Println("解析参数 --->")
	var name, ip, port string
	flag.StringVar(&ip, "ip", "127.0.0.1", "指定启动时的ip地址，默认为127.0.0.1")
	flag.StringVar(&name, "name", "org0.node0.order0", "本节点的名称，默认为org0.node0.order0")
	flag.StringVar(&port, "port", "7060", "端口号，默认为7060")
	flag.Parse()
	Myconfig = NewConfig(name, ip, port)
	fmt.Println("设置当前端口号：", port)

}
