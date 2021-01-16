/**
 * @note: 测试通道的用法
 * @author: zhangruiyuan
 * @date:2021/1/16
**/
package main

type user struct {
	name string
	age string
	c chan bool
}
func main()  {
	////c := make(chan bool)
	//var c chan bool
	//fmt.Println(c)
	//go func() {
	//	b := <- c
	//	fmt.Println(b)
	//	fmt.Println("hello")
	//}()
	//c <- false

	//c := make (chan *user)
	//go func() {
	//	u := <- c
	//	u.c <- true
	//	fmt.Println("循环体中已经执行完毕")
	//}()
	//u2 := &user{name: "zhangruiyuan",age: "12"}
	//u2.c = make(chan bool)
	//c <- u2
	//<- u2.c
	//fmt.Println("server已经处理完毕")

	u := user{name: "zhangruiyuan"}
	go func() {
		u.c <- false
	}()
	<-u.c
	for{}
}
