package controller

import (
	"time"
)

// Dispatch routes client commands to their respective handlers
func (c *Controller) Dispatch() error {
	c.conn.SetDeadline(time.Now().Add(3e9))

	return nil
}
