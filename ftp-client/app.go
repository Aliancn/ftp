package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	ftp *FTPClient
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		ftp: NewFTPClient(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) OpenAndUploadFile() (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a file to upload",
	})
	if err != nil || filePath == "" {
		return "", fmt.Errorf("no file selected or error occurred: %w", err)
	}
	return filePath, nil
}

