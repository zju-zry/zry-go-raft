/**
 * @note:
 * @author: zhangruiyuan
 * @date:2021/2/18
**/
package server

import (
	"encoding/json"
	"fmt"
	"time"
	"zry-raft/enum"
	pb "zry-raft/proto"
)

/**
 * @Description: 向日志文件中追加一条日志（命令行操作）
 * @author zhangruiyuan
 * @date 2021/2/18 8:25 下午
 */
func (s *server) PushMassage(data string) {

	// 1. 当前节点是leader节点的时候执行下面的操作
	if s.Name == s.CurrentLeader {
		s.DoPushMassage(data)
	} else {
		// 2. 当前节点不是leader节点的时候
		// 需要将该日志信息发送到leader中，并获取leader的返回值进行显示。
		// 日志的内容同步是另一个自动的步骤，在这里不需要等待同步leader中的数据成功之后再进行显示
		// 日志的同步在append中进行
		fmt.Println("正在向leader节点发送数据")
		for _, p := range s.Peers {
			if p.Name == s.CurrentLeader {
				res := p.PushMessageToLeader(&pb.PushMessageRequest{
					CommitName: s.Name,
					Data:       data,
				})
				if res.GetStatus() == enum.Success {
					fmt.Println("添加数据成功")
					fmt.Printf("日志编号：%d\n", res.Index)
					fmt.Printf("用户名字：%s\n", s.Name)
					fmt.Printf("leader名字：%s \n", res.LeaderName)
					fmt.Printf("任期：%d\n", res.Term)
					fmt.Printf("内容：%s\n", data)
					fmt.Printf("时间：%s\n", res.CommitTime)
				} else {
					fmt.Println("添加数据失败，可能是因为你找的这个节点并非leader节点")
				}
				break
			}
		}
	}
}

/**
 * @Description: 真正的提交message信息（ 当前节点为leader节点的时候本函数直接使用）
 * @author zhangruiyuan
 * @date 2021/2/19 12:06 上午
 */
func (s *server) DoPushMassage(data string) {
	s.MutexMessages.Lock()
	defer s.MutexMessages.Unlock()
	s.Messages = append(s.Messages, message{
		Index:      int64(len(s.Messages)),
		CommitName: s.Name,
		LeaderName: s.CurrentLeader,
		Term:       s.GetCurrentTerm(),
		Data:       data,
		CommitTime: time.Now(),
	})
}

/**
 * @Description: 展示日志文件（命令行操作）
 * @author zhangruiyuan
 * @date 2021/2/18 8:25 下午
 */
func (s *server) ShowMassage() {
	s.MutexMessages.Lock()
	defer s.MutexMessages.Unlock()
	fmt.Printf("系统中存在的日志数量: %d \n", len(s.Messages))
	fmt.Printf("下一个要记录的日志索引: %d \n", len(s.Messages)+1)
	fmt.Printf("系统中存在的日志信息如下：\n")
	for _, m := range s.Messages {
		fmt.Printf("日志编号：%d\n", m.Index)
		fmt.Printf("用户名字：%s\n", m.CommitName)
		fmt.Printf("leader名字：%s \n", m.LeaderName)
		fmt.Printf("任期：%d\n", m.Term)
		fmt.Printf("内容：%s\n", m.Data)
		fmt.Printf("时间：%s\n", m.CommitTime)
	}
}

/**
 * @Description: 获取日志数组的json字符串
 * @Param 以prevIndex开始的所有日志信息
 * @author zhangruiyuan
 * @date 2021/2/18 9:17 下午
 */
func (s *server) GetMassagesJsonData(prevIndex int64) string {
	s.MutexMessages.Lock()
	defer s.MutexMessages.Unlock()
	data, err := json.Marshal(s.Messages[prevIndex:])
	if err != nil {
		fmt.Printf("日志在json构造中出错\n")
		return ""
	}
	return string(data)
}

/**
 * @Description: 向本地插入获得的massage串信息
 * @Param 日志数组的json字符串
 * @author zhangruiyuan
 * @date 2021/2/18 9:19 下午
 */
func (s *server) InsertMassagesJsonData(data string) {
	s.MutexMessages.Lock()
	defer s.MutexMessages.Unlock()
	newMessage := []message{}
	err := json.Unmarshal([]byte(data), &newMessage)
	if err != nil {
		fmt.Printf("日志在json解析中出错\n")
	} else if len(newMessage) != 0 {
		fmt.Printf("成功同步日志信息\n")
		s.Messages = append(s.Messages, newMessage...)
	}
}

/**
 * @Description: 获取日志数组的高度
 * @Param 以prevIndex开始的所有日志信息
 * @author zhangruiyuan
 * @date 2021/2/18 9:17 下午
 */
func (s *server) GetMassagesHeight() int {
	s.MutexMessages.Lock()
	defer s.MutexMessages.Unlock()
	return len(s.Messages)
}
