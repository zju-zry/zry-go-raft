/**
 * @note:
 * @author: zhangruiyuan
 * @date:2021/1/15
**/
package main

import "fmt"

func main()  {
	//leaderHealth := 10
	//tick := time.Tick(1 * time.Second)
	//for {
	//	select{
	//	case <-tick:
	//		leaderHealth--
	//		if leaderHealth<0{
	//			fmt.Println("系统进入候选者状态")
	//		}
	//	}
	//}

	ps := []string{ "a", "b", "c"}
	for i,v := range ps{
		if v == "b"{
			ps = append(ps[:i],ps[i+1:]...)
		}
	}
	fmt.Println(ps)
}
