/**
 * @note: 实现随机数倒计时提醒的功能
 * @author: zhangruiyuan
 * @date:2021/1/16
**/
package util

import (
	"math/rand"
	"time"
)

/**
 * @Description: 实现倒计时随机提醒的功能
 * @Param 时间从min到max，这两个字段代表两个时间
 * @return 时间提醒的通道
 * @author zhangruiyuan
 * @date 2021/1/16 11:59 下午
 */
func afterBetween(min time.Duration, max time.Duration) <-chan time.Time {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	d, delta := min, (max - min)
	if delta > 0 {
		d += time.Duration(rand.Int63n(int64(delta)))
	}
	return time.After(d)
}
