package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RyanTarnowski/steambackup/internal/wol"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access variables using the standard os package
	mac := os.Getenv("SOURCE_MAC_ADDR")
	fmt.Println("MAC Address:", mac)

	fmt.Println("Sending WOL magic packet")

	err = wol.WakeOnLan("")
	if err != nil {
		fmt.Println("Error sending magic packet:", err)
	} else {
		fmt.Println("Magic packet sent successfully!")
	}

	//TODO:
	//1: WOL
	//2: SMB connection
	//3: File transfer
	//4: Logging
	//5: Email
	//7: Run on a schedule
	//8: Restore file back to source

}
