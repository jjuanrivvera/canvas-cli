package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// DefaultInstance is the default instance to use
	DefaultInstance string `yaml:"default_instance"`

	// Instances holds all configured Canvas instances
	Instances map[string]*Instance `yaml:"instances"`

	// Settings holds global settings
	Settings *Settings `yaml:"settings"`

	// configPath is the path to the config file
	configPath string
}

// Instance represents a Canvas instance configuration
type Instance struct {
	Name         string `yaml:"name"`
	URL          string `yaml:"url"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
	Token        string `yaml:"token,omitempty"` // API access token (alternative to OAuth)
	Description  string `yaml:"description,omitempty"`
}

// HasToken returns true if the instance has an API token configured
func (i *Instance) HasToken() bool {
	return i.Token != ""
}

// HasOAuth returns true if the instance has OAuth credentials configured
func (i *Instance) HasOAuth() bool {
	return i.ClientID != "" && i.ClientSecret != ""
}

// AuthType returns a string describing the authentication type
func (i *Instance) AuthType() string {
	if i.HasToken() {
		return "token"
	}
	if i.HasOAuth() {
		return "oauth"
	}
	return "none"
}

// Settings holds global application settings
type Settings struct {
	DefaultOutputFormat string  `yaml:"default_output_format"`
	RequestsPerSecond   float64 `yaml:"requests_per_second"`
	CacheEnabled        bool    `yaml:"cache_enabled"`
	CacheTTL            int     `yaml:"cache_ttl_minutes"`
	TelemetryEnabled    bool    `yaml:"telemetry_enabled"`
	LogLevel            string  `yaml:"log_level"`
}

// DefaultSettings returns the default settings
func DefaultSettings() *Settings {
	return &Settings{
		DefaultOutputFormat: "table",
		RequestsPerSecond:   5.0,
		CacheEnabled:        true,
		CacheTTL:            15,
		TelemetryEnabled:    false,
		LogLevel:            "info",
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".canvas-cli")
	return filepath.Join(configDir, "config.yaml"), nil
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config doesn't exist, return default
		return &Config{
			Instances:  make(map[string]*Instance),
			Settings:   DefaultSettings(),
			configPath: configPath,
		}, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.configPath = configPath

	// Ensure instances map exists
	if config.Instances == nil {
		config.Instances = make(map[string]*Instance)
	}

	// Ensure settings exist
	if config.Settings == nil {
		config.Settings = DefaultSettings()
	}

	return &config, nil
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	// Ensure config directory exists
	configDir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(c.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddInstance adds a new instance to the configuration
func (c *Config) AddInstance(instance *Instance) error {
	if instance.Name == "" {
		return fmt.Errorf("instance name is required")
	}

	if instance.URL == "" {
		return fmt.Errorf("instance URL is required")
	}

	// Check if instance already exists
	if _, exists := c.Instances[instance.Name]; exists {
		return fmt.Errorf("instance %q already exists", instance.Name)
	}

	c.Instances[instance.Name] = instance

	// If this is the first instance, set it as default
	if c.DefaultInstance == "" {
		c.DefaultInstance = instance.Name
	}

	return c.Save()
}

// UpdateInstance updates an existing instance
func (c *Config) UpdateInstance(name string, instance *Instance) error {
	if _, exists := c.Instances[name]; !exists {
		return fmt.Errorf("instance %q does not exist", name)
	}

	c.Instances[name] = instance
	return c.Save()
}

// RemoveInstance removes an instance from the configuration
func (c *Config) RemoveInstance(name string) error {
	if _, exists := c.Instances[name]; !exists {
		return fmt.Errorf("instance %q does not exist", name)
	}

	delete(c.Instances, name)

	// If this was the default instance, clear it
	if c.DefaultInstance == name {
		c.DefaultInstance = ""
		// Set a new default if instances remain
		for instanceName := range c.Instances {
			c.DefaultInstance = instanceName
			break
		}
	}

	return c.Save()
}

// GetInstance retrieves an instance by name
func (c *Config) GetInstance(name string) (*Instance, error) {
	instance, exists := c.Instances[name]
	if !exists {
		return nil, fmt.Errorf("instance %q does not exist", name)
	}
	return instance, nil
}

// GetDefaultInstance returns the default instance
func (c *Config) GetDefaultInstance() (*Instance, error) {
	if c.DefaultInstance == "" {
		return nil, fmt.Errorf("no default instance configured")
	}
	return c.GetInstance(c.DefaultInstance)
}

// SetDefaultInstance sets the default instance
func (c *Config) SetDefaultInstance(name string) error {
	if _, exists := c.Instances[name]; !exists {
		return fmt.Errorf("instance %q does not exist", name)
	}

	c.DefaultInstance = name
	return c.Save()
}

// ListInstances returns all configured instances
func (c *Config) ListInstances() []*Instance {
	instances := make([]*Instance, 0, len(c.Instances))
	for _, instance := range c.Instances {
		instances = append(instances, instance)
	}
	return instances
}

// UpdateSettings updates the global settings
func (c *Config) UpdateSettings(settings *Settings) error {
	c.Settings = settings
	return c.Save()
}
