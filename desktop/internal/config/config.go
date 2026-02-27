package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the desktop app configuration
type Config struct {
	// APIBaseURL is the base URL for the pawnshop API
	APIBaseURL string `json:"api_base_url"`

	// ThermalPrinter is the name of the thermal printer for receipts
	ThermalPrinter string `json:"thermal_printer"`

	// DefaultPrinter is the name of the default document printer
	DefaultPrinter string `json:"default_printer"`

	// WindowWidth is the saved window width
	WindowWidth int `json:"window_width"`

	// WindowHeight is the saved window height
	WindowHeight int `json:"window_height"`

	// path is the config file path (not serialized)
	path string `json:"-"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		APIBaseURL:   "http://localhost:8090/api/v1",
		WindowWidth:  1280,
		WindowHeight: 800,
	}
}

// LoadConfig loads the configuration from disk or returns defaults
func LoadConfig() *Config {
	configPath := getConfigPath()
	config := DefaultConfig()
	config.path = configPath

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file doesn't exist, create with defaults
		config.Save()
		return config
	}

	if err := json.Unmarshal(data, config); err != nil {
		// Invalid config file, return defaults
		return config
	}

	return config
}

// Save persists the configuration to disk
func (c *Config) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}

// GetPath returns the config file path
func (c *Config) GetPath() string {
	return c.path
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to executable directory
		exe, _ := os.Executable()
		return filepath.Join(filepath.Dir(exe), "config.json")
	}
	return filepath.Join(configDir, "PawnshopPOS", "config.json")
}
