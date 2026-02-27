package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"pawnshop-desktop/internal/config"
	"pawnshop-desktop/internal/printer"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	config  *config.Config
	printer *printer.PrinterService
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg := config.LoadConfig()
	return &App{
		config:  cfg,
		printer: printer.NewPrinterService(),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Cleanup resources if needed
}

// GetAPIBaseURL returns the configured API base URL
func (a *App) GetAPIBaseURL() string {
	return a.config.APIBaseURL
}

// SetAPIBaseURL updates the API base URL configuration
func (a *App) SetAPIBaseURL(url string) error {
	a.config.APIBaseURL = url
	return a.config.Save()
}

// GetAppVersion returns the application version
func (a *App) GetAppVersion() string {
	return "1.0.0"
}

// GetPrinters returns a list of available printers
func (a *App) GetPrinters() ([]printer.PrinterInfo, error) {
	return a.printer.GetPrinters()
}

// GetDefaultPrinter returns the default system printer name
func (a *App) GetDefaultPrinter() (string, error) {
	return a.printer.GetDefaultPrinter()
}

// GetThermalPrinter returns the configured thermal printer name
func (a *App) GetThermalPrinter() string {
	return a.config.ThermalPrinter
}

// SetThermalPrinter sets the thermal printer for receipts
func (a *App) SetThermalPrinter(printerName string) error {
	a.config.ThermalPrinter = printerName
	return a.config.Save()
}

// PrintDocument prints a PDF document to the default printer
// pdfBase64 is the base64-encoded PDF content
func (a *App) PrintDocument(pdfBase64 string) error {
	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	defaultPrinter, err := a.printer.GetDefaultPrinter()
	if err != nil {
		return fmt.Errorf("failed to get default printer: %w", err)
	}

	return a.printer.PrintPDF(pdfBytes, defaultPrinter)
}

// PrintThermalTicket prints a PDF to the configured thermal printer
// pdfBase64 is the base64-encoded PDF content
func (a *App) PrintThermalTicket(pdfBase64 string) error {
	if a.config.ThermalPrinter == "" {
		return fmt.Errorf("no thermal printer configured")
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	return a.printer.PrintPDF(pdfBytes, a.config.ThermalPrinter)
}

// PrintToSpecificPrinter prints a PDF to a specific printer
func (a *App) PrintToSpecificPrinter(pdfBase64 string, printerName string) error {
	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	return a.printer.PrintPDF(pdfBytes, printerName)
}

// ShowSaveDialog shows a native file save dialog and returns the selected path
func (a *App) ShowSaveDialog(defaultFilename string, filterDescription string, filterPattern string) (string, error) {
	return wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultFilename,
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: filterDescription,
				Pattern:     filterPattern,
			},
		},
	})
}

// SaveFile saves content to a file using native dialog
func (a *App) SaveFile(contentBase64 string, defaultFilename string) error {
	path, err := a.ShowSaveDialog(defaultFilename, "PDF Files (*.pdf)", "*.pdf")
	if err != nil {
		return err
	}
	if path == "" {
		return nil // User cancelled
	}

	content, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return fmt.Errorf("failed to decode content: %w", err)
	}

	return os.WriteFile(path, content, 0644)
}

// GetConfigPath returns the config file path
func (a *App) GetConfigPath() string {
	return a.config.GetPath()
}

// GetPlatform returns the current platform
func (a *App) GetPlatform() string {
	return runtime.GOOS
}

// OpenExternalURL opens a URL in the default browser
func (a *App) OpenExternalURL(url string) {
	wailsRuntime.BrowserOpenURL(a.ctx, url)
}

// GetDataDirectory returns the application data directory
func (a *App) GetDataDirectory() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "PawnshopPOS")
}

// IsDesktopApp returns true to indicate this is running as a desktop app
func (a *App) IsDesktopApp() bool {
	return true
}
