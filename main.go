package main

import (
	"log"
	"os"
	"github.com/RyanTarnowski/steambackup/internal/env"
)

const envPath = "./.env"

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

	//Check for command line args (index 0 is the program name)
	if len(os.Args) > 1 {
		startViaCmdLineArgs(os.Args[1])
		return
	}

	startRepl()
}
