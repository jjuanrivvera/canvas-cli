package diagnostics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "test",
	}

	doctor := New(cfg, nil)
	if doctor == nil {
		t.Fatal("expected non-nil doctor")
	}

	if doctor.config != cfg {
		t.Error("expected config to be set")
	}
}

func TestCheckEnvironment(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkEnvironment(ctx)

	if check.Name != "Environment" {
		t.Errorf("expected name 'Environment', got '%s'", check.Name)
	}

	if check.Status != StatusPass {
		t.Errorf("expected status PASS, got %s", check.Status)
	}

	if check.Message == "" {
		t.Error("expected non-empty message")
	}

	// Skip duration check on Windows as timing can be too fast to measure
	if runtime.GOOS != "windows" && check.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}

func TestCheckConfig_NoConfig(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkConfig(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL, got %s", check.Status)
	}

	if check.Message != "No configuration found" {
		t.Errorf("unexpected message: %s", check.Message)
	}
}

func TestCheckConfig_NoDefaultInstance(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "",
		Instances:       map[string]*config.Instance{},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConfig(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL, got %s", check.Status)
	}
}

func TestCheckConfig_Success(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "test",
		Instances: map[string]*config.Instance{
			"test": {
				Name: "test",
				URL:  "https://canvas.example.com",
			},
		},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConfig(ctx)

	if check.Status != StatusPass {
		t.Errorf("expected status PASS, got %s: %s", check.Status, check.Message)
	}
}

func TestCheckConfig_GetDefaultInstanceError(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "nonexistent",
		Instances:       map[string]*config.Instance{},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConfig(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "Failed to get default instance") {
		t.Errorf("unexpected message: %s", check.Message)
	}
}

func TestCheckConfig_EmptyURL(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "test",
		Instances: map[string]*config.Instance{
			"test": {
				Name: "test",
				URL:  "",
			},
		},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConfig(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL, got %s", check.Status)
	}

	if check.Message != "Instance URL not configured" {
		t.Errorf("unexpected message: %s", check.Message)
	}
}

func TestCheckConnectivity_NoConfig(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkConnectivity(ctx)

	if check.Status != StatusSkipped {
		t.Errorf("expected status SKIP, got %s", check.Status)
	}
}

func TestCheckConnectivity_NoInstance(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "",
		Instances:       map[string]*config.Instance{},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConnectivity(ctx)

	if check.Status != StatusSkipped {
		t.Errorf("expected status SKIP, got %s", check.Status)
	}
}

func TestCheckConnectivity_Success(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		DefaultInstance: "test",
		Instances: map[string]*config.Instance{
			"test": {
				Name: "test",
				URL:  server.URL,
			},
		},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConnectivity(ctx)

	if check.Status != StatusPass {
		t.Errorf("expected status PASS, got %s: %s", check.Status, check.Message)
	}
}

func TestCheckAuthentication_NoClient(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkAuthentication(ctx)

	if check.Status != StatusSkipped {
		t.Errorf("expected status SKIP, got %s", check.Status)
	}
}

func TestCheckAuthentication_Success(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/users/self" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1, "name": "Test User"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, _ := api.NewClient(api.ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	doctor := New(nil, client)
	ctx := context.Background()

	check := doctor.checkAuthentication(ctx)

	if check.Status != StatusPass {
		t.Errorf("expected status PASS, got %s: %s", check.Status, check.Message)
	}
}

func TestCheckAPIAccess_NoClient(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkAPIAccess(ctx)

	if check.Status != StatusSkipped {
		t.Errorf("expected status SKIP, got %s", check.Status)
	}
}

func TestCheckAPIAccess_Success(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/courses" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, _ := api.NewClient(api.ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	doctor := New(nil, client)
	ctx := context.Background()

	check := doctor.checkAPIAccess(ctx)

	if check.Status != StatusPass {
		t.Errorf("expected status PASS, got %s: %s", check.Status, check.Message)
	}
}

func TestCheckDiskSpace(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkDiskSpace(ctx)

	// Should at least not fail
	if check.Status == StatusFail {
		t.Errorf("unexpected FAIL status: %s", check.Message)
	}

	if check.Name != "Disk Space" {
		t.Errorf("expected name 'Disk Space', got '%s'", check.Name)
	}
}

func TestCheckPermissions(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkPermissions(ctx)

	// Should at least not fail
	if check.Status == StatusFail {
		t.Errorf("unexpected FAIL status: %s", check.Message)
	}

	if check.Name != "Permissions" {
		t.Errorf("expected name 'Permissions', got '%s'", check.Name)
	}
}

func TestCheckPermissions_WithConfigDir(t *testing.T) {
	// Create temp directory with correct permissions
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".canvas-cli")
	os.MkdirAll(configDir, 0700)

	// Temporarily change HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkPermissions(ctx)

	if check.Status == StatusFail {
		t.Errorf("unexpected FAIL status: %s", check.Message)
	}
}

func TestCheckPermissions_InsecurePerms(t *testing.T) {
	// Create temp directory with insecure permissions
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".canvas-cli")
	os.MkdirAll(configDir, 0755) // Insecure: world-readable

	// Temporarily change HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkPermissions(ctx)

	if check.Status != StatusWarning {
		t.Errorf("expected status WARNING for insecure permissions, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "insecure permissions") {
		t.Errorf("expected message about insecure permissions, got: %s", check.Message)
	}
}

func TestCheckPermissions_NoConfigDir(t *testing.T) {
	// Skip on Windows as HOME environment variable works differently
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows (test uses HOME environment variable)")
	}

	// Create temp directory WITHOUT .canvas-cli subdirectory
	tempDir := t.TempDir()

	// Temporarily change HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkPermissions(ctx)

	if check.Status != StatusWarning {
		t.Errorf("expected status WARNING when config dir doesn't exist, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "does not exist") {
		t.Errorf("expected message about directory not existing, got: %s", check.Message)
	}
}

func TestRun_AllChecks(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/users/self":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1, "name": "Test User"}`))
		case "/api/v1/courses":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		DefaultInstance: "test",
		Instances: map[string]*config.Instance{
			"test": {
				Name: "test",
				URL:  server.URL,
			},
		},
	}

	client, _ := api.NewClient(api.ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	doctor := New(cfg, client)
	ctx := context.Background()

	report, err := doctor.Run(ctx)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if report == nil {
		t.Fatal("expected non-nil report")
	}

	// Should have 7 checks
	if len(report.Checks) != 7 {
		t.Errorf("expected 7 checks, got %d", len(report.Checks))
	}

	// At least environment check should pass
	if report.PassCount == 0 {
		t.Error("expected at least one passing check")
	}

	if report.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}

func TestRun_NoConfig(t *testing.T) {
	doctor := New(nil, nil)
	ctx := context.Background()

	report, err := doctor.Run(ctx)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Should have some checks
	if len(report.Checks) == 0 {
		t.Error("expected some checks")
	}

	// Should have failures due to missing config
	if report.FailCount == 0 {
		t.Error("expected some failures without config")
	}
}

func TestCheckStatus_String(t *testing.T) {
	tests := []struct {
		status   CheckStatus
		expected string
	}{
		{StatusPass, "PASS"},
		{StatusFail, "FAIL"},
		{StatusWarning, "WARN"},
		{StatusSkipped, "SKIP"},
	}

	for _, tt := range tests {
		result := tt.status.String()
		if result != tt.expected {
			t.Errorf("expected '%s', got '%s'", tt.expected, result)
		}
	}
}

func TestReport_IsHealthy(t *testing.T) {
	// Healthy report
	report := &Report{
		PassCount: 5,
		FailCount: 0,
	}

	if !report.IsHealthy() {
		t.Error("expected report to be healthy")
	}

	// Unhealthy report
	report2 := &Report{
		PassCount: 3,
		FailCount: 2,
	}

	if report2.IsHealthy() {
		t.Error("expected report to be unhealthy")
	}
}

func TestReport_Summary(t *testing.T) {
	report := &Report{
		Checks: []Check{
			{Status: StatusPass},
			{Status: StatusPass},
			{Status: StatusFail},
			{Status: StatusWarning},
			{Status: StatusSkipped},
		},
		PassCount: 2,
		FailCount: 1,
		WarnCount: 1,
		SkipCount: 1,
	}

	summary := report.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Check that summary contains counts
	expected := "Total: 5, Pass: 2, Fail: 1, Warn: 1, Skip: 1"
	if summary != expected {
		t.Errorf("expected summary '%s', got '%s'", expected, summary)
	}
}

func TestCheck_Structure(t *testing.T) {
	check := Check{
		Name:        "Test Check",
		Description: "Test description",
		Status:      StatusPass,
		Message:     "Test message",
		Duration:    100 * time.Millisecond,
		Error:       nil,
	}

	if check.Name != "Test Check" {
		t.Errorf("expected name 'Test Check', got '%s'", check.Name)
	}

	if check.Status != StatusPass {
		t.Errorf("expected status PASS, got %s", check.Status)
	}

	if check.Duration != 100*time.Millisecond {
		t.Errorf("expected duration 100ms, got %v", check.Duration)
	}
}

func TestReport_StartTime(t *testing.T) {
	report := &Report{
		StartTime: time.Now(),
		Checks:    []Check{},
	}

	if report.StartTime.IsZero() {
		t.Error("expected non-zero start time")
	}
}

func TestCheckAuthentication_Failure(t *testing.T) {
	// Create test HTTP server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client, _ := api.NewClient(api.ClientConfig{
		BaseURL: server.URL,
		Token:   "invalid-token",
	})

	doctor := New(nil, client)
	ctx := context.Background()

	check := doctor.checkAuthentication(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL, got %s", check.Status)
	}

	if check.Error == nil {
		t.Error("expected error to be set")
	}
}

func TestCheckAPIAccess_Failure(t *testing.T) {
	// Create test HTTP server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client, _ := api.NewClient(api.ClientConfig{
		BaseURL: server.URL,
		Token:   "invalid-token",
	})

	doctor := New(nil, client)
	ctx := context.Background()

	check := doctor.checkAPIAccess(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL, got %s", check.Status)
	}
}

func TestCheckConnectivity_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		DefaultInstance: "test",
		Instances: map[string]*config.Instance{
			"test": {
				Name: "test",
				URL:  "http://invalid-url-that-does-not-exist.example.com",
			},
		},
	}

	doctor := New(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	check := doctor.checkConnectivity(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL for invalid URL, got %s", check.Status)
	}
}

func TestCheckDiskSpace_MkdirError(t *testing.T) {
	// Skip on Windows as HOME environment variable works differently
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows (test uses HOME environment variable)")
	}

	// Create a file where the cache directory should be, causing MkdirAll to fail
	tempDir := t.TempDir()

	// Create a file (not directory) at the cache path
	cacheFile := filepath.Join(tempDir, ".canvas-cli", "cache")
	os.MkdirAll(filepath.Dir(cacheFile), 0700)
	os.WriteFile(cacheFile, []byte("blocking file"), 0600)

	// Temporarily change HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkDiskSpace(ctx)

	if check.Status != StatusWarning {
		t.Errorf("expected status WARNING when cache dir not writable, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "not writable") {
		t.Errorf("expected message about directory not writable, got: %s", check.Message)
	}
}

func TestCheckPermissions_SecurePerms(t *testing.T) {
	// Skip on Windows as HOME environment variable and permission checks work differently
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows (test uses HOME environment variable and chmod)")
	}

	// Create temp directory with secure permissions (0700)
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".canvas-cli")
	os.MkdirAll(configDir, 0700)

	// Temporarily change HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkPermissions(ctx)

	if check.Status != StatusPass {
		t.Errorf("expected status PASS for secure permissions, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "secure") {
		t.Errorf("expected message about secure permissions, got: %s", check.Message)
	}
}

func TestCheckConnectivity_InvalidURLFormat(t *testing.T) {
	// Use a URL with invalid characters to trigger NewRequestWithContext error
	cfg := &config.Config{
		DefaultInstance: "test",
		Instances: map[string]*config.Instance{
			"test": {
				Name: "test",
				URL:  "ht\ttp://invalid", // Contains control character
			},
		},
	}

	doctor := New(cfg, nil)
	ctx := context.Background()

	check := doctor.checkConnectivity(ctx)

	if check.Status != StatusFail {
		t.Errorf("expected status FAIL for malformed URL, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "Failed to create request") {
		t.Errorf("expected message about request creation failure, got: %s", check.Message)
	}
}

func TestCheckDiskSpace_HomeDirError(t *testing.T) {
	// Temporarily set HOME to empty to trigger UserHomeDir error
	oldHome := os.Getenv("HOME")
	os.Unsetenv("HOME")
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkDiskSpace(ctx)

	if check.Status != StatusWarning {
		t.Errorf("expected status WARNING when home dir unavailable, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "home directory") {
		t.Errorf("expected message about home directory, got: %s", check.Message)
	}
}

func TestCheckPermissions_HomeDirError(t *testing.T) {
	// Temporarily set HOME to empty to trigger UserHomeDir error
	oldHome := os.Getenv("HOME")
	os.Unsetenv("HOME")
	defer os.Setenv("HOME", oldHome)

	doctor := New(nil, nil)
	ctx := context.Background()

	check := doctor.checkPermissions(ctx)

	if check.Status != StatusWarning {
		t.Errorf("expected status WARNING when home dir unavailable, got %s", check.Status)
	}

	if !strings.Contains(check.Message, "home directory") {
		t.Errorf("expected message about home directory, got: %s", check.Message)
	}
}
