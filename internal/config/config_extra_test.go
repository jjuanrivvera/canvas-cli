package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// ---------- Clone -----------------------------------------------------------

func TestConfig_Clone_Nil(t *testing.T) {
	var c *Config
	if c.Clone() != nil {
		t.Error("Clone of nil Config should return nil")
	}
}

func TestConfig_Clone_Full(t *testing.T) {
	original := &Config{
		DefaultInstance: "prod",
		Instances: map[string]*Instance{
			"prod": {Name: "prod", URL: "https://prod.example.com"},
		},
		Settings: DefaultSettings(),
		Aliases:  map[string]string{"ls": "courses list"},
		Context:  &Context{CourseID: 99},
	}

	clone := original.Clone()

	if clone == original {
		t.Error("Clone should return a different pointer")
	}
	if clone.DefaultInstance != original.DefaultInstance {
		t.Errorf("DefaultInstance mismatch: got %q", clone.DefaultInstance)
	}
	if clone.Instances["prod"] == original.Instances["prod"] {
		t.Error("Instances entries should be copies, not same pointer")
	}
	if clone.Settings == original.Settings {
		t.Error("Settings should be a copy, not same pointer")
	}
	if clone.Aliases["ls"] != "courses list" {
		t.Errorf("Alias not copied, got %q", clone.Aliases["ls"])
	}
	if clone.Context == original.Context {
		t.Error("Context should be a copy, not same pointer")
	}
	if clone.Context.CourseID != 99 {
		t.Errorf("Context.CourseID not copied, got %d", clone.Context.CourseID)
	}
}

func TestConfig_Clone_NilFields(t *testing.T) {
	original := &Config{
		DefaultInstance: "dev",
		Instances:       nil,
		Settings:        nil,
		Aliases:         nil,
		Context:         nil,
	}
	clone := original.Clone()
	if clone.Instances != nil {
		t.Error("expected nil Instances in clone")
	}
	if clone.Settings != nil {
		t.Error("expected nil Settings in clone")
	}
	if clone.Aliases != nil {
		t.Error("expected nil Aliases in clone")
	}
	if clone.Context != nil {
		t.Error("expected nil Context in clone")
	}
}

// ---------- Reload ----------------------------------------------------------

func TestReload(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".canvas-cli", "config.yaml")

	// Write initial config via Save
	cfg := &Config{
		DefaultInstance: "reload-test",
		Instances: map[string]*Instance{
			"reload-test": {Name: "reload-test", URL: "https://reload.example.com"},
		},
		Settings:   DefaultSettings(),
		configPath: configPath,
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Setenv("HOME", tempDir)
	ResetCache()
	defer ResetCache()

	loaded, err := Reload()
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}
	if loaded.DefaultInstance != "reload-test" {
		t.Errorf("expected 'reload-test', got %q", loaded.DefaultInstance)
	}
}

// ---------- ResetCache ------------------------------------------------------

func TestResetCache(t *testing.T) {
	// Prime the cache
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	ResetCache()
	defer ResetCache()

	// Load populates cache
	if _, err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Reset clears it — verify by checking cachedConfig is nil after reset
	ResetCache()
	cacheMu.Lock()
	isNil := cachedConfig == nil
	cacheMu.Unlock()
	if !isNil {
		t.Error("expected cachedConfig to be nil after ResetCache")
	}
}

// ---------- Load (cache path) -----------------------------------------------

func TestLoad_CachePath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	ResetCache()
	defer ResetCache()

	// First load — populates cache
	cfg1, err := Load()
	if err != nil {
		t.Fatalf("first Load: %v", err)
	}
	// Second load — returns clone from cache
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}
	// They must be equal in value but different pointers
	if cfg1 == cfg2 {
		t.Error("Load should return clones, not the same pointer")
	}
	if cfg1.DefaultInstance != cfg2.DefaultInstance {
		t.Errorf("DefaultInstance mismatch between cache loads")
	}
}

// ---------- GetConfigDir / GetConfigPath ------------------------------------

func TestGetConfigDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	expected := filepath.Join(tempDir, ".canvas-cli")
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

// ---------- Aliases ---------------------------------------------------------

func TestConfig_Aliases_SetGetListDelete(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	// SetAlias
	if err := cfg.SetAlias("ls", "courses list"); err != nil {
		t.Fatalf("SetAlias: %v", err)
	}
	if err := cfg.SetAlias("enroll", "enrollments list"); err != nil {
		t.Fatalf("SetAlias second: %v", err)
	}

	// GetAlias — found
	val, ok := cfg.GetAlias("ls")
	if !ok || val != "courses list" {
		t.Errorf("GetAlias('ls') = %q, %v; want 'courses list', true", val, ok)
	}

	// GetAlias — not found
	_, ok = cfg.GetAlias("nonexistent")
	if ok {
		t.Error("expected GetAlias to return false for nonexistent key")
	}

	// ListAliases
	aliases := cfg.ListAliases()
	if len(aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(aliases))
	}

	// DeleteAlias — success
	if err := cfg.DeleteAlias("ls"); err != nil {
		t.Fatalf("DeleteAlias: %v", err)
	}
	if _, ok := cfg.GetAlias("ls"); ok {
		t.Error("expected alias to be deleted")
	}

	// DeleteAlias — nonexistent
	if err := cfg.DeleteAlias("nonexistent"); err == nil {
		t.Error("expected error when deleting nonexistent alias")
	}
}

func TestConfig_ListAliases_Nil(t *testing.T) {
	cfg := &Config{Aliases: nil}
	aliases := cfg.ListAliases()
	if aliases == nil || len(aliases) != 0 {
		t.Errorf("expected empty map, got %v", aliases)
	}
}

func TestConfig_GetAlias_NilMap(t *testing.T) {
	cfg := &Config{Aliases: nil}
	_, ok := cfg.GetAlias("x")
	if ok {
		t.Error("expected false when Aliases map is nil")
	}
}

func TestConfig_DeleteAlias_NilMap(t *testing.T) {
	cfg := &Config{Aliases: nil}
	err := cfg.DeleteAlias("x")
	if err == nil {
		t.Error("expected error when Aliases map is nil")
	}
}

// ---------- Context ---------------------------------------------------------

func TestConfig_Context_SetGetClear(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	cfg := &Config{
		Instances:  make(map[string]*Instance),
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	// GetContext on nil → returns empty Context
	ctx := cfg.GetContext()
	if ctx == nil {
		t.Fatal("GetContext should never return nil")
	}
	if ctx.CourseID != 0 {
		t.Errorf("expected zero CourseID, got %d", ctx.CourseID)
	}

	// SetContext
	newCtx := &Context{CourseID: 42, AssignmentID: 7}
	if err := cfg.SetContext(newCtx); err != nil {
		t.Fatalf("SetContext: %v", err)
	}
	got := cfg.GetContext()
	if got.CourseID != 42 {
		t.Errorf("expected CourseID 42, got %d", got.CourseID)
	}

	// ClearContext
	if err := cfg.ClearContext(); err != nil {
		t.Fatalf("ClearContext: %v", err)
	}
	if cfg.Context != nil {
		t.Error("expected Context to be nil after clear")
	}
}

// ---------- RemoveInstance (last instance clears default) -------------------

func TestConfig_RemoveInstance_LastInstance(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		DefaultInstance: "only",
		Instances: map[string]*Instance{
			"only": {Name: "only", URL: "https://only.example.com"},
		},
		Settings:   DefaultSettings(),
		configPath: configPath,
	}

	if err := cfg.RemoveInstance("only"); err != nil {
		t.Fatalf("RemoveInstance: %v", err)
	}
	if cfg.DefaultInstance != "" {
		t.Errorf("expected empty DefaultInstance after removing last instance, got %q", cfg.DefaultInstance)
	}
	if len(cfg.Instances) != 0 {
		t.Errorf("expected 0 instances, got %d", len(cfg.Instances))
	}
}

// ---------- loadFromDisk (corrupt YAML) -------------------------------------

func TestLoadFromDisk_CorruptYAML(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}

	tempDir := t.TempDir()
	canvasDir := filepath.Join(tempDir, ".canvas-cli")
	if err := os.MkdirAll(canvasDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	// Write invalid YAML (unclosed brace triggers a parse error in go-yaml v3)
	configPath := filepath.Join(canvasDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("[invalid: {unclosed"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	t.Setenv("HOME", tempDir)
	ResetCache()
	defer ResetCache()

	_, err := Load()
	if err == nil {
		t.Error("expected error when loading corrupt YAML")
	}
}

// ---------- ValidateToken ---------------------------------------------------

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		errHint string
	}{
		{"empty", "", true, "empty"},
		{"whitespace only", "   ", true, "whitespace"},
		{"too short", "short", true, "too short"},
		{"too long", makeString('a', 501), true, "too long"},
		{"placeholder: token", "token", true, "placeholder"},
		{"placeholder: changeme", "changeme", true, "placeholder"},
		{"valid token", "7~abcdefghijklmnopqrstuvwxyz123456", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for token %q, got nil", tt.token)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for token %q: %v", tt.token, err)
				}
			}
		})
	}
}

func makeString(r rune, n int) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = r
	}
	return string(runes)
}

// ---------- ValidateInstance with token and OAuth ---------------------------

func TestValidateInstance_WithToken(t *testing.T) {
	inst := &Instance{
		Name:  "test",
		URL:   "https://canvas.example.com",
		Token: "7~averylongenoughtokenfortesting1234",
	}
	if err := ValidateInstance(inst); err != nil {
		t.Errorf("unexpected error for valid token: %v", err)
	}
}

func TestValidateInstance_WithShortToken(t *testing.T) {
	inst := &Instance{
		Name:  "test",
		URL:   "https://canvas.example.com",
		Token: "shorttoken",
	}
	err := ValidateInstance(inst)
	if err == nil {
		t.Error("expected error for too-short token")
	}
}

func TestValidateInstance_WithOAuth(t *testing.T) {
	inst := &Instance{
		Name:         "test",
		URL:          "https://canvas.example.com",
		ClientID:     "1234567890abcdef",
		ClientSecret: "secret1234567890",
	}
	if err := ValidateInstance(inst); err != nil {
		t.Errorf("unexpected error for valid OAuth: %v", err)
	}
}

func TestValidateInstance_ClientIDTooShort(t *testing.T) {
	inst := &Instance{
		Name:         "test",
		URL:          "https://canvas.example.com",
		ClientID:     "short1234567890",
		ClientSecret: "short",
	}
	// ClientSecret too short
	if err := ValidateInstance(inst); err == nil {
		t.Error("expected error for too-short client_secret")
	}
}

func TestValidateInstance_PublicClient(t *testing.T) {
	inst := &Instance{
		Name:         "test",
		URL:          "https://canvas.example.com",
		ClientID:     "1234567890abcdef",
		PublicClient: true,
	}
	if err := ValidateInstance(inst); err != nil {
		t.Errorf("unexpected error for public client without secret: %v", err)
	}
}

func TestValidateInstance_PublicClientWithSecret(t *testing.T) {
	inst := &Instance{
		Name:         "test",
		URL:          "https://canvas.example.com",
		ClientID:     "1234567890abcdef",
		ClientSecret: "secret1234567890",
		PublicClient: true,
	}
	if err := ValidateInstance(inst); err == nil {
		t.Error("expected error when public_client is combined with client_secret")
	}
}

func TestValidateInstance_PublicClientWithoutClientID(t *testing.T) {
	inst := &Instance{
		Name:         "test",
		URL:          "https://canvas.example.com",
		PublicClient: true,
	}
	if err := ValidateInstance(inst); err == nil {
		t.Error("expected error when public_client is set without client_id")
	}
}

func TestInstance_PublicClient_HasOAuthAndAuthType(t *testing.T) {
	inst := &Instance{
		Name:         "test",
		URL:          "https://canvas.example.com",
		ClientID:     "1234567890abcdef",
		PublicClient: true,
	}
	if !inst.HasOAuth() {
		t.Error("HasOAuth() = false for public client with client ID, want true")
	}
	if got := inst.AuthType(); got != "oauth" {
		t.Errorf("AuthType() = %q for public client, want %q", got, "oauth")
	}

	// Without the public-client marker a bare client ID is still not OAuth.
	inst.PublicClient = false
	if inst.HasOAuth() {
		t.Error("HasOAuth() = true for bare client ID without public_client, want false")
	}
}

func TestValidateInstance_NegativeAccountID(t *testing.T) {
	inst := &Instance{
		Name:             "test",
		URL:              "https://canvas.example.com",
		DefaultAccountID: -5,
	}
	if err := ValidateInstance(inst); err == nil {
		t.Error("expected error for negative DefaultAccountID")
	}
}

func TestConfig_Validate_WithNilSettings(t *testing.T) {
	cfg := &Config{
		Instances: map[string]*Instance{
			"test": {Name: "test", URL: "https://canvas.example.com"},
		},
		Settings: nil,
	}
	// nil Settings → ValidateSettings is skipped
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error with nil Settings: %v", err)
	}
}
