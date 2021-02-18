### 一、项目文件介绍
#### 1. work.md
`work.md`文件是我正在开发的工作记录，其中有详细的开发记录。
   
#### 2. 概念解释

##### 2.1 循环
   我一般在注释中使用【循环】来表示`server.loop()`的函数内执行的内容

### 二、配置文件信息
1. 端口号设置
   + peer 节点705x
   + order 节点706x
   + 客户端 节点 707x
   + 链码端 708x
   
2. 倒计时设置
   + HeartbeatInterval = 50 * time.Millisecond
   + ElectionTimeout   = 200 * time.Millisecond
   
3. 日志设置
   + 第一个日志的编号为1