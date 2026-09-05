package network

import (
	"fmt"
	"net"
	"time"
)

func CheckPortStatus(host, port string, attempts int, timeout, delay time.Duration) error {
	address := net.JoinHostPort(host, port)

	//Run a loop with a delay that attempts to connect to the host via the port
	//return when connection is successful or when attempts run out
	for i := range attempts {
		fmt.Printf("Attempt #%v checking port...\n", i+1)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			conn.Close()
			return nil
		} else {
			fmt.Println("Error:", err)
		}
		time.Sleep(delay)
	}

	return fmt.Errorf("Failed to establish connection to %s:%s after %v attempts.", host, port, attempts)
}
