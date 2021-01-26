package controller

// detach will try to restart processes on failures.
func detach(c *Controller) {
	for proc := range c.Terminated {
		if err := c.Start(proc); err != nil {
			c.log.Err(err)
		}
	}
}
