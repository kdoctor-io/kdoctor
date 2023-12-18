package connexus

import (
	"container/heap"
	"errors"
	"net"
	"sort"
	"sync"
	"time"
)

var (
	ErrClosed = errors.New("pool is closed")
)

const MaxConnectCount = 20000
const MaxIdleConnectCount = 25000

type Pool interface {
	// Get return a new item from the pool. Closing the item puts it back to the pool
	Get() (net.Conn, error)
	// Close close the pool and release all resources
	Close()
	// Len returns the number of items of the pool
	Len() int
}

type PoolConfig struct {
	Cap        int
	MaxIdleCap int
	Factory    func() (net.Conn, error)
}

type connexPool struct {
	mu         sync.Mutex
	freeConn   *PriorityQueue
	cap        int
	maxIdleCap int
	cleanerCh  chan struct{}
	count      int
	factory    func() (net.Conn, error)
}

func NewConnexPool(cfg PoolConfig) (Pool, error) {
	if cfg.Cap > cfg.MaxIdleCap {
		return nil, errors.New("Cap can not more than MaxIdleCap")
	}

	var poolCap int
	var idleConnCap int

	if cfg.Cap > MaxConnectCount {
		poolCap = MaxConnectCount
	} else {
		poolCap = cfg.Cap
	}

	if cfg.MaxIdleCap > MaxIdleConnectCount {
		idleConnCap = MaxIdleConnectCount
	} else {
		idleConnCap = cfg.Cap
	}

	cp := &connexPool{
		cap:        poolCap,
		cleanerCh:  make(chan struct{}, 1),
		factory:    cfg.Factory,
		maxIdleCap: idleConnCap,
	}

	pq := make(PriorityQueue, 0, cfg.Cap)
	heap.Init(&pq)
	cp.freeConn = &pq

	for i := 0; i < poolCap; i++ {
		conn, err := cfg.Factory()
		if err != nil {
			continue
		}
		cp.put(cp.wrapConn(conn).(*Connex))
		cp.count++
	}

	go cp.inducer()

	return cp, nil
}

func (cp *connexPool) Get() (net.Conn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.freeConn == nil {
		return nil, ErrClosed
	}

	if cp.freeConn.Len() > 0 {
		return heap.Pop(cp.freeConn).(*Connex), nil
	}

	if cp.count < cp.maxIdleCap {
		conn, err := cp.factory()
		cp.count++
		if err != nil {
			return nil, err
		}
		return cp.wrapConn(conn), nil
	}

	return nil, errors.New("can not allocate connect")
}

func (cp *connexPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.freeConn == nil {
		return
	}
	cp.cleanerCh <- struct{}{}
	cp.factory = nil
	for cp.freeConn.Len() > 0 {
		c := heap.Pop(cp.freeConn).(*Connex)
		c.cp = nil
		c.Conn.Close()
	}
	cp.freeConn = nil
}

func (cp *connexPool) put(conn *Connex) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.freeConn == nil {
		return ErrClosed
	}
	if cp.freeConn.Len() >= cp.cap {
		return errors.New("pool have been filled")
	}
	heap.Push(cp.freeConn, conn)
	return nil
}

func (cp *connexPool) Len() int {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.freeConn == nil {
		return 0
	}
	return cp.freeConn.Len()
}

func (cp *connexPool) inducer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cp.mu.Lock()
			if cp.freeConn != nil {
				sort.Sort(cp.freeConn)
			}
			cp.mu.Unlock()

		case <-cp.cleanerCh:
			return
		}
	}
}

func (cp *connexPool) wrapConn(conn net.Conn) net.Conn {
	p := &Connex{cp: cp, updatedTime: time.Now()}
	p.Conn = conn
	return p
}
