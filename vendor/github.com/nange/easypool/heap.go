package easypool

import (
	"container/heap"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type PriorityQueue []*PoolConn

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// we want to get the oldest item
	return pq[i].updatedTime.Sub(pq[j].updatedTime) < 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	pc := x.(*PoolConn)
	*pq = append(*pq, pc)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type heapPool struct {
	mu          sync.Mutex
	freeConn    *PriorityQueue
	initialCap  int
	maxCap      int
	maxIdle     int
	idleTime    time.Duration
	maxLifetime time.Duration
	cleanerCh   chan struct{}

	factory func() (net.Conn, error)
}

func NewHeapPool(config *PoolConfig) (Pool, error) {
	if config.InitialCap > config.MaxCap || config.Factory == nil {
		return nil, ErrConfigInvalid
	}

	initialCap := 5
	if config.InitialCap > 0 {
		initialCap = config.InitialCap
	}
	maxCap := 50
	if config.MaxCap > 0 {
		maxCap = config.MaxCap
	}
	maxIdle := 5
	if config.MaxIdle > 0 {
		maxIdle = config.MaxIdle
	}
	idleTime := 2 * time.Minute
	if config.Idletime > 0 {
		idleTime = config.Idletime
	}
	maxLifetime := 15 * time.Minute
	if config.MaxLifetime > 0 {
		maxLifetime = config.MaxLifetime
	}

	hp := &heapPool{
		initialCap:  initialCap,
		maxCap:      maxCap,
		maxIdle:     maxIdle,
		idleTime:    idleTime,
		maxLifetime: maxLifetime,
		cleanerCh:   make(chan struct{}),
		factory:     config.Factory,
	}

	pq := make(PriorityQueue, 0, maxCap)
	heap.Init(&pq)
	hp.freeConn = &pq

	type res struct {
		conn net.Conn
		err  error
	}
	ch := make(chan res, initialCap)
	for i := 0; i < initialCap; i++ {
		go func() {
			conn, err := hp.factory()
			ch <- res{conn: conn, err: err}
		}()
	}

	go func() {
		for i := 0; i < initialCap; i++ {
			ret := <-ch
			if ret.err != nil {
				log.Printf("init connection for easy pool err:%v", ret.err)
				continue
			}
			hp.put(hp.wrapConn(ret.conn).(*PoolConn))
		}
	}()

	go hp.cleaner()

	return hp, nil
}

func (hp *heapPool) Get() (net.Conn, error) {
	hp.mu.Lock()
	if hp.freeConn == nil {
		hp.mu.Unlock()
		return nil, ErrClosed
	}
	for hp.freeConn.Len() > 0 {
		pc := heap.Pop(hp.freeConn).(*PoolConn)
		if time.Now().Sub(pc.updatedTime) <= hp.maxLifetime {
			hp.mu.Unlock()
			return pc, nil
		}
		go pc.close()
	}
	hp.mu.Unlock()

	conn, err := hp.factory()
	if err != nil {
		return nil, err
	}
	return hp.wrapConn(conn), nil
}

func (hp *heapPool) Close() {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	if hp.freeConn == nil {
		return
	}
	hp.cleanerCh <- struct{}{}
	hp.factory = nil
	for hp.freeConn.Len() > 0 {
		pc := heap.Pop(hp.freeConn).(*PoolConn)
		pc.hp = nil
		pc.close()
	}
	hp.freeConn = nil
}

func (hp *heapPool) put(conn *PoolConn) error {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	if hp.freeConn == nil {
		return ErrClosed
	}
	if hp.freeConn.Len() >= hp.maxCap {
		return errors.New("pool have been filled")
	}
	heap.Push(hp.freeConn, conn)
	return nil
}

func (hp *heapPool) Len() int {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	if hp.freeConn == nil {
		return 0
	}
	return hp.freeConn.Len()
}

func (hp *heapPool) cleaner() {
	ticker := time.NewTicker(hp.idleTime / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			hp.mu.Lock()

			for hp.freeConn.Len() > 0 {
				pc := (*hp.freeConn)[0]
				interval := time.Now().Sub(pc.updatedTime)
				if interval >= hp.maxLifetime {
					_p := heap.Pop(hp.freeConn).(*PoolConn)
					go _p.close()
					continue
				}
				if interval >= hp.idleTime && hp.freeConn.Len() > hp.maxIdle {
					_p := heap.Pop(hp.freeConn).(*PoolConn)
					go _p.close()
					continue
				}
				break
			}

			hp.mu.Unlock()

		case <-hp.cleanerCh:
			log.Println("cleaner exited...")
			return
		}
	}
}

func (hp *heapPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn{hp: hp, updatedTime: time.Now()}
	p.Conn = conn
	return p
}
