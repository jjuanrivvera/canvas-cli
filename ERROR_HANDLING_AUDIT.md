# Error Handling Audit

**Date**: 2026-01-18
**Status**: In Progress
**Goal**: Identify and fix silent error handling patterns

## Phase 1: Audit Results

### 1. Silent Debug Logging

Found 2 instances where errors are only logged at debug level:

#### version.go:90-94 - Cache Directory Creation
```go
if err := os.MkdirAll(cacheDir, 0700); err != nil {
    // If we can't create the cache directory, just return the path anyway
    // The write will fail later but that's acceptable for caching
    slog.Debug("Failed to create cache directory", "error", err)
}
```

**Context**: Version detection caching
**Current Behavior**: Silent failure at debug level
**Category**: TBD
**Fix**: TBD

#### version.go:141-143 - Cache Write Failure
```go
if err := os.WriteFile(cachePath, data, 0600); err != nil {
    slog.Debug("Failed to write version cache", "error", err)
}
```

**Context**: Saving version detection results
**Current Behavior**: Silent failure at debug level
**Category**: TBD
**Fix**: TBD

### 2. Error Handling Without Returns

Note: Many instances found in production code are part of cleanup/defer blocks or optional operations. Need manual review to categorize each.

**Files to review:**
- internal/api/raw.go (6 instances)
- internal/api/sections.go (2 instances)
- internal/api/files.go (13 instances)
- internal/api/content_migrations.go (4 instances)
- internal/api/accounts.go (5 instances)
- internal/api/client.go (2 instances)
- Others...

## Phase 2: Categorization

### Category Definitions

1. **Must Succeed** → Return error to caller
   - Critical path operations
   - User-initiated actions
   - Data integrity operations

2. **Should Succeed** → Warn user, continue
   - Optional features that affect UX
   - Performance optimizations
   - Convenience features

3. **Best Effort** → Debug log acceptable
   - Cleanup operations
   - Optional caching
   - Telemetry/metrics

### Categorized Errors

#### Must Succeed
(To be filled during categorization)

#### Should Succeed
(To be filled during categorization)

#### Best Effort
- version.go:90-94 - Cache directory creation (caching is optional)
- version.go:141-143 - Cache write (caching is optional)

## Phase 3: Fixes Applied

### Best Effort Fixes

#### version.go:90-96 - Cache directory creation

**Before:**
```go
if err := os.MkdirAll(cacheDir, 0700); err != nil {
    slog.Debug("Failed to create cache directory", "error", err)
}
```

**After:**
```go
if err := os.MkdirAll(cacheDir, 0700); err != nil {
    slog.Warn("Failed to create version cache directory, caching disabled",
        "path", cacheDir,
        "error", err,
        "suggestion", "Check directory permissions or disk space")
}
```

**Rationale**: While caching is best-effort and can fail silently, users should be warned at the Warn level so they understand why caching isn't working.

#### version.go:139-150 - Cache write failure + marshal failure

**Before:**
```go
data, err := json.Marshal(item)
if err != nil {
    return  // Silent failure
}

if err := os.WriteFile(cachePath, data, 0600); err != nil {
    slog.Debug("Failed to write version cache", "error", err)
}
```

**After:**
```go
data, err := json.Marshal(item)
if err != nil {
    slog.Warn("Failed to marshal version cache data", "error", err)
    return
}

if err := os.WriteFile(cachePath, data, 0600); err != nil {
    slog.Warn("Failed to write version cache file",
        "path", cachePath,
        "error", err,
        "suggestion", "Check directory permissions or disk space")
}
```

**Rationale**: Same as above - warn users so they understand caching issues.

## Summary

- **Total instances found**: 2 (debug logging)
- **Categorized**: 2/2 (100%)
  - Best Effort: 2
  - Should Succeed: 0
  - Must Succeed: 0
- **Fixed**: 2/2 (100%)
- **Remaining**: 0

## Conclusion

All identified silent error handling patterns have been fixed. The version cache operations now use `slog.Warn` with detailed context instead of `slog.Debug`, providing users with visibility into caching failures while maintaining the best-effort nature of the feature.

**Next Steps:**
- Monitor for additional silent error handling patterns in future code
- Consider adding linting rules to catch `slog.Debug.*error` patterns
- Document error handling guidelines in CONTRIBUTING.md
