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
   + 某节点初次成为leader节点，leader节点会首先向所有的peer询问需要第几个日志，
     然后再传输数据，防止leader节点的大量没有必要的数据传输
   

### 三、项目启动
1. 测试启动
   
   请按照下面这张图片进行配合
   ![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219005416.png)
   
   即在本机中使用如下四个端口：
   + 127.0.0.1:7061
   + 127.0.0.1:7062
   + 127.0.0.1:7063
   + 127.0.0.1:7064
   
2. 项目配置
   
   2.1 配置节点的信息
      请在文件`zry-raft/server/peerServer.go`中配置节点的信息，并按照3.1样例对您自己的配置进行修改
   ![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219010120.png)
   

3. 启动演示
   
   节点启动之后，若一半以上的节点能够合作，能够达成leader的生成。如下图，order[7062]节点成为leader状态，
   ![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219010451.png)
   
   其他follow节点一般在收到leader节点发来的请求投票信息后一直作为follow节点（除非leader节点崩掉）。
   ![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219010658.png)
   
   leader节点崩坏后效果如下：
   ![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219010843.png)
   
   
4. 操作演示
   
```shell
   goLeader            主动申请成为leader节点（代码有问题，还在测试中）
   currentLeader       展示系统当前的leader节点
   addLog              在系统中添加日志信息（若本节点为leader则直接添加，若本节点是普通节点就需要向leader转发）
   showLog             展示系统中的日志信息
   showPrevLogIndex    展示系统中其他节点当前需要的日志索引（leader中使用有意义）
```

currentLeader:

![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219011836.png)

addLog:

![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219011901.png)

showLog:

![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219011946.png)

showPrevLogIndex:

![](https://zhangruiyuan.oss-cn-hangzhou.aliyuncs.com/picGo/images/20210219012104.png)

核心功能为日志的提交和查看，其余可以不用关心。