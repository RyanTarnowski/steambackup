package main

import (
  "os"
	"fmt"
  "log"
	"time"
	"errors"
	"path/filepath"
	"github.com/RyanTarnowski/steambackup/internal/network"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	srcMAC             string
	srcIP              string
	srcShareName       string
	srcBackupDir       string
	destIP             string
	destShareName      string
	destBackupDir      string
	shareUsername      string
	sharePassword      string
}

const horizontalLine = "*********************************************************"

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit Steam Backup Utility",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays available commands",
			callback:    commandHelp,
		},
		"wol": {
			name:        "wol",
			description: "Send wake on lan command to source",
			callback:    commandWOL,
		},
		"lsd": {
		 	name:        "lsd",
		 	description: "Lists source directory (lsd)",
		 	callback:    commandListSourceDirectory,
		},
		"lbd": {
		 	name:        "lbd",
		 	description: "Lists backup directory (lbd)",
		 	callback:    commandListDestinationDirectory,
		},
		"fb": {
		 	name:        "fb",
		 	description: "Performs a full backup (fb) of the source directory to the backup directory",
		 	callback:    commandFullBackup,
		},
		"bbn": {
		 	name:        "bbn",
		 	description: "Performs a backup by name (bbn) of the source directory to the backup directory",
		 	callback:    commandBackupByName,
		},
		"rbn": {
		 	name:        "rbn",
		 	description: "Performs a retore by name (rbn) of the backup directory to the source directory",
		 	callback:    commandRestoreByName,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Print("Closing the Steam Backup Utility...")
	os.Exit(0)

	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Usage:")
	fmt.Println(horizontalLine)
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println(horizontalLine)
	return nil
}

func commandWOL(cfg *config, args ...string) error {
	fmt.Println("Sending WOL magic packet...")
	fmt.Println(horizontalLine)
	err := network.WakeOnLan(cfg.srcMAC)
	if err != nil {
		log.Fatalf("Error sending magic packet: %v", err)
		return fmt.Errorf("Error sending magic packet: %v", err)
	} else {
		fmt.Println("Magic packet sent successfully!")
	}

	fmt.Println("Check port status...")
	err = network.CheckPortStatus(cfg.srcIP, "445", 10, 2*time.Second, 1 * time.Second)
	if err != nil {
		log.Fatalf("CheckPortStatus: %v", err)
		return fmt.Errorf("CheckPortStatus: %v", err)
	} else {
		fmt.Println("Server is responding on port 445")
	}
	
	fmt.Println(horizontalLine)
	return nil
}

func commandListSourceDirectory(cfg *config, args ...string) error {
	fmt.Printf("\nSource Directory:\n")
	fmt.Println(horizontalLine)
	err := network.PrintSMBDirectory(cfg.srcIP + smbPort, cfg.srcShareName, cfg.shareUsername, cfg.sharePassword, cfg.srcBackupDir)

	if err != nil {
		log.Fatalf("Failed to connect to source: %v", err)
		return fmt.Errorf("Failed to connect to source: %v", err)
	}

	fmt.Println(horizontalLine)
	return nil
}

func commandListDestinationDirectory(cfg *config, args ...string) error {
	fmt.Printf("\nBackup Directory:\n")
	fmt.Println(horizontalLine)
	err := network.PrintSMBDirectory(cfg.destIP + smbPort, cfg.destShareName, cfg.shareUsername, cfg.sharePassword, cfg.destBackupDir)

	if err != nil {
		log.Fatalf("Failed to connect to backup: %v", err)
		return fmt.Errorf("Failed to connect to backup: %v", err)
	}

	fmt.Println(horizontalLine)
	return nil
}

func commandFullBackup(cfg *config, args ...string) error {
	fmt.Printf("\nFull Backup:\n")
	fmt.Println(horizontalLine)

	fmt.Println("This feature is not ready")

	fmt.Println(horizontalLine)
	return nil
}

func commandBackupByName(cfg *config, args ...string) error {
	fmt.Printf("\nBackup by Name:\n")
	fmt.Println(horizontalLine)

	if len(args) != 1 || args[0] == "" {
		return errors.New("Source directory name required. Use command 'lsd' to view source directories")
	}

	//Connect to source
	srcFS, srcCleanup, err := network.ConnectSMB(cfg.srcIP + smbPort, cfg.shareUsername, cfg.sharePassword, cfg.srcShareName)
	if err != nil {
		log.Fatalf("Failed to connect to source: %v", err)
		return fmt.Errorf("Failed to connect to source: %v", err)
	}
	defer srcCleanup()

	//Connect to destination
	dstFS, dstCleanup, err := network.ConnectSMB(cfg.destIP + smbPort, cfg.shareUsername, cfg.sharePassword, cfg.destShareName)
	if err != nil {
		log.Fatalf("Failed to connect to backup: %v", err)
		return fmt.Errorf("Failed to connect to backup: %v", err)
	}
	defer dstCleanup()

	srcFolder := args[0]
	dstFolder := filepath.ToSlash(filepath.Join(cfg.destBackupDir, srcFolder))

	fmt.Printf("Starting backup of (%s) to (%s)...\n", srcFolder, dstFolder)
	
	//recursive copy operation
	err = network.CopyFolderRemote(srcFS, dstFS, srcFolder, dstFolder)
	if err != nil {
		log.Fatalf("Backup failed: %v", err)
		return fmt.Errorf("Backup failed: %v", err)
	}

	fmt.Printf("\nBackup completed successfully!\n")
	fmt.Println(horizontalLine)
	return nil
}

func commandRestoreByName(cfg *config, args ...string) error {
	fmt.Printf("\nRestore by Name:\n")
	fmt.Println(horizontalLine)

	if len(args) != 1 || args[0] == "" {
		return errors.New("Source directory name required. Use command 'lbd' to view backup directories")
	}

	//Connect to source
	srcFS, srcCleanup, err := network.ConnectSMB(cfg.srcIP + smbPort, cfg.shareUsername, cfg.sharePassword, cfg.srcShareName)
	if err != nil {
		log.Fatalf("Failed to connect to source: %v", err)
		return fmt.Errorf("Failed to connect to source: %v", err)
	}
	defer srcCleanup()

	//Connect to destination
	dstFS, dstCleanup, err := network.ConnectSMB(cfg.destIP + smbPort, cfg.shareUsername, cfg.sharePassword, cfg.destShareName)
	if err != nil {
		log.Fatalf("Failed to connect to backup: %v", err)
		return fmt.Errorf("Failed to connect to backup: %v", err)
	}
	defer dstCleanup()

	dstFolder := filepath.ToSlash(filepath.Join(cfg.destBackupDir, args[0]))
	// srcFolder := filepath.ToSlash(filepath.Join(cfg.srcBackupDir, "BackupTest", args[0]))
	srcFolder := filepath.ToSlash(filepath.Join(cfg.srcBackupDir, args[0]))

	fmt.Printf("Starting backup of (%s) to (%s)...\n", dstFolder, srcFolder)
	
	//recursive copy operation
	err = network.CopyFolderRemote(dstFS, srcFS, dstFolder, srcFolder)
	if err != nil {
		log.Fatalf("Restore failed: %v", err)
		return fmt.Errorf("Restore failed: %v", err)
	}

	fmt.Printf("\nRestore completed successfully!\n")
	fmt.Println(horizontalLine)
	return nil
}
