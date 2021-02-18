/**
 * @note: 本节点处理循环内容
 * @author: zhangruiyuan
 * @date:2021/2/18
**/
package server

import (
	"fmt"
	"zry-raft/enum"
	"zry-raft/proto"
	"zry-raft/util"
)

/**
 * @Description: 循环主体
 * @author zhangruiyuan
 * @date 2021/2/18 1:10 上午
 */
func (s *server) Loop() {
	// 循环应该处理两个内容
	// 1 处理倒计时结束成为候选节点
	// 2 处理倒计时结束之后向其他节点发送append信息
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
		fmt.Println(s.State)
		switch s.GetState() {
		case StateFollower:
			s.followerLoop()
		case StateCandidate:
			s.candidateLoop()
		case StateLeader:
			s.leaderLoop()
		}
	}
}

func (s *server) followerLoop() {
	s.FTimeoutChan = util.AfterBetween(ElectionTimeout, ElectionTimeout*2)
	for s.GetState() == StateFollower {
		//fmt.Println("正在执行follow循环")
		select {
		case <-s.FTimeoutChan:
			if true {
				s.SetState(StateCandidate)
				fmt.Println("接收消息超时，成为候选者")
			} else {
				fmt.Println("接收消息超时，继续等待连接")
			}

		default:

		}
	}
}

func (s *server) candidateLoop() {
	doVote := true
	votesGranted := 1
	newTerm := s.CurrentTerm
	var respChan chan *proto.VoteReply
	CTimeoutChan := util.AfterBetween(ElectionTimeout, ElectionTimeout*2)

	for s.GetState() == StateCandidate {
		//fmt.Println("正在执行candidate循环")

		if doVote {
			newTerm++
			respChan = make(chan *proto.VoteReply, len(s.Peers))
			for _, peer := range s.Peers {
				go peer.SendVoteRequest(&proto.VoteRequest{
					Term:          newTerm,
					LastLogIndex:  0,
					LastLogTerm:   0,
					CandidateName: s.Name,
				}, respChan)
				fmt.Println("发送请求投票信息到：", peer.ConnectionString)
			}
			// 修改标记是否拉票变量
			votesGranted = 1
			doVote = false
			CTimeoutChan = util.AfterBetween(ElectionTimeout, ElectionTimeout*2)
		}
		select {
		case res := <-respChan:
			if res.Term == newTerm && res.VoteGranted == enum.Success {
				votesGranted++
				fmt.Println("收到了一封没有过期的选票，当前收到的票数为", votesGranted)
			}
			if votesGranted > (len(s.Peers)+1)/2 {
				fmt.Println("收到了足够的票数，成为leader节点，当前任期为", newTerm)
				s.SetState(StateLeader)
				s.SetCurrentTerm(newTerm)
				s.CurrentLeader = s.Name
			}

		case <-CTimeoutChan:
			// 选举超时，进行下一轮选举
			fmt.Println("选举超时，进行下一轮选举")
			doVote = true

		default:

		}
	}
}

func (s *server) leaderLoop() {
	// 成功竞争为领导者，先向大家询问一下大家都是什么状态,默认是false
	for _, p0 := range s.Peers {
		p0.Lock()
		defer p0.Unlock()
		p0.IfAsk = false
	}
	// 向其他节点发送append消息
	var respChan chan *proto.AppendEntriesReply
	LTimeoutChan := util.AfterBetween(HeartbeatInterval, HeartbeatInterval*2)
	for s.GetState() == StateLeader {
		select {
		case <-LTimeoutChan:
			respChan = make(chan *proto.AppendEntriesReply, len(s.Peers))
			//fmt.Println("我在执行leader事情")
			// leader节点应该向其他节点发送数据列表信息
			for _, p := range s.Peers {
				p.Lock()
				defer p.Unlock()
				jdata := "[]"
				// 经过询问之后再发送相应的数据
				if p.IfAsk {
					jdata = s.GetMassagesJsonData(p.PrevLogIndex)
				}
				go p.sentHeartbeat(&proto.AppendEntriesRequest{
					Term:       s.GetCurrentTerm(),
					LeaderName: s.Name,
					Massages:   jdata,
				}, respChan)
			}
			LTimeoutChan = util.AfterBetween(HeartbeatInterval, HeartbeatInterval*2)

		case res := <-respChan:
			// 接收到其他节点发送的返回值，用以更新本机记录的pervLogIndex等信息
			//fmt.Printf("接收到其他节点发送的返回值，用以更新本机记录的pervLogIndex等信息\n")
			//fmt.Printf(res.PeerName,res.PrevLogIndex,res.Term,"\n")
			for _, p := range s.Peers {
				p.Lock()
				defer p.Unlock()
				if p.Name == res.PeerName {
					p.PrevLogIndex = res.PrevLogIndex
					p.IfAsk = true
					fmt.Printf("将%s的状态修改了.\n", p.Name)
					break
				}
			}

		default:

		}
	}

}
