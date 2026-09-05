# Steam Backup

## Motivation
Have you seen storage prices lately? My motivation for "Steam Backup" was twofold: 
   1. Create a way to transfer raw game files to and from my gaming PC to my NAS.
   2. Develop a project in GO, which I'm currently learning via Boot.dev

This could certainly have been accomplished with some scripting and robocopy, but I was looking for a practical application to help develop my GO skills.


## Prerequisites
- [Install Go](https://golang.org/dl/) (version 1.21+ recommended)

## Build & Run
1. **Initialize module** (if not already done):
   ```bash
   go mod init your-module-name
   ```

2. **Build the application**:
   ```bash
   go build -o steambackup
   ```

3. **Run the application**:
   ```bash
   ./steambackup
   ```

## Dependencies
```bash
go get -v
```

## Environment Variables

Create a `.env` file in the project root with the following configuration:

### Source Configuration
- `SRC_MAC_ADDR`: MAC address of the source machine (e.g., `11:22:33:44:55:66`)
- `SRC_IP_ADDR`: IP address of the source machine (e.g., `192.168.1.123`)
- `SRC_SHARENAME`: Name of the shared folder on the source (e.g., `steam` this should be a share of your current steam game files location)
- `SRC_BACKUP_DIR`: Relative path for backup storage on source (e.g., `.` for current directory)

### Destination Configuration
- `DEST_IP_ADDR`: IP address of the destination machine (e.g., `192.168.1.321`)
- `DEST_SHARENAME`: Name of the shared folder on the destination (e.g., `sharename`)
- `DEST_BACKUP_DIR`: Relative path for backup storage on destination (e.g., `./Games/SteamBackup`)

### Share Credentials
- `SHARE_USERNAME`: Username for network share access (e.g., `shareusername`)
- `SHARE_PASSWORD`: Password for network share access (e.g., `sharepassword`)


## Available Commands

| Command | Description |
|--------|-------------|
| `exit` | Exits the Steam Backup Utility |
| `help` | Displays this list of available commands |
| `wol`  | Sends a wake-on-LAN command to the source device |
| `lsd`  | Lists the contents of the source directory |
| `lbd`  | Lists the contents of the backup directory |
| `fb`   | Performs a full backup of the source directory to the backup location |
| `bbn`  | Performs a backup by name (specify a custom backup name) |
| `rbn`  | Restores a backup by name from the backup directory to the source |
