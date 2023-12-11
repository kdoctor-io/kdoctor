package connexus

import (
	"net"
	"time"
)

type Connex struct {
	net.Conn
	cp          *connexPool
	updatedTime time.Time
	unusable    bool
}

// Close put the connection back to pool if possible.
// Executed by multi times is ok.
func (c *Connex) Close() error {
	if c.unusable {
		return c.close()
	}

	c.updatedTime = time.Now()

	if c.cp.Len() > c.cp.MaxIdleCap {
		c.cp = nil
		return c.Conn.Close()
	}

	if err := c.cp.put(c); err != nil {
		c.cp = nil
		return c.Conn.Close()
	}
	return nil
}

func (c *Connex) MarkUnusable() {
	c.unusable = true
}

func (c *Connex) IsUnusable() bool {
	return c.unusable
}

func (c *Connex) close() error {
	return c.Conn.Close()
}
