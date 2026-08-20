package network

import (
	"fmt"
	"net"
	"time"
)

func CheckPortStatus(host, port string, timeout time.Duration) error {
	address := net.JoinHostPort(host, port)

	//Run a loop with a delay that attempts to connect to the host via the port
	//return when connection is successful
	//TODO: Only allow X attempts before returning an err
	for {
		fmt.Println("checking port...")
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			conn.Close()
			return nil
		} else {
			fmt.Println("Error:", err)
		}
		time.Sleep(1 * time.Second)
	}
}
