package easypool

import (
	"errors"
	"net"
	"time"
)

var (
	ErrClosed        = errors.New("pool is closed")
	ErrConfigInvalid = errors.New("config is invalid")
)

type Pool interface {
	// Get return a new item from the pool. Closing the item puts it back to the pool
	Get() (net.Conn, error)
	// Close close the pool and release all resources
	Close()
	// Len returns the number of items of the pool
	Len() int
}

type PoolConfig struct {
	InitialCap  int
	MaxCap      int
	MaxIdle     int
	Idletime    time.Duration
	MaxLifetime time.Duration
	Factory     func() (net.Conn, error)
}
