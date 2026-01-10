# Canvas CLI - Improvement Items

**Created:** 2026-01-09
**Status:** Day 1 Development Review
**Source:** Automated code review by expert agents

---

## Summary

This document contains all improvement items identified during the initial code review. Items are categorized by priority and area. Each item includes the file location, description, and recommended fix.

**Total Items:** 23 code fixes + 1 major feature
- Critical (Security): 3
- High Priority: 5
- Medium Priority: 8
- Low Priority (Nice to Have): 7
- New Feature: Resource Context System (24 sub-items)

---

## Critical Priority (Security)

### 1. OAuth State Parameter Hardcoded (CSRF Vulnerability)

**File:** `internal/auth/oauth.go:169-172`

**Current Code:**
```go
authURL := f.oauth2Config.AuthCodeURL(
    "state",  // Hardcoded!
    oauth2.SetAuthURLParam("code_challenge", f.pkce.Challenge),
    oauth2.SetAuthURLParam("code_challenge_method", f.pkce.Method),
)
```

**Issue:** Using a hardcoded state parameter enables CSRF attacks. An attacker could forge authorization requests.

**Fix:** Generate a cryptographically secure random state value and validate it in the callback.

```go
state, err := generateSecureState()
if err != nil {
    return nil, fmt.Errorf("failed to generate state: %w", err)
}

authURL := f.oauth2Config.AuthCodeURL(
    state,
    oauth2.SetAuthURLParam("code_challenge", f.pkce.Challenge),
    oauth2.SetAuthURLParam("code_challenge_method", f.pkce.Method),
)
```

---

### 2. Weak Encryption Key Derivation

**File:** `internal/auth/encryption.go:16-33`

**Current Code:**
```go
func deriveEncryptionKey() ([]byte, error) {
    machineID, err := getMachineID()
    // ...
    combined := machineID + ":" + username
    hash := sha256.Sum256([]byte(combined))  // Not a proper KDF
    return hash[:], nil
}
```

**Issue:** SHA-256 is not a proper Key Derivation Function. No salt, no iterations, vulnerable to rainbow table attacks.

**Fix:** Use PBKDF2, scrypt, or Argon2 with:
- Random salt stored with the encrypted file
- High iteration count (100,000+ for PBKDF2)
- Proper key stretching

```go
import "golang.org/x/crypto/pbkdf2"

func deriveEncryptionKey(salt []byte) ([]byte, error) {
    machineID, err := getMachineID()
    if err != nil {
        return nil, err
    }
    username := getUsername()
    combined := machineID + ":" + username

    // PBKDF2 with 100,000 iterations
    key := pbkdf2.Key([]byte(combined), salt, 100000, 32, sha256.New)
    return key, nil
}
```

---

### 3. Weak Machine ID Fallback

**File:** `internal/auth/encryption.go:70-76`

**Current Code:**
```go
// Fallback: use hostname
hostname, err := os.Hostname()
if err != nil {
    return "", fmt.Errorf("failed to get hostname: %w", err)
}
return hostname, nil
```

**Issue:** Hostnames are often guessable (e.g., "MacBook-Pro"), making the encryption key predictable.

**Fix:** Either fail closed if proper machine ID is unavailable, or use a randomly generated salt stored alongside the encrypted file.

---

## High Priority

### 4. openBrowser() Function Not Implemented

**File:** `internal/auth/oauth.go:268-273`

**Current Code:**
```go
func openBrowser(url string) {
    // This is a best-effort attempt - we don't error if it fails
    // The user can always use the printed URL
    // Implementation varies by OS...
}
```

**Issue:** Function body is empty. Browser never opens automatically.

**Fix:** Implement cross-platform browser opening:

```go
import "github.com/pkg/browser"

func openBrowser(url string) {
    _ = browser.OpenURL(url)  // Best effort, ignore errors
}
```

Or manual implementation:
```go
func openBrowser(url string) {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "darwin":
        cmd = exec.Command("open", url)
    case "linux":
        cmd = exec.Command("xdg-open", url)
    case "windows":
        cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
    }
    if cmd != nil {
        _ = cmd.Start()
    }
}
```

---

### 5. Sync Conflict Detection Uses Wrong Field

**File:** `internal/batch/sync.go:96`

**Current Code:**
```go
if s.interactive && sourceCourse.CreatedAt != targetCourse.CreatedAt {
```

**Issue:** Comparing `CreatedAt` doesn't detect modifications. Should compare `UpdatedAt`.

**Fix:**
```go
if s.interactive && sourceCourse.UpdatedAt != targetCourse.UpdatedAt {
```

---

### 6. createAssignmentInTarget Not Implemented

**File:** `internal/batch/sync.go:72-76`

**Current Code:**
```go
func (s *SyncOperation) createAssignmentInTarget(ctx context.Context, courseID int64, assignment *api.Assignment) error {
    return fmt.Errorf("assignment creation in target not yet implemented - requires API client update")
}
```

**Issue:** Core sync functionality is incomplete. The entire sync feature is unusable.

**Fix:** Implement the function using AssignmentsService.Create() or mark the sync feature as experimental/incomplete in documentation.

---

### 7. Error Helper Functions Don't Handle Wrapped Errors

**File:** `internal/api/errors.go:48-69`

**Current Code:**
```go
func IsRateLimitError(err error) bool {
    if apiErr, ok := err.(*APIError); ok {  // Type assertion
        return apiErr.StatusCode == http.StatusTooManyRequests
    }
    return false
}
```

**Issue:** Type assertion only works for exact types. Won't work if error is wrapped with `fmt.Errorf("context: %w", err)`.

**Fix:** Use `errors.As()` for proper wrapped error handling:

```go
func IsRateLimitError(err error) bool {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        return apiErr.StatusCode == http.StatusTooManyRequests
    }
    return false
}
```

Apply same fix to `IsAuthError()` and `IsNotFoundError()`.

---

### 8. Retry Policy Error Comparison Not Idiomatic

**File:** `internal/api/retry.go:40-41`

**Current Code:**
```go
if err == context.Canceled || err == context.DeadlineExceeded {
    return false
}
```

**Issue:** Uses `==` for error comparison. While it works for sentinel errors, `errors.Is()` is more idiomatic and handles wrapped errors.

**Fix:**
```go
if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
    return false
}
```

---

## Medium Priority

### 9. Missing HTTP Transport Configuration

**File:** `internal/api/client.go:134-136`

**Current Code:**
```go
client := &Client{
    httpClient: &http.Client{
        Timeout: config.Timeout,
    },
    // ...
}
```

**Issue:** Uses default transport with no control over connection pooling.

**Fix:** Add explicit transport configuration:

```go
transport := &http.Transport{
    MaxIdleConns:          10,
    MaxIdleConnsPerHost:   5,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ResponseHeaderTimeout: 10 * time.Second,
}

client := &Client{
    httpClient: &http.Client{
        Timeout:   config.Timeout,
        Transport: transport,
    },
    // ...
}
```

---

### 10. NewClient Should Accept Context Parameter

**File:** `internal/api/client.go:112`

**Current Code:**
```go
func NewClient(config ClientConfig) (*Client, error) {
```

**Issue:** Version detection uses `context.Background()` instead of accepting a parent context.

**Fix:**
```go
func NewClient(ctx context.Context, config ClientConfig) (*Client, error) {
    // ...
    versionCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    version, err := DetectCanvasVersion(versionCtx, client.httpClient, config.BaseURL)
    // ...
}
```

---

### 11. Cache Missing Close() Method

**File:** `internal/cache/cache.go:23-33`

**Current Code:**
```go
func New(ttl time.Duration) *Cache {
    c := &Cache{
        items: make(map[string]*item),
        ttl:   ttl,
    }
    go c.cleanup()  // Goroutine starts but never stops
    return c
}
```

**Issue:** No way to stop the cleanup goroutine. Goroutine leak in long-running applications.

**Fix:** Add stop channel and Close() method:

```go
type Cache struct {
    items  map[string]*item
    mu     sync.RWMutex
    ttl    time.Duration
    stopCh chan struct{}
}

func New(ttl time.Duration) *Cache {
    c := &Cache{
        items:  make(map[string]*item),
        ttl:    ttl,
        stopCh: make(chan struct{}),
    }
    go c.cleanup()
    return c
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            c.removeExpired()
        case <-c.stopCh:
            return
        }
    }
}

func (c *Cache) Close() error {
    close(c.stopCh)
    return nil
}
```

---

### 12. Hardcoded Canvas Quota

**File:** `internal/api/client.go:222`

**Current Code:**
```go
c.rateLimiter.AdjustRate(remainingFloat, 700.0)  // Hardcoded
```

**Issue:** Canvas quota of 700 is hardcoded. Different instances may have different quotas.

**Fix:** Make configurable or detect from headers:

```go
// Try to get quota total from X-Rate-Limit-Total header
if total := resp.Header.Get("X-Rate-Limit-Total"); total != "" {
    if totalFloat, err := strconv.ParseFloat(total, 64); err == nil {
        c.quotaTotal = totalFloat
    }
}
c.rateLimiter.AdjustRate(remainingFloat, c.quotaTotal)
```

---

### 13. Pagination URL Handling

**File:** `internal/api/client.go:322-326`

**Current Code:**
```go
nextURL, err := url.Parse(links.Next)
if err != nil {
    return fmt.Errorf("failed to parse next URL: %w", err)
}
currentURL = nextURL.Path + "?" + nextURL.RawQuery
```

**Issue:** If `RawQuery` is empty, results in trailing `?`.

**Fix:**
```go
if nextURL.RawQuery != "" {
    currentURL = nextURL.Path + "?" + nextURL.RawQuery
} else {
    currentURL = nextURL.Path
}
```

---

### 14. CSV File Permissions

**File:** `internal/batch/csv.go:97`

**Current Code:**
```go
file, err := os.Create(filename)  // Uses default permissions
```

**Issue:** Files created with default permissions (0666 umask). Should use secure permissions.

**Fix:**
```go
file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
```

---

### 15. Interactive Prompts Without Timeout

**File:** `internal/batch/sync.go:139, 159`

**Current Code:**
```go
var choice string
fmt.Scanln(&choice)  // Blocks indefinitely
```

**Issue:** `fmt.Scanln` blocks indefinitely. Could cause hangs in non-TTY environments.

**Fix:** Add context-aware input with timeout, or validate TTY availability before enabling interactive mode.

---

### 16. SyncAssignments Not Using Batch Processor

**File:** `internal/batch/sync.go:178-209`

**Current Code:**
```go
for _, assignment := range assignments {
    err := s.CopyAssignment(ctx, sourceCourseID, targetCourseID, assignment.ID)
    // ...
}
```

**Issue:** Processes assignments sequentially instead of using the concurrent batch processor.

**Fix:** Use the Processor for concurrent syncing:

```go
processor := New(5, false, nil)
items := make([]interface{}, len(assignments))
for i, a := range assignments {
    items[i] = a
}
summary, err := processor.Process(ctx, items, func(ctx context.Context, item interface{}) error {
    assignment := item.(*api.Assignment)
    return s.CopyAssignment(ctx, sourceCourseID, targetCourseID, assignment.ID)
})
```

---

## Low Priority (Nice to Have)

### 17. Define API Client Interface

**File:** `internal/api/client.go`

**Issue:** `Client` is a concrete type with no interface. Makes mocking difficult for testing.

**Fix:** Define interface for better testability:

```go
type HTTPClient interface {
    Get(ctx context.Context, path string) (*http.Response, error)
    Post(ctx context.Context, path string, body io.Reader) (*http.Response, error)
    Put(ctx context.Context, path string, body io.Reader) (*http.Response, error)
    Delete(ctx context.Context, path string) (*http.Response, error)
    GetJSON(ctx context.Context, path string, result interface{}) error
    PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error
    PutJSON(ctx context.Context, path string, body interface{}, result interface{}) error
    GetAllPages(ctx context.Context, path string, result interface{}) error
}

var _ HTTPClient = (*Client)(nil)
```

---

### 18. Split Large Types File

**File:** `internal/api/types.go` (431 lines)

**Issue:** All domain types in one file.

**Fix:** Split by domain:
```
internal/api/
├── types_course.go
├── types_assignment.go
├── types_user.go
├── types_submission.go
```

---

### 19. Extract Display Logic from Commands

**File:** `commands/courses.go:107-119` (and similar in other command files)

**Issue:** Display/formatting logic embedded in command handlers.

**Fix:** Extract to `internal/output/` package formatters:

```go
type CourseFormatter interface {
    FormatList(courses []Course) string
    FormatDetail(course *Course) string
}
```

---

### 20. Eliminate Duplication in Courses Create/Update

**File:** `internal/api/courses.go:128-232` and `264-359`

**Issue:** Near-identical field mapping logic in Create and Update methods.

**Fix:** Extract shared logic:

```go
func mapCourseParams(params interface{}) map[string]interface{} {
    // Shared field mapping
}
```

Or use Builder pattern for complex parameter structs.

---

### 21. Add Benchmark Tests

**Files:** Various `*_test.go` files

**Issue:** No `Benchmark*` functions found. Performance characteristics not validated.

**Fix:** Add benchmarks for critical paths:

```go
func BenchmarkCache_Get(b *testing.B) {
    c := New(time.Hour)
    c.Set("key", []byte("value"))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        c.Get("key")
    }
}

func BenchmarkProcessor_Process(b *testing.B) {
    // ...
}
```

---

### 22. Improve Commands Test Coverage

**File:** `commands/` directory (12.9% coverage)

**Issue:** CLI command handlers have critically low test coverage.

**Fix:** Add integration tests for CLI commands:
- Test flag parsing and validation
- Mock file system operations
- Test output formatting
- Test error handling paths

---

### 23. Implement REPL Syntax Highlighter

**File:** `internal/repl/highlighter.go` (missing)

**Issue:** Syntax highlighting mentioned in spec but not implemented.

**Fix:** Create highlighter using chroma library:

```go
package repl

import "github.com/alecthomas/chroma/quick"

type Highlighter struct {
    // ...
}

func (h *Highlighter) Highlight(input string) string {
    // Highlight commands in blue, flags in green, values in yellow
}
```

---

## Implementation Checklist

### Phase 1: Security Fixes
- [x] 1. Generate random OAuth state parameter ✅
- [x] 2. Implement proper KDF with salt (PBKDF2 with 100k iterations) ✅
- [x] 3. Improve machine ID fallback (fail closed, no hostname) ✅

### Phase 2: Core Functionality
- [x] 4. Implement openBrowser() ✅
- [x] 5. Fix sync conflict detection field (UpdatedAt) ✅
- [x] 6. Implement createAssignmentInTarget ✅
- [x] 7. Use errors.As() in error helpers ✅
- [x] 8. Use errors.Is() in retry policy ✅

### Phase 3: Code Quality
- [x] 9. Add HTTP Transport configuration ✅
- [ ] 10. Accept context in NewClient (deferred - breaking change)
- [x] 11. Add Cache Close() method ✅
- [x] 12. Make Canvas quota configurable ✅
- [x] 13. Fix pagination URL handling ✅
- [x] 14. Set secure CSV file permissions ✅
- [x] 15. Add timeout to interactive prompts ✅
- [x] 16. Use batch processor in SyncAssignments ✅

### Phase 4: Architecture Improvements
- [x] 17. Define API Client interface (HTTPClient) ✅
- [ ] 18. Split types.go file (deferred - risk vs benefit)
- [ ] 19. Extract display logic from commands (deferred - low priority)
- [x] 20. Eliminate Create/Update duplication ✅

### Phase 5: Testing & Polish
- [x] 21. Add benchmark tests (cache, batch) ✅
- [ ] 22. Improve commands test coverage (ongoing)
- [x] 23. Implement REPL syntax highlighter ✅

---

## New Feature: Resource Context System

**Reference:** `docs/SCOPE_ANALYSIS.md`

### Problem Statement

Canvas resources exist at different scope levels. Currently, `canvas courses list` returns only the user's enrolled courses. Administrators need to list ALL account courses, but the question "show me all courses" is meaningless without knowing *which account*.

```
Root Account (ID: 1)
├── Sub-Account A (ID: 5)
│   ├── Course 101
│   └── Course 102
├── Sub-Account B (ID: 8)
│   └── Course 201
└── Course 001 (root level)
```

An admin at Sub-Account A sees courses 101 and 102, but NOT 201 or 001.

### Solution: Context Flags

Instead of an ambiguous `--scope` flag, use **explicit context flags** that change which API endpoint is called:

```bash
--account <id>    # Use /accounts/<id>/... endpoints
--course <id>     # Use /courses/<id>/... endpoints
--user <id>       # Use /users/<id>/... endpoints
--as-user <id>    # Masquerade (adds ?as_user_id=<id> param)
```

### Usage Examples

```bash
# Courses
canvas courses list                    # GET /courses (my courses)
canvas courses list --account 1        # GET /accounts/1/courses (all account courses)
canvas courses list --account 1 --search "Biology"

# Users
canvas users list --course 123         # GET /courses/123/users (course roster)
canvas users list --account 1          # GET /accounts/1/users (account directory)

# Files
canvas files list                      # GET /users/self/files (my files)
canvas files list --course 123         # GET /courses/123/files

# Groups
canvas groups list --account 1         # GET /accounts/1/groups
canvas groups list --course 123        # GET /courses/123/groups

# Masquerading (different - adds query param)
canvas courses list --as-user 456      # GET /courses?as_user_id=456
```

### Default Account Shorthand

```bash
# One-time setup
canvas accounts list                   # Discover available accounts
canvas config set default-account 1    # Set default

# Then use shorthand
canvas courses list --account          # Uses account 1 from config
canvas courses list --account 5        # Override with explicit ID
```

### Context Flag Matrix

| Command | `--account` | `--course` | `--user` | Notes |
|---------|-------------|------------|----------|-------|
| `courses list` | ✅ | ❌ | ❌ | Courses belong to accounts |
| `users list` | ✅ | ✅ | ❌ | Users in accounts OR courses |
| `files list` | ❌ | ✅ | ✅ | Files belong to courses OR users |
| `groups list` | ✅ | ✅ | ❌ | Groups in accounts OR courses |
| `enrollments list` | ❌ | ✅ | ✅ | Enrollments per-course OR per-user |

### Implementation Checklist

#### Phase A: Foundation
- [x] A1. Add `Account` type to `internal/api/types.go` ✅
- [x] A2. Create `internal/api/accounts.go` with `AccountsService` ✅
- [x] A3. Add `canvas accounts list` command ✅
- [x] A4. Add `canvas accounts get <id>` command ✅
- [ ] A5. Add `default-account` to config (deferred - low priority)
- [ ] A6. Add `canvas config set default-account` command (deferred - low priority)

#### Phase B: Courses Context
- [x] B1. Add `--account` flag to `courses list` ✅
- [x] B2. Create `ListAccountCoursesOptions` struct ✅
- [x] B3. Implement `AccountsService.ListCourses()` ✅
- [x] B4. Add `--search` flag (account-level only) ✅
- [ ] B5. Add `--by-teacher` flag (account-level only) (deferred)
- [x] B6. Add `--sort` and `--order` flags ✅
- [x] B7. Update help text with context examples ✅

#### Phase C: Users Context
- [x] C1. Add `--account` flag to `users list` ✅ (as --account-id)
- [x] C2. Add `--course` flag to `users list` ✅ (as --course-id)
- [x] C3. Implement `UsersService.ListCourseUsers()` ✅
- [x] C4. Implement `UsersService.List()` ✅

#### Phase D: Other Resources
- [x] D1. Files: Add `--course` and `--user` flags ✅ (as --course-id, --user-id)
- [ ] D2. Groups: Add `--account` and `--course` flags (deferred - no groups API service)
- [x] D3. Enrollments: Add `--course` and `--user` flags ✅ (created enrollments command)

#### Phase E: Masquerading
- [x] E1. Add global `--as-user` flag ✅
- [x] E2. Modify API client to append `as_user_id` param ✅
- [x] E3. Add permission checking/warnings ✅ (verbose warning when active)
- [x] E4. Document audit trail implications ✅ (warning message mentions audit log)

#### Phase F: UX Polish
- [x] F1. Helpful error messages when context is ambiguous ✅
- [x] F2. Suggestions in output ("Tip: use --account for all courses") ✅
- [x] F3. Shell completion (basic) ✅ (dynamic ID completion deferred - requires API calls)
- [x] F4. Mutual exclusivity validation for conflicting flags ✅ (implemented in users, files, enrollments)

---

## Notes

### Code Fixes (Items 1-23)
- Items 1-3 are security-critical and should be addressed before any public release
- Items 4-8 affect core functionality and user experience
- Items 9-16 improve code quality and robustness
- Items 17-23 are architectural improvements and polish

### Resource Context System (New Feature)
- Phase A (Foundation) should be implemented first - provides account discovery
- Phase B (Courses) is the highest-value deliverable for administrators
- Phases C-D extend the pattern to other resources
- Phase E (Masquerading) is optional but valuable for support workflows
- Phase F (UX Polish) can be done incrementally

### Implementation Order Recommendation

1. **Security First** (Items 1-3) - Critical before any release
2. **Context Foundation** (Phase A) - Enables admin workflows
3. **Courses Context** (Phase B) - Highest user value
4. **Core Functionality** (Items 4-8) - Complete existing features
5. **Other Contexts** (Phases C-D) - Extend pattern
6. **Code Quality** (Items 9-23) - Polish and robustness
7. **Masquerading & UX** (Phases E-F) - Advanced features

All code fix items were identified through automated code review and verified against the actual codebase. The Resource Context System was designed based on analysis of Canvas API structure, canvas-lms-kit patterns, and mature CLI conventions (kubectl, gh, aws).
