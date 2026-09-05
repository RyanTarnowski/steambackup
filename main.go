package main

import (
	"log"
	"os"
	"github.com/RyanTarnowski/steambackup/internal/env"
)

const envPath = "./.env"

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
