DFS = distributed file systems
NFS = network ~ 使用RPC作为client server通信
AFS = andrew ~ 弱一致性 ==》2.修改回写

CAP权衡 = consistency availability partitoin-resilience分区弹性

**1.NFS分布式系统 文件修改handle failures**
client cache
e.g. 当文件is closed,所有被修改的blocks推送给server，收到server replay否则得到failures

AFS作者认为 比如数据库之类的高度并发的共享访问的应用需要不同的模型处理 和dfs设计理念不同


**2.file access consistency**
2.1 unix local file systems只是使用sequential consistency 顺序一致性==> kernel lock the file vnode
同时NFS以及AFS都都没有提供对于file的并发控制

2.2 writeback dirty data
假设write confilcts rare少数情况
NFS、AFS==》中央服务器的读写瓶颈

2.3 out-of-date file blocks

**3.name space、user access**
AAA system kerberos
ticket server 分别给clinet file_server发送key 在进行对应验证

**4.coda support短暂断开连接的dfs**
client断开之前获取obj的lock,同时在limited time恢复,file change会被server整合   ==>  file未被close一直缓存在server 超过cache size会被Equilibrium均衡处理
coda process重新连接后=每一个volume卷执行日志重播算法log replay algorithm

**5.LBFS 弱连接的dfs**
rabin哈希只进行更改文件的部分内容(基于内容的块定义)，减少文件的传输时的带宽占用
