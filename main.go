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
	//- Restore by name
	//Run on a schedule
	//Rename all refs to destination to backup (ldd -> lbd)

//TEST:
//backup to dest when folder does not exist on dest
//backup to dest when folder does exist on dest
//restore to src when folder does not exist on src
//restore to src when folder does exist on src
//full backup to dest for src

func main() {
	err := env.LoadEnvironmentFile(envPath)
	if err != nil {
		log.Fatalf("Error reading .env file %v", err)
	}

	startRepl()
}
