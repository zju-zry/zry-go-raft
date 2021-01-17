/**
 * @note: 服务的循环实现。
 * @author: zhangruiyuan
 * @date:2021/1/17
**/
package server

import (
	"fmt"
	en "zry-raft/common/enum"
	"zry-raft/config"
	"zry-raft/proto"
	"zry-raft/util"
)

/**
 * @Description: server对象的主要循环
 * @author zhangruiyuan
 * @date 2021/1/16 8:41 下午
 */
func (s *server) Loop() {
	// 本段代码的主要逻辑为判断当前的状态信息，然后决定走哪一个循环
	// 以下是raft的状态转换图
	//--------------------------------------
	//               ________
	//            --|Snapshot|                 timeout
	//            |  --------                  ______
	// recover    |       ^                   |      |
	// snapshot / |       |snapshot           |      |
	// higher     |       |                   v      |     recv majority votes
	// term       |    --------    timeout    -----------                        -----------
	//            |-> |Follower| ----------> | Candidate |--------------------> |  Leader   |
	//                 --------               -----------                        -----------
	//                    ^          higher term/ |                         higher term |
	//                    |            new leader |                                     |
	//                    |_______________________|____________________________________ |

	defer fmt.Println("循环关闭")
	for {
		switch config.Myconfig.State {
		case config.StateFollower:
			s.followerLoop()
		case config.StateCandidate:
			s.candidateLoop()
		case config.StateLeader:
			s.leaderLoop()
			break
		}
	}
}

/**
 * @Description: 追随者的循环
 * @author zhangruiyuan
 * @date 2021/1/16 8:53 下午
 */
func (s *server) followerLoop() {
	// 用以判断选举超时的通道，这里是启动一个倒计时装置
	timeoutChan := util.AfterBetween(config.ElectionTimeout, config.ElectionTimeout*2)
	for config.Myconfig.State == config.StateFollower {
		// 用以判断这个本次循环中是否收到leader的消息
		update := false
		select {
		case e := <-s.c:
			switch req := e.Target.(type) {
			case *proto.AppendEntriesRequest:
				fmt.Println(req.GetLeaderName())
				update = true

			case *proto.VoteRequest:
				if req.Term > config.Myconfig.CurrentTerm {
					e.ReturnValue = &proto.VoteReply{Term: config.Myconfig.CurrentTerm, VoteGranted: en.Success}
				} else {
					e.ReturnValue = &proto.VoteReply{Term: config.Myconfig.CurrentTerm, VoteGranted: en.Fail}
				}
				e.ReturnC <- true
			}
		case <-timeoutChan:
			// 发现了超时，本节点就成为候选者
			config.Myconfig.State = config.StateCandidate
			fmt.Println("wo bian le ")
		}
		if update {
			// 注意，在这里重新进行倒计时，之前的那个通道就丢弃了。
			timeoutChan = util.AfterBetween(config.ElectionTimeout, config.ElectionTimeout*2)
		}
	}
}

/**
 * @Description: 候选者的循环
 * @author zhangruiyuan
 * @date 2021/1/16 8:53 下午
 */
func (s *server) candidateLoop() {
	fmt.Println("进入选举循环")
	// 用以判断选举超时的通道，这里是启动一个倒计时装置
	timeoutChan := util.AfterBetween(config.ElectionTimeout, config.ElectionTimeout*2)
	// 标记是否需要申请投票 true投票 false不投票
	doVote := true
	// 存放其他人的投票结果的通道
	var respChan chan *proto.VoteReply
	// 获取的投票数量
	votesGranted := 1
	// 申请leader的任期
	newTerm := config.Myconfig.CurrentTerm
	// 循环主要内容
	for config.Myconfig.State == config.StateCandidate {
		// 进行投票
		if doVote {
			// 首先是当前的任期++
			// 当前任期这个值是在追加数据成功之后进行修改的
			newTerm++
			respChan = make(chan *proto.VoteReply, len(s.Peers))
			for _, peer := range s.Peers {
				go peer.SendVoteRequest(&proto.VoteRequest{
					Term:          newTerm,
					LastLogIndex:  0,
					LastLogTerm:   0,
					CandidateName: config.Myconfig.Name,
				}, respChan)
			}
			// 修改标记是否拉票变量
			votesGranted = 1
			doVote = false
			timeoutChan = util.AfterBetween(config.ElectionTimeout, config.ElectionTimeout*2)
		}
		// 判断是不是满足了选票的数量
		if votesGranted >= s.QuorumSize() {
			config.Myconfig.CurrentTerm = newTerm
			config.Myconfig.State = config.StateLeader
		}
		select {
		case resp := <-respChan:
			// 节点同意了我在newTerm中成为leader的请求
			if resp.GetVoteGranted() == en.Success && resp.GetTerm() == newTerm {
				votesGranted++
				fmt.Println("当前申请的任期是", newTerm, "已经收到了", votesGranted)
			}
		case e := <-s.c:
			switch req := e.Target.(type) {
			case *proto.AppendEntriesRequest:
			case *proto.VoteRequest:
				fmt.Println("接收到别人的拉票信息")
				fmt.Println(req)
				if req.Term > newTerm {
					config.Myconfig.CurrentTerm = newTerm
					config.Myconfig.State = config.StateFollower
					fmt.Println("发现对面的任期比自己高，所以主动退位")
				}
			}
		// 时间到了进行投票
		case <-timeoutChan:
			//选举超时，进行重新选举
			doVote = true
			fmt.Println("选举超时，进行重新选举")
		}
	}
}

/**
 * @Description: 领导者的循环
 * @author zhangruiyuan
 * @date 2021/1/16 8:54 下午
 */
func (s *server) leaderLoop() {
	fmt.Println("我成为leader 任期是", config.Myconfig.CurrentTerm)
	// 开始向其他节点发送append请求。
	for _, p := range s.Peers {
		p.startHeartbeat()
	}
}
