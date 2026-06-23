package main

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	initialPath string
}

// NewApp creates a new App application struct
func NewApp(initialPath string) *App {
	return &App{initialPath: initialPath}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// FileInfo holds the details of a markdown file
type FileInfo struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Name    string `json:"name"`
}

// GetInitialFile returns the file passed via CLI on startup
func (a *App) GetInitialFile() FileInfo {
	if a.initialPath == "" {
		return FileInfo{}
	}
	content, err := os.ReadFile(a.initialPath)
	if err != nil {
		return FileInfo{Path: a.initialPath, Name: filepath.Base(a.initialPath)}
	}
	return FileInfo{
		Path:    a.initialPath,
		Content: string(content),
		Name:    filepath.Base(a.initialPath),
	}
}

// ReadFile reads the content of a file
func (a *App) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SaveFile writes content to a file
func (a *App) SaveFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// OpenFileDialog opens a system file dialog to select a markdown file
func (a *App) OpenFileDialog() (FileInfo, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Markdown File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files (*.md, *.markdown)",
				Pattern:     "*.md;*.markdown",
			},
		},
	})
	if err != nil {
		return FileInfo{}, err
	}
	if path == "" {
		return FileInfo{}, nil // user cancelled
	}
	content, err := a.ReadFile(path)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Path:    path,
		Content: content,
		Name:    filepath.Base(path),
	}, nil
}

// SaveFileDialog opens a system file dialog to save a new markdown file
func (a *App) SaveFileDialog(content string) (FileInfo, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Markdown File",
		DefaultFilename: "document.md",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files (*.md, *.markdown)",
				Pattern:     "*.md;*.markdown",
			},
		},
	})
	if err != nil {
		return FileInfo{}, err
	}
	if path == "" {
		return FileInfo{}, nil // user cancelled
	}
	err = a.SaveFile(path, content)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Path:    path,
		Content: content,
		Name:    filepath.Base(path),
	}, nil
}

// ExportPdfDialog opens a save dialog for PDF and writes the base64 content
func (a *App) ExportPdfDialog(base64Content string, defaultName string) (FileInfo, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export PDF",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files (*.pdf)",
				Pattern:     "*.pdf",
			},
		},
	})
	if err != nil {
		return FileInfo{}, err
	}
	if path == "" {
		return FileInfo{}, nil // user cancelled
	}

	// Decode base64 content
	dec, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return FileInfo{}, err
	}

	// Write to PDF file
	err = os.WriteFile(path, dec, 0644)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Path: path,
		Name: filepath.Base(path),
	}, nil
}
