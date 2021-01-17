# 本文记录，我在开发中遇到的问题

## 1. import循环

问题发生的原因：
我在config中定义了peers，初始化了其他节点的信息，所以config包引用了server包中的peer，
在server包中使用了config包的配置信息。
所以造成了引用循环。

所以我最终的解决办法就是：
在config中定义全局使用的配置文件，在server中定义一些服务，
服务使用config中的定义信息，但是在config中不会再包含server的结构。

## 2. 结构体实现的方法不与结构体一个包，不能引用

问题发生的原因：
我打算在common包中定义peer、server的结构体，
然后在server包中进行真正的实现，但是最终没有成功。

解决办法：
放弃抵抗，直接就在server包中进行
