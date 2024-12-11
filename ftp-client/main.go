package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

var MyLogger *MySLog

func init() {
	MyLogger = NewMySlog("info", "log.log")
}

func main() {
	// Create an instance of the app structure
	app := NewApp()
	ftpClient := &FTPClient{}
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ftp-client",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		AlwaysOnTop:      false,
		Bind: []interface{}{
			app,
			ftpClient,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "FTP Application",
				Message: "© 2024 Aliancn",
			},
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "wails",
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
