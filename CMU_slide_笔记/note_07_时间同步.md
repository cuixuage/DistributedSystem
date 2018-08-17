**1.时间同步的方法**
NTP algorithm==timestamp
udp传递发送.接收时间戳
RTT = (client_wait_time - server_proc处理_time)/2

**2.Lamport clocks**
logical clocks
缺点: L(e) < L(e')并不是意味 e一定发生在e'之前
例如 有可能是并发的
**3.vector clocks**
