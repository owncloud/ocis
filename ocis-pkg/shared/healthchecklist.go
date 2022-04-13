package shared

import "net"

// Check is a single health-check
type Check func() error

// RunChecklist runs all the given checks
func RunChecklist(checks ...Check) error {
	for _, c := range checks {
		err := c()
		if err != nil {
			return err
		}
	}
	return nil
}

// TCPConnect connects to a given tcp endpoint
func TCPConnect(host string) Check {
	return func() error {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			return err
		}
		defer conn.Close()
		return nil
	}
}
