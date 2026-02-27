package printer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// PrinterInfo contains information about a printer
type PrinterInfo struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	Status    string `json:"status"`
}

// PrinterService handles native printing operations
type PrinterService struct{}

// NewPrinterService creates a new printer service
func NewPrinterService() *PrinterService {
	return &PrinterService{}
}

// GetPrinters returns a list of available printers
func (p *PrinterService) GetPrinters() ([]PrinterInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return p.getWindowsPrinters()
	case "darwin":
		return p.getMacPrinters()
	case "linux":
		return p.getLinuxPrinters()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// GetDefaultPrinter returns the name of the default system printer
func (p *PrinterService) GetDefaultPrinter() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return p.getWindowsDefaultPrinter()
	case "darwin":
		return p.getMacDefaultPrinter()
	case "linux":
		return p.getLinuxDefaultPrinter()
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// PrintPDF prints a PDF file to the specified printer
func (p *PrinterService) PrintPDF(pdfBytes []byte, printerName string) error {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "pawnshop-print-*.pdf")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pdfBytes); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	switch runtime.GOOS {
	case "windows":
		return p.printWindowsPDF(tmpFile.Name(), printerName)
	case "darwin":
		return p.printMacPDF(tmpFile.Name(), printerName)
	case "linux":
		return p.printLinuxPDF(tmpFile.Name(), printerName)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// Windows-specific implementations

func (p *PrinterService) getWindowsPrinters() ([]PrinterInfo, error) {
	// Use PowerShell to get printers
	cmd := exec.Command("powershell", "-Command",
		"Get-Printer | Select-Object Name, Default, PrinterStatus | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get printers: %w", err)
	}

	// Parse output manually since it might be a single object or array
	outputStr := strings.TrimSpace(string(output))
	var printers []PrinterInfo

	if outputStr == "" {
		return printers, nil
	}

	// Simple parsing - PowerShell output format
	lines := strings.Split(outputStr, "\n")
	var currentPrinter PrinterInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"Name\"") {
			// Extract name
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				name = strings.Trim(name, "\",")
				currentPrinter.Name = name
			}
		} else if strings.HasPrefix(line, "\"Default\"") {
			if strings.Contains(line, "true") {
				currentPrinter.IsDefault = true
			}
		} else if strings.HasPrefix(line, "\"PrinterStatus\"") {
			// Extract status
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				status := strings.TrimSpace(parts[1])
				status = strings.Trim(status, "\",")
				currentPrinter.Status = status
			}
		} else if line == "}" || line == "}," {
			if currentPrinter.Name != "" {
				printers = append(printers, currentPrinter)
				currentPrinter = PrinterInfo{}
			}
		}
	}

	// Handle last printer if not added
	if currentPrinter.Name != "" {
		printers = append(printers, currentPrinter)
	}

	return printers, nil
}

func (p *PrinterService) getWindowsDefaultPrinter() (string, error) {
	cmd := exec.Command("powershell", "-Command",
		"(Get-Printer | Where-Object {$_.Default -eq $true}).Name")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default printer: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (p *PrinterService) printWindowsPDF(filePath, printerName string) error {
	// Try using SumatraPDF if available (recommended for silent printing)
	sumatraPath := findSumatraPDF()
	if sumatraPath != "" {
		cmd := exec.Command(sumatraPath,
			"-print-to", printerName,
			"-silent",
			filePath)
		return cmd.Run()
	}

	// Fallback to Adobe Reader if available
	adobePath := findAdobeReader()
	if adobePath != "" {
		cmd := exec.Command(adobePath,
			"/t", filePath, printerName)
		return cmd.Run()
	}

	// Fallback to Windows print command (opens dialog)
	cmd := exec.Command("rundll32", "mshtml.dll,PrintHTML", filePath)
	return cmd.Run()
}

// macOS-specific implementations

func (p *PrinterService) getMacPrinters() ([]PrinterInfo, error) {
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get printers: %w", err)
	}

	var printers []PrinterInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "printer ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				printers = append(printers, PrinterInfo{
					Name:   parts[1],
					Status: "Ready",
				})
			}
		}
	}

	// Mark default
	defaultPrinter, _ := p.getMacDefaultPrinter()
	for i := range printers {
		if printers[i].Name == defaultPrinter {
			printers[i].IsDefault = true
		}
	}

	return printers, nil
}

func (p *PrinterService) getMacDefaultPrinter() (string, error) {
	cmd := exec.Command("lpstat", "-d")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default printer: %w", err)
	}
	// Output format: "system default destination: PrinterName"
	parts := strings.Split(string(output), ":")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[1]), nil
	}
	return "", nil
}

func (p *PrinterService) printMacPDF(filePath, printerName string) error {
	cmd := exec.Command("lpr", "-P", printerName, filePath)
	return cmd.Run()
}

// Linux-specific implementations

func (p *PrinterService) getLinuxPrinters() ([]PrinterInfo, error) {
	// Same as macOS - uses CUPS
	return p.getMacPrinters()
}

func (p *PrinterService) getLinuxDefaultPrinter() (string, error) {
	return p.getMacDefaultPrinter()
}

func (p *PrinterService) printLinuxPDF(filePath, printerName string) error {
	return p.printMacPDF(filePath, printerName)
}

// Helper functions

func findSumatraPDF() string {
	// Common installation paths for SumatraPDF
	paths := []string{
		filepath.Join(os.Getenv("ProgramFiles"), "SumatraPDF", "SumatraPDF.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "SumatraPDF", "SumatraPDF.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "SumatraPDF", "SumatraPDF.exe"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func findAdobeReader() string {
	// Common installation paths for Adobe Reader
	paths := []string{
		filepath.Join(os.Getenv("ProgramFiles"), "Adobe", "Acrobat Reader DC", "Reader", "AcroRd32.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Adobe", "Acrobat Reader DC", "Reader", "AcroRd32.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Adobe", "Reader 11.0", "Reader", "AcroRd32.exe"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
