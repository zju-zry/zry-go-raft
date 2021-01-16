/**
 * @note:
 * @author: zhangruiyuan
 * @date:2021/1/15
**/
package common

import (
	"sync"
	"time"
)

// A peer is a reference to another server involved in the consensus protocol.
type Peer struct {
	// Name：peer的名称
	Name              string `json:"name"`
	// ConnectionString：peer的ip地址，形式为”ip:port”
	ConnectionString  string `json:"connectionString"`
	// prevLogIndex：这个很关键，记录了该peer的当前日志index，接下来leader将该index之后的日志继续发往该peer
	prevLogIndex      uint64
	// 停止通道？ 不知道是做什么用的
	stopChan          chan bool
	// 心脏跳动的时间间隔
	heartbeatInterval time.Duration
	//lastActivity：记录peer的上次活跃时间
	lastActivity      time.Time
	// 互斥保护当前Peer
	sync.RWMutex
}