package main

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var filePath string
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			// Resolve absolute path if possible
			if absPath, err := filepath.Abs(arg); err == nil {
				filePath = absPath
			} else {
				filePath = arg
			}
			break
		}
	}

	// Create an instance of the app structure
	app := NewApp(filePath)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "mdviewer",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 30, G: 30, B: 30, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
