package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	cachedConfig *Config
	cacheMu      sync.Mutex
)

// Config represents the application configuration
type Config struct {
	// DefaultInstance is the default instance to use
	DefaultInstance string `yaml:"default_instance"`

	// Instances holds all configured Canvas instances
	Instances map[string]*Instance `yaml:"instances"`

	// Settings holds global settings
	Settings *Settings `yaml:"settings"`

	// Aliases holds user-defined command aliases
	Aliases map[string]string `yaml:"aliases,omitempty"`

	// Context holds the current context (course_id, assignment_id, etc.)
	Context *Context `yaml:"context,omitempty"`

	// configPath is the path to the config file
	configPath string
}

// Context holds the current working context
type Context struct {
	CourseID     int64 `yaml:"course_id,omitempty"`
	AssignmentID int64 `yaml:"assignment_id,omitempty"`
	UserID       int64 `yaml:"user_id,omitempty"`
	AccountID    int64 `yaml:"account_id,omitempty"`
}

// Instance represents a Canvas instance configuration
type Instance struct {
	Name             string `yaml:"name"`
	URL              string `yaml:"url"`
	ClientID         string `yaml:"client_id,omitempty"`
	ClientSecret     string `yaml:"client_secret,omitempty"`
	Token            string `yaml:"token,omitempty"` // API access token (alternative to OAuth)
	Description      string `yaml:"description,omitempty"`
	DefaultAccountID int64  `yaml:"default_account_id,omitempty"` // Default account ID for API calls
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

// HasDefaultAccountID returns true if the instance has a default account ID configured
func (i *Instance) HasDefaultAccountID() bool {
	return i.DefaultAccountID > 0
}

// Settings holds global application settings
type Settings struct {
	DefaultOutputFormat   string  `yaml:"default_output_format"`
	RequestsPerSecond     float64 `yaml:"requests_per_second"`
	CacheEnabled          bool    `yaml:"cache_enabled"`
	CacheTTL              int     `yaml:"cache_ttl_minutes"`
	TelemetryEnabled      bool    `yaml:"telemetry_enabled"`
	LogLevel              string  `yaml:"log_level"`
	AutoUpdateEnabled     bool    `yaml:"auto_update_enabled"`
	AutoUpdateIntervalMin int     `yaml:"auto_update_interval_minutes"`
}

// DefaultSettings returns the default settings
func DefaultSettings() *Settings {
	return &Settings{
		DefaultOutputFormat:   "table",
		RequestsPerSecond:     5.0,
		CacheEnabled:          true,
		CacheTTL:              15,
		TelemetryEnabled:      false,
		LogLevel:              "info",
		AutoUpdateEnabled:     true,
		AutoUpdateIntervalMin: 60, // Check every hour
	}
}

// Clone returns a deep copy of the Config.
func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}

	clone := &Config{
		DefaultInstance: c.DefaultInstance,
		configPath:      c.configPath,
	}

	// Clone Instances map
	if c.Instances != nil {
		clone.Instances = make(map[string]*Instance, len(c.Instances))
		for k, v := range c.Instances {
			if v != nil {
				inst := *v
				clone.Instances[k] = &inst
			}
		}
	}

	// Clone Settings
	if c.Settings != nil {
		s := *c.Settings
		clone.Settings = &s
	}

	// Clone Aliases map
	if c.Aliases != nil {
		clone.Aliases = make(map[string]string, len(c.Aliases))
		for k, v := range c.Aliases {
			clone.Aliases[k] = v
		}
	}

	// Clone Context
	if c.Context != nil {
		ctx := *c.Context
		clone.Context = &ctx
	}

	return clone
}

// GetConfigDir returns the path to the config directory (~/.canvas-cli)
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".canvas-cli"), nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// Load loads the configuration from the config file.
// The result is cached after the first successful load; subsequent calls
// return a clone of the cached value. Use Reload() to force a fresh read.
// The returned Config is safe to mutate without affecting other callers.
func Load() (*Config, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cachedConfig != nil {
		return cachedConfig.Clone(), nil
	}

	cfg, err := loadFromDisk()
	if err != nil {
		return nil, err
	}

	cachedConfig = cfg
	return cachedConfig.Clone(), nil
}

// Reload forces a fresh read of the config file, ignoring the cache.
// The returned Config is safe to mutate without affecting other callers.
func Reload() (*Config, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	cfg, err := loadFromDisk()
	if err != nil {
		return nil, err
	}

	cachedConfig = cfg
	return cachedConfig.Clone(), nil
}

// ResetCache clears the cached config. Intended for tests.
func ResetCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cachedConfig = nil
}

// loadFromDisk reads and parses the config file from disk.
func loadFromDisk() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Instances:  make(map[string]*Instance),
			Settings:   DefaultSettings(),
			configPath: configPath,
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.configPath = configPath

	if config.Instances == nil {
		config.Instances = make(map[string]*Instance)
	}

	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}

	if config.Settings == nil {
		config.Settings = DefaultSettings()
	}

	return &config, nil
}

// Save saves the configuration to the config file and invalidates the cache.
func (c *Config) Save() error {
	configDir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Invalidate cache so next Load() picks up the saved changes
	cacheMu.Lock()
	cachedConfig = c
	cacheMu.Unlock()

	return nil
}

// AddInstance adds a new instance to the configuration
func (c *Config) AddInstance(instance *Instance) error {
	// Validate instance configuration
	if err := ValidateInstance(instance); err != nil {
		return fmt.Errorf("invalid instance configuration: %w", err)
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
	// Validate instance configuration
	if err := ValidateInstance(instance); err != nil {
		return fmt.Errorf("invalid instance configuration: %w", err)
	}

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

// SetDefaultAccountID sets the default account ID for an instance
func (c *Config) SetDefaultAccountID(instanceName string, accountID int64) error {
	instance, exists := c.Instances[instanceName]
	if !exists {
		return fmt.Errorf("instance %q does not exist", instanceName)
	}

	instance.DefaultAccountID = accountID
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

// SetAlias creates or updates an alias
func (c *Config) SetAlias(name, expansion string) error {
	if c.Aliases == nil {
		c.Aliases = make(map[string]string)
	}
	c.Aliases[name] = expansion
	return c.Save()
}

// GetAlias retrieves an alias by name
func (c *Config) GetAlias(name string) (string, bool) {
	if c.Aliases == nil {
		return "", false
	}
	expansion, exists := c.Aliases[name]
	return expansion, exists
}

// DeleteAlias removes an alias
func (c *Config) DeleteAlias(name string) error {
	if c.Aliases == nil {
		return fmt.Errorf("alias %q does not exist", name)
	}
	if _, exists := c.Aliases[name]; !exists {
		return fmt.Errorf("alias %q does not exist", name)
	}
	delete(c.Aliases, name)
	return c.Save()
}

// ListAliases returns all configured aliases
func (c *Config) ListAliases() map[string]string {
	if c.Aliases == nil {
		return make(map[string]string)
	}
	return c.Aliases
}

// SetContext updates the current context
func (c *Config) SetContext(ctx *Context) error {
	c.Context = ctx
	return c.Save()
}

// GetContext returns the current context
func (c *Config) GetContext() *Context {
	if c.Context == nil {
		return &Context{}
	}
	return c.Context
}

// ClearContext clears the current context
func (c *Config) ClearContext() error {
	c.Context = nil
	return c.Save()
}
