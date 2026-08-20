package network

import (
	"fmt"
	"net"
	"strings"
)

func WakeOnLan(macAddr string) error {
	//validate the macAddr
	var mac [6]byte
	parts := strings.Split(macAddr, ":")
	if len(parts) != 6 {
		return fmt.Errorf("invalid MAC address format")
	}

	//for each part of the split macAddr parse to a byte and add to the mac slice
	for i, part := range parts {
		var b byte
		fmt.Sscanf(part, "%2x", &b)
		mac[i] = b
	}

	//create the magic packet
	//magic packet consists of 6 bytes of 0xFF
	//followed by the target device's 6-byte MAC address repeated 16 times.
	var packet [102]byte //6 bytes of 0xFF + 96 (6 * 16) = 102
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:], mac[:])
	}

	//send the magic packet over udp
	/*
	* The User Datagram Protocol (UDP) is a fast,
	* lightweight network protocol that sends data directly to a target
	* without establishing a connection or verifying delivery.
	* It prioritizes speed and low latency over reliability,
	* making it ideal for time-sensitive applications.
	 */
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(packet[:])
	return err
}
