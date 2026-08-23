package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/RyanTarnowski/steambackup/internal/env"
	"github.com/RyanTarnowski/steambackup/internal/network"
)

const envPath = "./.env"
const smbPort = ":445"

func main() {
	err := env.LoadEnvironmentFile(envPath)
	if err != nil {
		log.Fatalf("Error reading .env file %v", err)
	}

	srcMAC := os.Getenv("SRC_MAC_ADDR")
	srcIP := os.Getenv("SRC_IP_ADDR")
	srcShareName := os.Getenv("SRC_SHARENAME")
	srcBackupDir := os.Getenv("SRC_BACKUP_DIR")

  destIP := os.Getenv("DEST_IP_ADDR")
	destShareName := os.Getenv("DEST_SHARENAME")
	destBackupDir := os.Getenv("DEST_BACKUP_DIR")
	
	shareUsername := os.Getenv("SHARE_USERNAME")
	sharePassword := os.Getenv("SHARE_PASSWORD")


	// fmt.Println("MAC Address:", src_mac)
	// fmt.Println("IP Address:", src_ip)
	
	fmt.Println("Sending WOL magic packet...")
	err = network.WakeOnLan(srcMAC)
	if err != nil {
		log.Fatalf("Error sending magic packet: %v", err)
	} else {
		fmt.Println("Magic packet sent successfully!")
	}

	fmt.Println("Check port status...")
	err = network.CheckPortStatus(srcIP, "445", 5, 2*time.Second, 1 * time.Second)
	if err != nil {
		log.Fatalf("CheckPortStatus: %v", err)
	} else {
		fmt.Println("Server is responding on port 445")
	}





	fmt.Printf("\nSource share files:\n")
	fmt.Println("*********************************************************")
	err = network.PrintSMBDirectory(srcIP + smbPort, srcShareName, shareUsername, sharePassword, srcBackupDir)
	if err != nil {
		log.Fatalf("Failed to connect to source: %v", err)
	}

	fmt.Printf("\nDestination share files:\n")
	fmt.Println("*********************************************************")
	err = network.PrintSMBDirectory(destIP + smbPort, destShareName, shareUsername, sharePassword, destBackupDir)
	if err != nil {
		log.Fatalf("Failed to connect to destination: %v", err)
	}






	//Connect to source
	srcFS, srcCleanup, err := network.ConnectSMB(srcIP + smbPort, shareUsername, sharePassword, srcShareName)
	if err != nil {
		log.Fatalf("Failed to connect to source: %v", err)
	}
	defer srcCleanup()

	//Connect to destination
	dstFS, dstCleanup, err := network.ConnectSMB(destIP + smbPort, shareUsername, sharePassword, destShareName)
	if err != nil {
		log.Fatalf("Failed to connect to destination: %v", err)
	}
	defer dstCleanup()

	//recursive copy operation
	srcFolder := "Arc Raiders"
	dstFolder := filepath.ToSlash(filepath.Join(destBackupDir, srcFolder))
	dstFolder = filepath.ToSlash(filepath.Join(destBackupDir, "BackupTest"))

	fmt.Printf("Starting transfer from Server A (%s) to Server B (%s)...\n", srcFolder, dstFolder)
	err = network.CopyFolderRemote(srcFS, dstFS, srcFolder, dstFolder)
	if err != nil {
		log.Fatalf("Copy failed: %v", err)
	}
	fmt.Println("Transfer completed successfully!")





	//TODO:
	//env vars - Done
	//WOL - Done
	//CheckPortStatus - Done
	//SMB connection - Done
	//File transfer - Done
	//Menu options
	//- WOL Source
	//- List Source
	//- List Destination
	//- Full Backup
	//- Backup by name
	//- Restore by name
	//Run on a schedule

}
