package main

import (
	"log"
	"github.com/RyanTarnowski/steambackup/internal/env"
)

const envPath = "./.env"
const smbPort = ":445"


//TODO:
//env vars - Done
//WOL - Done
//CheckPortStatus - Done
//SMB connection - Done
//File transfer - Done
//Menu options
//- WOL Source - Done
//- List Source - Done
//- List Destination - Done
//- Full Backup
//- Backup by name -Done
//- Restore by name -Done
//Run on a schedule - Look into Ofelia docker
//Rename all refs to destination to backup (ldd -> lbd) -done

//TEST:
//full backup to dest for src

func main() {
	err := env.LoadEnvironmentFile(envPath)
	if err != nil {
		log.Fatalf("Error reading .env file %v", err)
	}

	startRepl()
}
