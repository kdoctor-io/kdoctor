package easypool

import (
	"log"
	"net"
	"time"
)

type PoolConn struct {
	net.Conn
	hp          *heapPool
	updatedTime time.Time
	unusable    bool
}

// Close put the connection back to pool if possible.
// Executed by multi times is ok.
func (pc *PoolConn) Close() error {
	if pc.unusable {
		return pc.close()
	}

	pc.updatedTime = time.Now()
	if err := pc.hp.put(pc); err != nil {
		log.Printf("put conn failed:%v\n", err)
		pc.hp = nil
		return pc.close()
	}
	return nil
}

func (pc *PoolConn) MarkUnusable() {
	pc.unusable = true
}

func (pc *PoolConn) IsUnusable() bool {
	return pc.unusable
}

func (pc *PoolConn) close() error {
	return pc.Conn.Close()
}
