package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings.DefaultOutputFormat != "table" {
		t.Errorf("expected default output format 'table', got '%s'", settings.DefaultOutputFormat)
	}

	if settings.RequestsPerSecond != 5.0 {
		t.Errorf("expected requests per second 5.0, got %.1f", settings.RequestsPerSecond)
	}

	if !settings.CacheEnabled {
		t.Error("expected cache to be enabled by default")
	}

	if settings.CacheTTL != 15 {
		t.Errorf("expected cache TTL 15, got %d", settings.CacheTTL)
	}

	if settings.TelemetryEnabled {
		t.Error("expected telemetry to be disabled by default")
	}
}

func TestConfig_AddInstance(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	instance := &Instance{
		Name: "test",
		URL:  "https://test.instructure.com",
	}

	err := cfg.AddInstance(instance)
	if err != nil {
		t.Fatalf("AddInstance failed: %v", err)
	}

	// Should set as default since it's the first
	if cfg.DefaultInstance != "test" {
		t.Errorf("expected default instance 'test', got '%s'", cfg.DefaultInstance)
	}

	// Verify instance was added
	if _, exists := cfg.Instances["test"]; !exists {
		t.Error("expected instance to exist")
	}
}

func TestConfig_AddInstance_Duplicate(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	instance := &Instance{
		Name: "test",
		URL:  "https://test.instructure.com",
	}

	err := cfg.AddInstance(instance)
	if err != nil {
		t.Fatalf("AddInstance failed: %v", err)
	}

	// Try to add same instance again
	err = cfg.AddInstance(instance)
	if err == nil {
		t.Error("expected error when adding duplicate instance")
	}
}

func TestConfig_GetInstance(t *testing.T) {
	cfg := &Config{
		Instances: map[string]*Instance{
			"test": {
				Name: "test",
				URL:  "https://test.instructure.com",
			},
		},
	}

	instance, err := cfg.GetInstance("test")
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	if instance.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", instance.Name)
	}

	// Try non-existent instance
	_, err = cfg.GetInstance("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent instance")
	}
}

func TestConfig_GetDefaultInstance(t *testing.T) {
	cfg := &Config{
		DefaultInstance: "test",
		Instances: map[string]*Instance{
			"test": {
				Name: "test",
				URL:  "https://test.instructure.com",
			},
		},
	}

	instance, err := cfg.GetDefaultInstance()
	if err != nil {
		t.Fatalf("GetDefaultInstance failed: %v", err)
	}

	if instance.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", instance.Name)
	}
}

func TestConfig_GetDefaultInstance_NotConfigured(t *testing.T) {
	cfg := &Config{
		Instances: make(map[string]*Instance),
	}

	_, err := cfg.GetDefaultInstance()
	if err == nil {
		t.Error("expected error when no default instance configured")
	}
}

func TestConfig_SetDefaultInstance(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances: map[string]*Instance{
			"test1": {Name: "test1", URL: "https://test1.com"},
			"test2": {Name: "test2", URL: "https://test2.com"},
		},
		configPath: configPath,
	}

	err := cfg.SetDefaultInstance("test2")
	if err != nil {
		t.Fatalf("SetDefaultInstance failed: %v", err)
	}

	if cfg.DefaultInstance != "test2" {
		t.Errorf("expected default 'test2', got '%s'", cfg.DefaultInstance)
	}
}

func TestConfig_SetDefaultInstance_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		configPath: configPath,
	}

	err := cfg.SetDefaultInstance("nonexistent")
	if err == nil {
		t.Error("expected error when setting non-existent instance as default")
	}
}

func TestConfig_RemoveInstance(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		DefaultInstance: "test1",
		Instances: map[string]*Instance{
			"test1": {Name: "test1", URL: "https://test1.com"},
			"test2": {Name: "test2", URL: "https://test2.com"},
		},
		configPath: configPath,
	}

	err := cfg.RemoveInstance("test1")
	if err != nil {
		t.Fatalf("RemoveInstance failed: %v", err)
	}

	// Should update default to another instance
	if cfg.DefaultInstance == "test1" {
		t.Error("expected default to change after removing default instance")
	}

	// Verify instance was removed
	if _, exists := cfg.Instances["test1"]; exists {
		t.Error("expected instance to be removed")
	}
}

func TestConfig_ListInstances(t *testing.T) {
	cfg := &Config{
		Instances: map[string]*Instance{
			"test1": {Name: "test1", URL: "https://test1.com"},
			"test2": {Name: "test2", URL: "https://test2.com"},
			"test3": {Name: "test3", URL: "https://test3.com"},
		},
	}

	instances := cfg.ListInstances()

	if len(instances) != 3 {
		t.Errorf("expected 3 instances, got %d", len(instances))
	}
}

func TestConfig_UpdateSettings(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	newSettings := &Settings{
		DefaultOutputFormat: "json",
		RequestsPerSecond:   10.0,
		CacheEnabled:        false,
		CacheTTL:            30,
		TelemetryEnabled:    true,
		LogLevel:            "debug",
	}

	err := cfg.UpdateSettings(newSettings)
	if err != nil {
		t.Fatalf("UpdateSettings failed: %v", err)
	}

	if cfg.Settings.DefaultOutputFormat != "json" {
		t.Errorf("expected format 'json', got '%s'", cfg.Settings.DefaultOutputFormat)
	}

	if cfg.Settings.RequestsPerSecond != 10.0 {
		t.Errorf("expected requests per second 10.0, got %.1f", cfg.Settings.RequestsPerSecond)
	}
}

func TestLoad_NewConfig(t *testing.T) {
	// Use a temporary directory that doesn't exist
	tempHome := filepath.Join(os.TempDir(), "canvas-cli-test-"+t.Name())
	defer os.RemoveAll(tempHome)

	// Set temporary HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should return default config
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.Instances == nil {
		t.Error("expected instances map to be initialized")
	}

	if cfg.Settings == nil {
		t.Error("expected settings to be initialized")
	}
}

func TestConfig_Save(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		DefaultInstance: "test",
		Instances: map[string]*Instance{
			"test": {
				Name: "test",
				URL:  "https://test.instructure.com",
			},
		},
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestConfig_UpdateInstance(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances: map[string]*Instance{
			"test": {
				Name: "test",
				URL:  "https://test.instructure.com",
			},
		},
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	// Update instance
	updated := &Instance{
		Name: "test",
		URL:  "https://updated.instructure.com",
	}

	err := cfg.UpdateInstance("test", updated)
	if err != nil {
		t.Fatalf("UpdateInstance failed: %v", err)
	}

	// Verify update
	instance, err := cfg.GetInstance("test")
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	if instance.URL != "https://updated.instructure.com" {
		t.Errorf("expected URL 'https://updated.instructure.com', got '%s'", instance.URL)
	}
}

func TestConfig_UpdateInstance_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	updated := &Instance{
		Name: "nonexistent",
		URL:  "https://test.instructure.com",
	}

	err := cfg.UpdateInstance("nonexistent", updated)
	if err == nil {
		t.Error("expected error when updating nonexistent instance")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if path == "" {
		t.Error("expected non-empty config path")
	}

	// Path should end with config.yaml
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected filename 'config.yaml', got '%s'", filepath.Base(path))
	}
}

func TestLoad_ExistingConfig(t *testing.T) {
	// Skip on Windows as HOME environment variable works differently
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows (test uses HOME environment variable)")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a config file
	cfg := &Config{
		DefaultInstance: "test",
		Instances: map[string]*Instance{
			"test": {
				Name: "test",
				URL:  "https://test.instructure.com",
			},
		},
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Change HOME to temp dir so Load finds our config
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Create .canvas-cli directory structure
	canvasDir := filepath.Join(tempDir, ".canvas-cli")
	os.MkdirAll(canvasDir, 0700)

	// Copy config to expected location
	expectedPath := filepath.Join(canvasDir, "config.yaml")
	data, _ := os.ReadFile(configPath)
	os.WriteFile(expectedPath, data, 0600)

	// Load config
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.DefaultInstance != "test" {
		t.Errorf("expected default instance 'test', got '%s'", loaded.DefaultInstance)
	}

	if _, exists := loaded.Instances["test"]; !exists {
		t.Error("expected 'test' instance to exist")
	}
}

func TestConfig_SaveAndLoad_RoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create config
	original := &Config{
		DefaultInstance: "prod",
		Instances: map[string]*Instance{
			"prod": {
				Name: "prod",
				URL:  "https://prod.instructure.com",
			},
			"dev": {
				Name: "dev",
				URL:  "https://dev.instructure.com",
			},
		},
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	// Save
	err := original.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}
	loaded.configPath = configPath

	// Verify
	if loaded.DefaultInstance != original.DefaultInstance {
		t.Errorf("expected default instance '%s', got '%s'",
			original.DefaultInstance, loaded.DefaultInstance)
	}

	if len(loaded.Instances) != len(original.Instances) {
		t.Errorf("expected %d instances, got %d",
			len(original.Instances), len(loaded.Instances))
	}
}

func TestConfig_GetInstance_NonExistent(t *testing.T) {
	cfg := &Config{
		Instances: make(map[string]*Instance),
	}

	_, err := cfg.GetInstance("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent instance")
	}
}

func TestConfig_ListInstances_Empty(t *testing.T) {
	cfg := &Config{
		Instances: make(map[string]*Instance),
	}

	instances := cfg.ListInstances()
	if len(instances) != 0 {
		t.Errorf("expected 0 instances, got %d", len(instances))
	}
}

func TestConfig_ListInstances_Multiple(t *testing.T) {
	cfg := &Config{
		Instances: map[string]*Instance{
			"prod": {Name: "prod", URL: "https://prod.instructure.com"},
			"dev":  {Name: "dev", URL: "https://dev.instructure.com"},
			"test": {Name: "test", URL: "https://test.instructure.com"},
		},
	}

	instances := cfg.ListInstances()
	if len(instances) != 3 {
		t.Errorf("expected 3 instances, got %d", len(instances))
	}
}

func TestConfig_UpdateSettings_Nil(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Settings:   DefaultSettings(),
		Instances:  make(map[string]*Instance),
		configPath: configPath,
	}

	// UpdateSettings doesn't validate nil - it just accepts it
	err := cfg.UpdateSettings(nil)
	if err != nil {
		t.Errorf("UpdateSettings failed: %v", err)
	}

	// Verify settings were updated to nil
	if cfg.Settings != nil {
		t.Error("expected settings to be nil after update")
	}
}

func TestSettings_Validation(t *testing.T) {
	settings := DefaultSettings()

	// Test valid settings
	if settings.RequestsPerSecond <= 0 {
		t.Error("requests per second should be positive")
	}

	if settings.CacheTTL < 0 {
		t.Error("cache TTL should not be negative")
	}

	if settings.DefaultOutputFormat == "" {
		t.Error("default output format should not be empty")
	}
}

func TestInstance_Validation(t *testing.T) {
	instance := &Instance{
		Name: "test",
		URL:  "https://test.instructure.com",
	}

	if instance.Name == "" {
		t.Error("instance name should not be empty")
	}

	if instance.URL == "" {
		t.Error("instance URL should not be empty")
	}
}

func TestConfig_RemoveInstance_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	err := cfg.RemoveInstance("nonexistent")
	// RemoveInstance returns error for nonexistent instances
	if err == nil {
		t.Error("expected error when removing nonexistent instance")
	}
}

func TestConfig_AddInstance_WithoutURL(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	instance := &Instance{
		Name: "test",
		URL:  "",
	}

	err := cfg.AddInstance(instance)
	if err == nil {
		t.Error("expected error when adding instance without URL")
	}
}

func TestConfig_AddInstance_WithoutName(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	instance := &Instance{
		Name: "",
		URL:  "https://test.instructure.com",
	}

	err := cfg.AddInstance(instance)
	if err == nil {
		t.Error("expected error when adding instance without name")
	}
}
