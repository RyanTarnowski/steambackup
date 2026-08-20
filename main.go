package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RyanTarnowski/steambackup/internal/env"
	"github.com/RyanTarnowski/steambackup/internal/network"
)

const envPath = "./.env"

func main() {
	err := env.LoadEnvironmentFile(envPath)
	if err != nil {
		log.Fatalf("Error reading .env file %v", err)
	}

	mac := os.Getenv("SOURCE_MAC_ADDR")
	ip := os.Getenv("SOURCE_IP_ADDR")

	fmt.Println("MAC Address:", mac)
	fmt.Println("IP Address:", ip)

	//return //Testing env reader

	fmt.Println("Sending WOL magic packet...")
	err = network.WakeOnLan(mac)
	if err != nil {
		fmt.Println("Error sending magic packet:", err)
	} else {
		fmt.Println("Magic packet sent successfully!")
	}

	fmt.Println("Check port status...")
	err = network.CheckPortStatus(ip, "445", 2*time.Second)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Server is responding on port 445")
	}

	//TODO:
	//env vars - Done
	//WOL - Done
	//CheckPortStatus - in progress, add attempts args
	//SMB connection
	//File transfer
	//Logging
	//Email
	//Run on a schedule
	//Restore file back to source

}
