[![Build Status](https://travis-ci.org/nange/easypool.svg?branch=master)](https://travis-ci.org/nange/easypool)

# easypool
tcp连接池，最初是为[easyss](http://github.com/nange/easyss)设计，作为其tcp连接池管理使用。
也可用于其他需要tcp连接池的场景。


## 特性
* TCP连接池
* 支持基于时间的优先队列，越早创建和使用的连接，越早从Pool中取出
* 支持连接生命周期(最大空闲连接，连接最大存活时间等)

## 基本用法
```go
factory := func() (net.Conn, error) { return net.Dial("tcp", "localhost:7777") }
config := &PoolConfig{
	InitialCap:  5,
	MaxCap:      20,
	MaxIdle:     5,
	Idletime:    10 * time.Second,
	MaxLifetime: 10 * time.Minute,
	Factory:     factory,
}

pool, err := NewHeapPool(config)
if err != nil {
	log.Printf("err:%v\n", err)
	return
}

conn, err := pool.Get()
if err != nil {
	log.Printf("err:%v\n", err)
	return
}

// do sth... with conn
// 使用完连接之后，调用Close发放，当前连接会重新放回到pool中
conn.Close()

// 如果需要直接关闭底层连接(比如底层连接已经失效)，则：
conn.(*PoolConn).MarkUnusable()
conn.Close()

// 调用MarkUnusable()后，会返回true
conn.(*PoolConn).IsUnusable()

// 释放当前连接池中所有连接
pool.Close()

// 查看连接池中连接的个数
pool.Len()

```

## License

The MIT License (MIT) 
