// Package config provides Viper-based configuration management for acorn CLI.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// AppName is the application name used for config directories
	AppName = "acorn"
)

// Config holds the application configuration
type Config struct {
	// Debug enables verbose output
	Debug bool `mapstructure:"debug"`
	// DotfilesRoot is the path to the dotfiles repository
	DotfilesRoot string `mapstructure:"dotfiles_root"`
	// Editor is the preferred text editor
	Editor string `mapstructure:"editor"`
	// Shell is the preferred shell
	Shell string `mapstructure:"shell"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Debug:        false,
		DotfilesRoot: filepath.Join(home, ".config", "dotfiles"),
		Editor:       getEnvOrDefault("EDITOR", "vim"),
		Shell:        getEnvOrDefault("SHELL", "/bin/bash"),
	}
}

// Load initializes Viper and loads the configuration
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	defaults := DefaultConfig()
	v.SetDefault("debug", defaults.Debug)
	v.SetDefault("dotfiles_root", defaults.DotfilesRoot)
	v.SetDefault("editor", defaults.Editor)
	v.SetDefault("shell", defaults.Shell)

	// Config file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Search paths (XDG compliant)
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	v.AddConfigPath(filepath.Join(configHome, AppName))
	v.AddConfigPath(".")

	// Environment variable binding
	v.SetEnvPrefix("ACORN")
	v.AutomaticEnv()

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		// Config file not found is OK - use defaults
	}

	// Unmarshal into struct
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return cfg, nil
}

// ConfigDir returns the XDG-compliant config directory for acorn
func ConfigDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, AppName)
}

// DataDir returns the XDG-compliant data directory for acorn
func DataDir() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, AppName)
}

// CacheDir returns the XDG-compliant cache directory for acorn
func CacheDir() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, _ := os.UserHomeDir()
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, AppName)
}

// EnsureDirs creates all XDG directories if they don't exist
func EnsureDirs() error {
	dirs := []string{ConfigDir(), DataDir(), CacheDir()}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
