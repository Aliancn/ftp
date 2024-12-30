package main

import (
	"context"
	"fmt"
)

type FTPClient struct {
	*FTPConn
	DEFAULT_CONTROL_PORT int
	DEFAULT_ADDR         string
	ctx                  context.Context
	cancel               context.CancelFunc
}

// NewFTPClient initializes a new FTP client
func NewFTPClient() *FTPClient {
	return &FTPClient{
		FTPConn:              NewFTPConn(),
		DEFAULT_CONTROL_PORT: 2121,
		DEFAULT_ADDR:         "127.0.0.1",
	}
}

// Connect to FTP server
func (a *App) Connect(address, username, password string) error {
	err := a.ftp.Dial(address)
	if err != nil {
		MyLogger.Info("failed to connect", err)
		return fmt.Errorf("failed to connect: %v", err)
	}
	MyLogger.Info("password ", password)
	if err := a.ftp.Login(username, password); err != nil {
		MyLogger.Info("failed to login: ", err)
		return fmt.Errorf("failed to login: %v", err)
	}

	return nil
}

// List files and directories
func (a *App) List(path string) ([]string, error) {
	if a.ftp.controlConn == nil {
		MyLogger.Info("not connected")
		return nil, fmt.Errorf("not connected")
	}

	entries, err := a.ftp.ListFiles(path)
	if err != nil {
		MyLogger.Info("failed to list directory: ", err)
		return nil, fmt.Errorf("failed to list directory: %v", err)
	}

	return entries, nil
}

// Upload file
func (a *App) Upload(localFile, remotePath string) error {
	if a.ftp.controlConn == nil {
		MyLogger.Info("not connected")
		return fmt.Errorf("not connected")
	}

	if err := a.ftp.STOR(localFile, remotePath); err != nil {
		MyLogger.Info("failed to upload file: %v", err)
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

// Download file
func (a *App) Download(remotePath, localPath string, size int64) error {
	if a.ftp.controlConn == nil {
		MyLogger.Info("not connected")
		return fmt.Errorf("not connected")
	}

	if a.ftp.ctx == nil {
		a.ftp.ctx, a.ftp.cancel = context.WithCancel(context.Background())
	} else {
		a.ftp.ctx, a.ftp.cancel = context.WithCancel(context.Background())
	}

	// 获取本地文件大小
	localFileSize, err := GetDownloadedOffset(localPath)
	if err != nil {
		MyLogger.Info("获取本地文件大小失败: ", err)
		return err
	}

	// 恢复下载
	err = a.ftp.REST_RETR(remotePath, localPath, localFileSize, a.ctx)
	if err != nil {
		MyLogger.Info("恢复下载失败: ", err)
		return err
	}

	return nil
}

func (a *App) StopDownload() error {
	if a.ftp.controlConn == nil {
		MyLogger.Info("not connected")
		return fmt.Errorf("not connected")
	}
	if a.ftp.cancel != nil {
		fmt.Println("cancel")
		a.ftp.cancel()
		a.ftp.cancel = nil
	}

	return nil
}

// Create folder
func (a *App) CreateFolder(path string) error {
	if a.ftp.controlConn == nil {
		return fmt.Errorf("not connected")
	}

	if err := a.ftp.MakeDir(path); err != nil {
		MyLogger.Info("failed to create folder: ", err)
		return fmt.Errorf("failed to create folder: %v", err)
	}
	return nil
}

// Delete folder or file
func (a *App) Delete(path string) error {
	if a.ftp.controlConn == nil {
		return fmt.Errorf("not connected")
	}

	if err := a.ftp.Dele(path, false); err != nil {
		MyLogger.Info("failed to delete: ", err)
		return fmt.Errorf("failed to delete: %v", err)
	}
	return nil
}

// Disconnect from FTP server
func (a *App) Disconnect() error {
	if a.ftp.controlConn != nil {
		return a.ftp.Close()
	}
	return nil
}
