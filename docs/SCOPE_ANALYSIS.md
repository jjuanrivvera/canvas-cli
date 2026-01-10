# Canvas CLI - Resource Scope Analysis (Revised)

**Created:** 2026-01-09
**Revised:** 2026-01-09
**Purpose:** Deep analysis of Canvas resource scoping for CLI design

---

## Executive Summary

Canvas resources exist in a **hierarchical context system**, not a simple binary scope:

```
Root Account (ID: 1)
├── Sub-Account A (ID: 5)
│   ├── Sub-Sub-Account A1 (ID: 12)
│   │   └── Course 101
│   └── Course 102
├── Sub-Account B (ID: 8)
│   └── Course 201
└── Course 001 (root level)
```

**Key insight**: An admin at Sub-Account A sees courses 101 and 102, but NOT 201 or 001. The question "show me all courses" is meaningless without knowing *which account*.

**Recommendation**: Use **context flags** (`--account`, `--course`, `--user`) that change which API endpoint is called, rather than an ambiguous `--scope` flag.

---

## The Problem with `--scope`

My initial proposal suggested:
```bash
canvas courses list --scope account
```

**This is flawed because:**

1. **Which account?** User may have admin access to multiple accounts
2. **Implicit magic**: Would need to auto-detect account, which is error-prone
3. **Not generalizable**: Doesn't extend to course-scoped or user-scoped resources
4. **Ambiguous**: `--scope user` vs `--scope account` doesn't capture sub-accounts

---

## How Mature CLIs Solve This

### kubectl (Kubernetes)
```bash
kubectl get pods                      # Current namespace from context
kubectl get pods -n kube-system       # Explicit namespace flag
kubectl get pods --all-namespaces     # ALL namespaces
```
**Pattern**: Explicit `-n` flag OR context-based default

### gh (GitHub CLI)
```bash
gh repo list                          # My repos
gh repo list myorg                    # Org's repos (positional argument)
```
**Pattern**: Positional argument changes context

### aws CLI
```bash
aws s3 ls                             # Current profile
aws --profile admin s3 ls             # Different account via profile
```
**Pattern**: Profile system for account switching

### gcloud
```bash
gcloud compute instances list --project=myproj
gcloud config set project myproj      # Set default
```
**Pattern**: Explicit flag OR config default

---

## Canvas API Structure

Canvas has three primary context levels:

| Context | Endpoint Pattern | Example |
|---------|-----------------|---------|
| **User** | `/api/v1/{resource}` | `GET /courses` → my enrolled courses |
| **Account** | `/api/v1/accounts/{id}/{resource}` | `GET /accounts/1/courses` → all account courses |
| **Course** | `/api/v1/courses/{id}/{resource}` | `GET /courses/123/users` → course roster |

Resources exist at different context levels:

| Resource | User Context | Course Context | Account Context |
|----------|--------------|----------------|-----------------|
| Courses | ✅ My enrollments | ❌ | ✅ All courses |
| Users | ✅ Self only | ✅ Course roster | ✅ Account directory |
| Files | ✅ My files | ✅ Course files | ❌ |
| Groups | ✅ My groups | ✅ Course groups | ✅ Account groups |
| Enrollments | ✅ My enrollments | ✅ Course enrollments | ❌ |
| Assignments | ❌ | ✅ Course assignments | ❌ |
| Rubrics | ❌ | ✅ Course rubrics | ✅ Account rubrics |

---

## Revised Recommendation: Context Flags

### Core Design

**Global context flags** that change the API endpoint:

```bash
--account <id>    # Use /accounts/<id>/... endpoints
--course <id>     # Use /courses/<id>/... endpoints
--user <id>       # Use /users/<id>/... endpoints
--as-user <id>    # Masquerade (adds ?as_user_id=<id> param)
```

### Usage Examples

**Courses:**
```bash
# User context (default)
canvas courses list                    # GET /courses
                                       # Returns: my enrolled courses

# Account context (explicit)
canvas courses list --account 1        # GET /accounts/1/courses
                                       # Returns: all courses in account 1

# Account context with search
canvas courses list --account 1 --search "Biology" --state available
```

**Users:**
```bash
# Self (default)
canvas users me                        # GET /users/self

# Course context
canvas users list --course 123         # GET /courses/123/users
                                       # Returns: course roster

# Account context
canvas users list --account 1          # GET /accounts/1/users
                                       # Returns: account directory
```

**Files:**
```bash
# User context (default)
canvas files list                      # GET /users/self/files

# Course context
canvas files list --course 123         # GET /courses/123/files

# User context (other user, requires admin)
canvas files list --user 456           # GET /users/456/files
```

**Groups:**
```bash
# User context (default)
canvas groups list                     # GET /users/self/groups

# Course context
canvas groups list --course 123        # GET /courses/123/groups

# Account context
canvas groups list --account 1         # GET /accounts/1/groups
```

### Default Account Shorthand

For convenience, allow `--account` without a value to use config default:

```bash
# One-time setup
canvas accounts list                   # Discover available accounts
canvas config set default-account 1    # Set default

# Then use shorthand
canvas courses list --account          # Uses account 1 from config
canvas courses list --account 5        # Override with explicit ID
```

### Error Handling

When `--account` is used without value and no default is configured:

```
Error: No account specified and no default configured.

Available accounts (from `canvas accounts list`):
  ID    NAME                   ROLE
  1     Example University     Account Admin
  5     Biology Department     Sub-Account Admin

Set a default with: canvas config set default-account <id>
Or specify explicitly: canvas courses list --account <id>
```

---

## Implementation Architecture

### Flag Processing Flow

```go
func runCoursesList(cmd *cobra.Command, args []string) error {
    accountID, _ := cmd.Flags().GetInt64("account")
    accountFlag := cmd.Flags().Changed("account")

    if accountFlag {
        if accountID == 0 {
            // --account without value: use config default
            accountID = config.GetDefaultAccountID()
            if accountID == 0 {
                return errors.New("no default account configured")
            }
        }
        // Account context: GET /accounts/{id}/courses
        return listAccountCourses(ctx, accountID, opts)
    }

    // User context (default): GET /courses
    return listUserCourses(ctx, opts)
}
```

### Service Layer Design

```go
// internal/api/courses.go

// User context - existing method
func (s *CoursesService) List(ctx context.Context, opts *ListCoursesOptions) ([]Course, error) {
    // GET /courses
}

// Account context - new method
func (s *CoursesService) ListByAccount(ctx context.Context, accountID int64, opts *ListAccountCoursesOptions) ([]Course, error) {
    endpoint := fmt.Sprintf("/accounts/%d/courses", accountID)
    // GET /accounts/{id}/courses
}
```

### Options Structs

```go
// User-context options (existing)
type ListCoursesOptions struct {
    EnrollmentType  string   `url:"enrollment_type,omitempty"`
    EnrollmentState string   `url:"enrollment_state,omitempty"`
    Include         []string `url:"include[],omitempty"`
    State           []string `url:"state[],omitempty"`
    PerPage         int      `url:"per_page,omitempty"`
}

// Account-context options (new - different params available!)
type ListAccountCoursesOptions struct {
    // Shared
    Include []string `url:"include[],omitempty"`
    State   []string `url:"state[],omitempty"`
    PerPage int      `url:"per_page,omitempty"`

    // Account-specific (not available on user endpoint)
    SearchTerm    string  `url:"search_term,omitempty"`
    ByTeachers    []int64 `url:"by_teachers[],omitempty"`
    BySubaccounts []int64 `url:"by_subaccounts[],omitempty"`
    Sort          string  `url:"sort,omitempty"`      // course_name, sis_course_id, teacher, account_name
    Order         string  `url:"order,omitempty"`     // asc, desc

    // These are NOT available at account level
    // EnrollmentType  - only on user endpoint
    // EnrollmentState - only on user endpoint
}
```

---

## User Personas and Workflows

### Persona 1: New Administrator

**Journey:**
1. Gets admin access to Canvas
2. Installs CLI, runs `canvas auth login`
3. Runs `canvas courses list` - sees only 2 enrolled courses
4. Confused: "Where are the other 500 courses?"

**Solution:**
```bash
$ canvas courses list
Showing your enrolled courses. For all account courses, use --account <id>
Tip: Run `canvas accounts list` to see available accounts.

ID      NAME                    CODE      STATE
12345   Intro to Biology        BIO101    available
12346   Advanced Chemistry      CHEM301   available

$ canvas accounts list
ID    NAME                   ROLE
1     State University       Account Admin

$ canvas courses list --account 1
ID      NAME                    CODE      STATE       TERM
...     (500 courses)
```

### Persona 2: Multi-Account Administrator

**Journey:**
1. Manages Biology Dept (ID:5) and Chemistry Dept (ID:8)
2. Needs to switch between contexts frequently

**Solution:**
```bash
# Explicit context - always works
$ canvas courses list --account 5    # Biology courses
$ canvas courses list --account 8    # Chemistry courses

# Or set default for session
$ canvas config set default-account 5
$ canvas courses list --account       # Now shows Biology
```

### Persona 3: Developer/Scripter

**Journey:**
1. Writing automation scripts
2. Needs deterministic, explicit commands

**Solution:**
```bash
#!/bin/bash
# Script is explicit - no ambiguity
ACCOUNT_ID=1

canvas courses list --account $ACCOUNT_ID --state available --format json | \
  jq '.[] | select(.enrollment_count > 100)'
```

---

## Mutual Exclusivity Rules

Context flags are mutually exclusive where it doesn't make sense:

| Command | `--account` | `--course` | `--user` | Notes |
|---------|-------------|------------|----------|-------|
| `courses list` | ✅ | ❌ | ❌ | Courses belong to accounts |
| `users list` | ✅ | ✅ | ❌ | Users can be in accounts OR courses |
| `files list` | ❌ | ✅ | ✅ | Files belong to courses OR users |
| `groups list` | ✅ | ✅ | ❌ | Groups can be in accounts OR courses |
| `enrollments list` | ❌ | ✅ | ✅ | Enrollments are per-course OR per-user |

When mutually exclusive flags are both provided:
```
Error: --account and --course cannot be used together for this command.
Use --account for account-level groups, or --course for course-level groups.
```

---

## Masquerading (`--as-user`)

Different from context flags - adds query parameter, doesn't change endpoint:

```bash
# See what student 12345 sees
canvas courses list --as-user 12345

# API call: GET /courses?as_user_id=12345
# Returns: Courses as if student 12345 made the request
```

**Important distinctions:**
- `--user 12345` → Changes endpoint to `/users/12345/...`
- `--as-user 12345` → Adds `?as_user_id=12345` to ANY endpoint

Masquerading requires "Become other users" permission.

---

## Implementation Checklist

### Phase 1: Foundation
- [ ] Add `Account` type to `internal/api/types.go`
- [ ] Create `internal/api/accounts.go` with `AccountsService`
- [ ] Add `canvas accounts list` command
- [ ] Add `canvas accounts get <id>` command
- [ ] Add `default-account` to config

### Phase 2: Courses Context
- [ ] Add `--account` flag to `courses list`
- [ ] Create `ListAccountCoursesOptions` struct
- [ ] Implement `CoursesService.ListByAccount()`
- [ ] Add account-specific flags: `--search`, `--by-teacher`, `--sort`
- [ ] Update help text with context examples

### Phase 3: Users Context
- [ ] Add `--account` flag to `users list`
- [ ] Add `--course` flag to `users list`
- [ ] Implement `UsersService.ListByCourse()`
- [ ] Implement `UsersService.ListByAccount()`

### Phase 4: Other Resources
- [ ] Files: Add `--course` and `--user` flags
- [ ] Groups: Add `--account` and `--course` flags
- [ ] Enrollments: Add `--course` and `--user` flags

### Phase 5: Masquerading
- [ ] Add global `--as-user` flag
- [ ] Modify API client to append `as_user_id` param
- [ ] Add permission checking/warnings
- [ ] Document audit trail implications

### Phase 6: UX Polish
- [ ] Helpful error messages when context is ambiguous
- [ ] Suggestions in output ("Tip: use --account for all courses")
- [ ] Shell completion for account/course IDs
- [ ] Config shortcuts (`canvas config set default-account`)

---

## Comparison: Original vs Revised Approach

| Aspect | Original (`--scope`) | Revised (Context Flags) |
|--------|---------------------|------------------------|
| Clarity | Ambiguous: which account? | Explicit: `--account 1` |
| Generalizability | Only user/account | Any context level |
| API alignment | Abstracts away endpoints | Mirrors API structure |
| Discoverability | Hidden magic | `canvas accounts list` |
| Scripting | Unpredictable | Deterministic |
| Implementation | Complex scope logic | Simple flag → endpoint |

---

## References

- [Canvas Courses API](https://canvas.instructure.com/doc/api/courses.html)
- [Canvas Accounts API](https://canvas.instructure.com/doc/api/accounts.html)
- [Canvas Masquerading](https://canvas.instructure.com/doc/api/file.masquerading.html)
- [kubectl Namespace Pattern](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
- [gh CLI Patterns](https://cli.github.com/manual/)
