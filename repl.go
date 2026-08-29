package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{
		srcMAC: os.Getenv("SRC_MAC_ADDR"),
		srcIP: os.Getenv("SRC_IP_ADDR"),
		srcShareName: os.Getenv("SRC_SHARENAME"),
		srcBackupDir: os.Getenv("SRC_BACKUP_DIR"),
		destIP: os.Getenv("DEST_IP_ADDR"),
		destShareName: os.Getenv("DEST_SHARENAME"),
		destBackupDir: os.Getenv("DEST_BACKUP_DIR"),
		shareUsername: os.Getenv("SHARE_USERNAME"),
		sharePassword: os.Getenv("SHARE_PASSWORD"),
	}

	for {
		fmt.Println()
		fmt.Print("Steam Backup Utility > ")

		if scanner.Scan() {
			//Cut the cmd out on the input string and consider the remaining text as args
			cmd, args, _ := strings.Cut(scanner.Text(), " ")

			if command, ok := getCommands()[cmd]; ok {
					err := command.callback(&cfg, args)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("Unknown command")
				}
		}
	}
}
