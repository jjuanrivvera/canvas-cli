package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// Check represents a diagnostic check
type Check struct {
	Name        string
	Description string
	Status      CheckStatus
	Message     string
	Duration    time.Duration
	Error       error
}

// CheckStatus represents the status of a diagnostic check
type CheckStatus string

const (
	StatusPass    CheckStatus = "PASS"
	StatusFail    CheckStatus = "FAIL"
	StatusWarning CheckStatus = "WARN"
	StatusSkipped CheckStatus = "SKIP"
)

// Report represents a collection of diagnostic checks
type Report struct {
	Checks    []Check
	StartTime time.Time
	Duration  time.Duration
	PassCount int
	FailCount int
	WarnCount int
	SkipCount int
}

// Doctor performs system diagnostics
type Doctor struct {
	config *config.Config
	client *api.Client
}

// New creates a new Doctor instance
func New(cfg *config.Config, client *api.Client) *Doctor {
	return &Doctor{
		config: cfg,
		client: client,
	}
}

// Run runs all diagnostic checks
func (d *Doctor) Run(ctx context.Context) (*Report, error) {
	report := &Report{
		StartTime: time.Now(),
		Checks:    make([]Check, 0),
	}

	// Run checks
	checks := []func(context.Context) Check{
		d.checkEnvironment,
		d.checkConfig,
		d.checkConnectivity,
		d.checkAuthentication,
		d.checkAPIAccess,
		d.checkDiskSpace,
		d.checkPermissions,
	}

	for _, checkFn := range checks {
		check := checkFn(ctx)
		report.Checks = append(report.Checks, check)

		switch check.Status {
		case StatusPass:
			report.PassCount++
		case StatusFail:
			report.FailCount++
		case StatusWarning:
			report.WarnCount++
		case StatusSkipped:
			report.SkipCount++
		}
	}

	report.Duration = time.Since(report.StartTime)
	return report, nil
}

// checkEnvironment checks system environment
func (d *Doctor) checkEnvironment(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "Environment",
		Description: "System environment and runtime",
		Status:      StatusPass,
	}

	info := fmt.Sprintf("OS: %s, Arch: %s, Go: %s",
		runtime.GOOS, runtime.GOARCH, runtime.Version())

	check.Message = info
	check.Duration = time.Since(start)
	return check
}

// checkConfig checks configuration
func (d *Doctor) checkConfig(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "Configuration",
		Description: "Configuration file and settings",
		Status:      StatusPass,
	}

	if d.config == nil {
		check.Status = StatusFail
		check.Message = "No configuration found"
		check.Duration = time.Since(start)
		return check
	}

	// Check if default instance is configured
	if d.config.DefaultInstance == "" {
		check.Status = StatusFail
		check.Message = "No default instance configured"
		check.Duration = time.Since(start)
		return check
	}

	// Get default instance
	instance, err := d.config.GetDefaultInstance()
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Failed to get default instance: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	if instance.URL == "" {
		check.Status = StatusFail
		check.Message = "Instance URL not configured"
		check.Duration = time.Since(start)
		return check
	}

	check.Message = fmt.Sprintf("Instance: %s, URL: %s", instance.Name, instance.URL)
	check.Duration = time.Since(start)
	return check
}

// checkConnectivity checks network connectivity
func (d *Doctor) checkConnectivity(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "Connectivity",
		Description: "Network connectivity to Canvas",
		Status:      StatusPass,
	}

	if d.config == nil {
		check.Status = StatusSkipped
		check.Message = "Configuration not available"
		check.Duration = time.Since(start)
		return check
	}

	// Get default instance
	instance, err := d.config.GetDefaultInstance()
	if err != nil {
		check.Status = StatusSkipped
		check.Message = "No instance configured"
		check.Duration = time.Since(start)
		return check
	}

	// Test connectivity
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", instance.URL, nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Failed to create request: %v", err)
		check.Error = err
		check.Duration = time.Since(start)
		return check
	}

	resp, err := client.Do(req)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Connection failed: %v", err)
		check.Error = err
		check.Duration = time.Since(start)
		return check
	}
	defer resp.Body.Close()

	check.Message = fmt.Sprintf("Connected successfully (status: %d)", resp.StatusCode)
	check.Duration = time.Since(start)
	return check
}

// checkAuthentication checks API authentication
func (d *Doctor) checkAuthentication(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "Authentication",
		Description: "API token authentication",
		Status:      StatusPass,
	}

	if d.client == nil {
		check.Status = StatusSkipped
		check.Message = "API client not available"
		check.Duration = time.Since(start)
		return check
	}

	// Create users service
	usersService := api.NewUsersService(d.client)

	// Test authentication by fetching user profile
	user, err := usersService.GetCurrentUser(ctx)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Authentication failed: %v", err)
		check.Error = err
		check.Duration = time.Since(start)
		return check
	}

	check.Message = fmt.Sprintf("Authenticated as: %s", user.Name)
	check.Duration = time.Since(start)
	return check
}

// checkAPIAccess checks API accessibility
func (d *Doctor) checkAPIAccess(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "API Access",
		Description: "Canvas API accessibility",
		Status:      StatusPass,
	}

	if d.client == nil {
		check.Status = StatusSkipped
		check.Message = "API client not available"
		check.Duration = time.Since(start)
		return check
	}

	// Create courses service
	coursesService := api.NewCoursesService(d.client)

	// Test API by listing courses
	_, err := coursesService.List(ctx, nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("API access failed: %v", err)
		check.Error = err
		check.Duration = time.Since(start)
		return check
	}

	check.Message = "API accessible"
	check.Duration = time.Since(start)
	return check
}

// checkDiskSpace checks available disk space
func (d *Doctor) checkDiskSpace(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "Disk Space",
		Description: "Available disk space for cache",
		Status:      StatusPass,
	}

	// Check if directory is writable
	configDir, err := config.GetConfigDir()
	if err != nil {
		check.Status = StatusWarning
		check.Message = fmt.Sprintf("Could not determine config directory: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	cacheDir := configDir + "/cache"
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		check.Status = StatusWarning
		check.Message = fmt.Sprintf("Cache directory not writable: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	check.Message = fmt.Sprintf("Cache directory: %s", cacheDir)
	check.Duration = time.Since(start)
	return check
}

// checkPermissions checks file permissions
func (d *Doctor) checkPermissions(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:        "Permissions",
		Description: "File and directory permissions",
		Status:      StatusPass,
	}

	// Check config directory
	configDir, err := config.GetConfigDir()
	if err != nil {
		check.Status = StatusWarning
		check.Message = fmt.Sprintf("Could not determine config directory: %v", err)
		check.Duration = time.Since(start)
		return check
	}
	if info, err := os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			check.Status = StatusWarning
			check.Message = "Config directory does not exist"
			check.Duration = time.Since(start)
			return check
		}
		check.Status = StatusWarning
		check.Message = fmt.Sprintf("Could not check permissions: %v", err)
		check.Duration = time.Since(start)
		return check
	} else {
		// Check permissions
		mode := info.Mode().Perm()
		if mode&0077 != 0 {
			check.Status = StatusWarning
			check.Message = "Config directory has insecure permissions (should be 0700)"
		} else {
			check.Message = "Permissions are secure"
		}
	}

	check.Duration = time.Since(start)
	return check
}

// String returns a string representation of the check status
func (s CheckStatus) String() string {
	return string(s)
}

// IsHealthy returns true if all checks passed
func (r *Report) IsHealthy() bool {
	return r.FailCount == 0
}

// Summary returns a summary of the report
func (r *Report) Summary() string {
	total := len(r.Checks)
	return fmt.Sprintf("Total: %d, Pass: %d, Fail: %d, Warn: %d, Skip: %d",
		total, r.PassCount, r.FailCount, r.WarnCount, r.SkipCount)
}
