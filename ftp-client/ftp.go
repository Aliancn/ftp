package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPClient struct {
	connection *ftp.ServerConn
}

// Connect to FTP server
func (f *FTPClient) Connect(address, username, password string) error {
	conn, err := ftp.Dial(address, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	if err := conn.Login(username, password); err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	f.connection = conn
	return nil
}

// List files and directories
func (f *FTPClient) List(path string) ([]*ftp.Entry, error) {
	if f.connection == nil {
		return nil, fmt.Errorf("not connected")
	}

	entries, err := f.connection.List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %v", err)
	}

	// var files []string
	// for _, entry := range entries {
	// 	MyLogger.Info(entry.Name)
	// 	files = append(files, entry.Name)
	// }
	// return files, nil

	return entries, nil
}

// Upload file
func (f *FTPClient) Upload(localFile, remotePath string) error {
	if f.connection == nil {
		MyLogger.Info("not connected")
		return fmt.Errorf("not connected")
	}

	file, err := os.Open(localFile)
	if err != nil {
		MyLogger.Info("failed to open local file: ", localFile)
		MyLogger.Info("failed to open local file: ", err)
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	if err := f.connection.Stor(remotePath, file); err != nil {
		MyLogger.Info("failed to upload file: %v", err)
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

// Download file
func (f *FTPClient) Download(remotePath, localPath string) error {
	if f.connection == nil {
		return fmt.Errorf("not connected")
	}

	resp, err := f.connection.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Close()

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	return nil
}

// Create folder
func (f *FTPClient) CreateFolder(path string) error {
	if f.connection == nil {
		return fmt.Errorf("not connected")
	}

	if err := f.connection.MakeDir(path); err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}
	return nil
}

// Delete folder or file
func (f *FTPClient) Delete(path string) error {
	if f.connection == nil {
		return fmt.Errorf("not connected")
	}

	if err := f.connection.Delete(path); err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}
	return nil
}

// Disconnect from FTP server
func (f *FTPClient) Disconnect() error {
	if f.connection != nil {
		return f.connection.Quit()
	}
	return nil
}
