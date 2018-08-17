**1.集中式的互斥**
```go
//server
while true:
   m = Receive()
   If m == (Request, i):
     If Available():
	   Send (Grant) to I
   else:
     Add i to Q
   If m == (Release)&&!empty(Q):
     Remove ID j from Q
     Send (Grant) to j  重新分配
//client acquire
Send (Request, i) to coordinator
Wait for reply
//client replay
Send (Release) to coordinator
```

**2.有序的多播**
based lamport

**3.分布式互斥**
