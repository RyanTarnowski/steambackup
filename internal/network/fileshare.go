package network

import (
	"fmt"
	"io"
	"net"
	"os"
	"errors"
	"path/filepath"

	"github.com/hirochachacha/go-smb2"
)

//ConnectToSMB connects to an SMB share
func PrintSMBDirectory(server, share, user, password, backupDir string) error {
	//TCP connection to the SMB server (Port 445)
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return err
	}
	defer conn.Close()

	//setup auth 
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: password,
		},
	}

	//create session
	session, err := d.Dial(conn)
	if err != nil {
		return err
	}
	defer session.Logoff()

	//mount the target network share folder
	fs, err := session.Mount(share)
	if err != nil {
		return err
	}
	defer fs.Umount()

	//read backupDir into entries
	entries, err := fs.ReadDir(backupDir)
	if err != nil {
		return err
	}

	//print out entries
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	return nil
}


//ConnectSMB sets up the raw TCP connection, SMB session, and mounts the share
func ConnectSMB(addr, user, pass, share string) (*smb2.RemoteFileSystem, func(), error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: pass,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	fs, err := s.Mount(share)
	if err != nil {
		s.Logoff()
		conn.Close()
		return nil, nil, err
	}

	//clear up the connection (defer this after the call to connectSMB)
	cleanup := func() {
		fs.Umount()
		s.Logoff()
		conn.Close()
	}

	return fs, cleanup, nil
}

//CopyFolderRemote recursively reads the source share and recreates items on the destination share
func CopyFolderRemote(srcFS *smb2.RemoteFileSystem, dstFS *smb2.RemoteFileSystem, srcDir, dstDir string) error {
	//ensure the root destination directory exists
	err := dstFS.MkdirAll(dstDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination dir %s: %w", dstDir, err)
	}

	//read all files and folders inside the source directory
	entries, err := srcFS.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source dir %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.ToSlash(filepath.Join(srcDir, entry.Name()))
		dstPath := filepath.ToSlash(filepath.Join(dstDir, entry.Name()))

		if entry.IsDir() {
			//recurse subdirectories
			err = CopyFolderRemote(srcFS, dstFS, srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			//copy files to destination
			fmt.Printf("Copying: %s -> %s\n", srcPath, dstPath)
			err = copyFileRemote(srcFS, dstFS, srcPath, dstPath)
			if err != nil {
				return fmt.Errorf("failed copying file %s: %w", srcPath, err)
			}
		}
	}
	return nil
}

//copyFileRemote streams file content directly from the source pointer to the destination pointer
func copyFileRemote(srcFS *smb2.RemoteFileSystem, dstFS *smb2.RemoteFileSystem, srcPath, dstPath string) error {
	srcFile, err := srcFS.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := dstFS.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil //skip over the file if it exists
		}
		return err
	}
	defer dstFile.Close()

	//io.Copy uses an internal 32KB buffer to stream data smoothly without blowing up system RAM
	_, err = io.Copy(dstFile, srcFile)
	return err
}
