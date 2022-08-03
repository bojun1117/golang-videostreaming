package main //流量控制

import "log"

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool { //判斷使用者是否過多
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation.")
		return false
	}
	cl.bucket <- len(cl.bucket) + 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("New connction coming: %d", c)
}
