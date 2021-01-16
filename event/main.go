package main

import (
	"fmt"
	"time"
)

const HELLO_WORLD = "helloWorld"

func main() {
	// 创建一个新的调度器
	dispatcher := NewEventDispatcher()
	dispatcher.name = "zhangruiyuan"
	// 创建一个新的监听器
	listener := NewEventListener(myEventListener)
	// 将监听器添加到调度器中 事件的名称是HELLO_WORLD
	dispatcher.AddEventListener(HELLO_WORLD, listener)
	// 停留两秒
	time.Sleep(time.Second * 2)
	//dispatcher.RemoveEventListener(HELLO_WORLD, listener)
	//启动一个新的事件，事件名称是HELLO_WORLD，事件中携带的数据为nil
	dispatcher.DispatchEvent(NewEvent(HELLO_WORLD, nil))
}

func myEventListener(event Event) {
	fmt.Println(event.Type, event.Object, event.Target.GetName())
}