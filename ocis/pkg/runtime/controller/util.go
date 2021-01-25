package controller

// detach will try to restart processes on failures.
func detach(c *Controller) {
	for {
		select {
		case proc := <-c.Terminated:
			if err := c.Start(proc); err != nil {
				c.log.Err(err)
			}
		}
	}
}
