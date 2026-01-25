# CLI UX Improvements for Canvas CLI

Inspired by modern CLI best practices from gh, stripe-cli, kubectl, and others.

## Research Sources

Based on analysis of:
- [GitHub CLI best features](https://onlyutkarsh.com/posts/2026/github-cli-power-tips/)
- [Stripe CLI event forwarding](https://docs.stripe.com/stripe-cli/use-cli)
- [kubectl context switching](https://home.robusta.dev/blog/switching-kubernets-context)
- [Modern CLI UX patterns](https://clig.dev/)
- [CLI UX best practices](https://lucasfcosta.com/2022/06/01/ux-patterns-cli-tools.html)

---

## Current Canvas CLI Strengths

Canvas CLI already has excellent foundations:
- ✅ Multiple output formats (table, JSON, YAML, CSV)
- ✅ Interactive REPL with history
- ✅ Multi-instance profiles
- ✅ Smart caching and rate limiting
- ✅ Structured logging
- ✅ Dry-run mode (partial)
- ✅ Global flags (--limit, --verbose, --as-user)

---

## Recommended UX Improvements

### 🔥 HIGH IMPACT (Must Haves)

#### 1. **Shell Completion** (HIGHEST PRIORITY)
**Inspired by:** All modern CLIs
**Status:** ❌ Missing

**Implementation:**
```bash
# Generate completion for current shell
canvas completion bash > /etc/bash_completion.d/canvas
canvas completion zsh > /usr/local/share/zsh/site-functions/_canvas
canvas completion fish > ~/.config/fish/completions/canvas.fish
canvas completion powershell > canvas.ps1
```

**Why it matters:**
- Cobra has this built-in (2-3 hours to implement)
- Massive UX improvement for daily usage
- Tab-complete commands, flags, and even resource IDs
- Industry standard for modern CLIs

**Advanced:** Dynamic completions for resource IDs
```bash
canvas assignments get --course-id <TAB>
# Shows list of course IDs from config or recent usage
```

---

#### 2. **Built-in Aliases System**
**Inspired by:** [GitHub CLI aliases](https://cli.github.com/manual/gh_alias_set)
**Status:** ⚠️ Partial (shell aliases work, but not built-in)

**Current workaround:**
```bash
alias ca='canvas assignments'
alias cc='canvas courses'
```

**Better approach:**
```bash
# Set aliases in config
canvas alias set ca 'assignments list --course-id 123'
canvas alias set bugs 'assignments list --bucket ungraded'
canvas alias set whoami 'users me'

# Use them
canvas ca
canvas bugs
canvas whoami
```

**Why it matters:**
- Repetitive commands are common (same course ID, same filters)
- Aliases can include flags and parameters
- Shareable via config file
- gh users LOVE this feature

**Storage:** `~/.canvas-cli/config.yaml`
```yaml
aliases:
  ca: assignments list --course-id 123
  bugs: assignments list --bucket ungraded
  grade: submissions bulk-grade --csv
```

---

#### 3. **Output Filtering & Column Selection**
**Inspired by:** [Modern CLI filtering patterns](https://clig.dev/)
**Status:** ⚠️ Partial (--limit exists, but no filtering)

**Proposed flags:**
```bash
# Filter results (simple substring matching)
canvas courses list --filter "CS 101"
canvas users list --filter "@university.edu"

# Select specific columns
canvas assignments list --columns id,name,due_at,points_possible

# Sort output
canvas assignments list --sort due_at
canvas courses list --sort -enrollment_count  # descending

# Combine them
canvas assignments list --course-id 123 \
  --filter "homework" \
  --columns name,due_at,points_possible \
  --sort due_at \
  --limit 10
```

**Why it matters:**
- Piping to `jq` requires learning jq syntax
- Simple filters cover 80% of use cases
- Table output becomes much more useful
- Can still pipe to jq for complex cases

**Implementation notes:**
- Apply filters AFTER fetching data (client-side)
- Works with all output formats
- Case-insensitive substring matching by default
- `--filter-field name=homework` for field-specific filtering

---

#### 4. **JQ-like Path Expressions** (Light Version)
**Inspired by:** jq, but simpler
**Status:** ❌ Missing

**Proposed flag:**
```bash
# Extract specific fields from JSON output
canvas courses list --select '.[] | {id, name, enrollment_count}'
canvas assignments get 123 --select '.name, .points_possible'

# Works with any output format
canvas users list --select '.[] | select(.login_id | contains("@edu"))'
```

**Why it matters:**
- Many users don't know jq syntax
- Common extractions should be easy
- Alternative: use actual jq (`canvas courses list -o json | jq ...`)

**Decision:** **MAYBE** - Only if we can make it significantly simpler than jq
- **PRO:** Integrated, consistent
- **CON:** Reinventing jq is scope creep
- **ALTERNATIVE:** Excellent documentation for piping to jq

---

#### 5. **Smart Error Messages & Suggestions**
**Inspired by:** [Git's "did you mean" suggestions](https://clig.dev/)
**Status:** ⚠️ Partial (Cobra provides some suggestions)

**Current behavior:**
```bash
$ canvas assignemnt list
Error: unknown command "assignemnt" for "canvas"
Run 'canvas --help' for usage.
```

**Better behavior:**
```bash
$ canvas assignemnt list
Error: unknown command "assignemnt"

Did you mean one of these?
  • assignment
  • assignments

Run 'canvas --help' for usage.
```

**Additional improvements:**
- Suggest required flags when missing: "Missing --course-id. Did you mean to add --course-id 123?"
- Hint at common mistakes: "Note: Use 'assignments' (plural) not 'assignment'"
- Context-aware help: "To create an assignment, use: canvas assignments create --help"

**Implementation:** Levenshtein distance for command suggestions (many Go libraries available)

---

#### 6. **Context/Session Management** (Enhanced)
**Inspired by:** [kubectl contexts](https://github.com/ahmetb/kubectx)
**Status:** ✅ Exists in REPL, but could be global

**Current:** REPL has session variables (`course_id`, `user_id`, `assignment_id`)

**Proposed enhancement:**
```bash
# Set global context (persisted across commands)
canvas context set course 123
canvas context set assignment 456

# Now commands use context automatically
canvas assignments list  # uses course_id from context
canvas submissions list  # uses course_id + assignment_id from context

# Override when needed
canvas assignments list --course-id 789

# Show current context
canvas context show
# Output:
# course_id: 123
# assignment_id: 456
# instance: production

# Clear context
canvas context clear
```

**Why it matters:**
- Typing `--course-id 123` for every command is tedious
- kubectl's context system is beloved by users
- Natural for workflows (work on one course, run multiple commands)
- REPL already has this—extend to global CLI

**Storage:** `~/.canvas-cli/context.yaml` or in-memory for session

**Advanced:** Named contexts
```bash
canvas context save cs101
canvas context use cs101
canvas context list
```

---

### ⚡ MEDIUM IMPACT (Nice to Have)

#### 7. **Interactive Prompts for Required Flags**
**Inspired by:** [gh pr create interactive mode](https://cli.github.com/manual/gh_pr_create)
**Status:** ⚠️ Partial (some commands prompt, not consistent)

**Proposed behavior:**
```bash
$ canvas assignments create
? Course ID: 123
? Assignment name: Homework 1
? Points possible: 100
? Due date (YYYY-MM-DD HH:MM): 2026-02-15 23:59
? Submission types (comma-separated): online_text_entry,online_upload
? Description: |
  Complete the exercises in Chapter 1.
  Submit as PDF.

Creating assignment...
✓ Created assignment "Homework 1" (ID: 789)
```

**With flags (skip prompts):**
```bash
canvas assignments create \
  --course-id 123 \
  --name "Homework 1" \
  --points 100 \
  --due "2026-02-15T23:59:00Z"
```

**Why it matters:**
- Friendlier for new users and occasional operations
- Still supports non-interactive mode for scripts
- gh users love this for `gh pr create`

**Implementation:** Use [survey](https://github.com/AlecAivazis/survey) or [promptui](https://github.com/manifoldco/promptui) library

---

#### 8. **Diff/Preview Mode** (Enhanced Dry-Run)
**Inspired by:** [terraform plan](https://developer.hashicorp.com/terraform/cli/commands/plan)
**Status:** ✅ Exists for bulk-grade, expand to more commands

**Proposed:**
```bash
# Show what would change
canvas assignments update 123 --name "New Name" --dry-run
# Output:
# Would update assignment 123:
#   name: "Old Name" → "New Name"
#
# Run without --dry-run to apply changes.

# For bulk operations
canvas submissions bulk-grade --csv grades.csv --dry-run
# Output:
# Would update 25 submissions:
#   ✓ User 101: 0 → 95
#   ✓ User 102: 0 → 87
#   ✗ User 103: Invalid score "ABC"
#   ...
#
# Summary: 24 valid, 1 error
```

**Apply to:**
- Create, update, delete commands
- Bulk operations
- Configuration changes

---

#### 9. **Watch Mode for Real-Time Monitoring**
**Inspired by:** kubectl, docker
**Status:** ❌ Missing

**Proposed:**
```bash
# Watch submissions for an assignment
canvas submissions list --assignment-id 123 --watch
# Updates every 5 seconds, shows new submissions

# Watch grades being updated
canvas grades history --course-id 123 --watch

# Custom interval
canvas courses list --watch --interval 10s
```

**Why it matters:**
- Useful during active grading periods
- Monitor submission activity in real-time
- See changes as they happen

**Implementation:** Loop with clear screen, configurable interval

---

#### 10. **Progress Indicators for Long Operations**
**Inspired by:** [CLI progress display patterns](https://evilmartians.com/chronicles/cli-ux-best-practices-3-patterns-for-improving-progress-displays)
**Status:** ⚠️ Partial (batch operations have some progress)

**Current:** Silent for most operations

**Better:**
```bash
$ canvas submissions list --course-id 123 --all
Fetching submissions... ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 85% (850/1000)

$ canvas files download --all
Downloading files...
✓ lecture-01.pdf (2.5 MB)
✓ homework-01.docx (156 KB)
⠋ lecture-02.pdf (15.2 MB) [━━━━━━━━━━          ] 60% 9.1 MB
  homework-02.docx (pending)
```

**Types of progress:**
- **Spinner:** For indeterminate operations
- **Progress bar:** For paginated fetches
- **File-by-file:** For downloads/uploads
- **Summary:** Total items processed

**Implementation:** [progressbar](https://github.com/schollz/progressbar) or [mpb](https://github.com/vbauerster/mpb)

---

#### 11. **Fuzzy Finding with FZF Integration**
**Inspired by:** [kubectl + fzf](https://github.com/ahmetb/kubectx)
**Status:** ❌ Missing

**Proposed:**
```bash
# If fzf is installed, enable interactive selection
canvas courses list --fzf
# Opens fzf with course list, select one, outputs course ID

# Pipe to other commands
COURSE_ID=$(canvas courses list --fzf)
canvas assignments list --course-id $COURSE_ID

# Built-in for common workflows
canvas assignments list --course-id $(canvas courses list --fzf)
```

**Why it matters:**
- Remembering course/assignment IDs is hard
- fzf is extremely popular in developer workflows
- Natural for multi-step operations

**Implementation:** Detect if fzf is installed, pipe JSON to fzf with preview

---

#### 12. **Quick Action Shortcuts**
**Inspired by:** gh shortcuts
**Status:** ❌ Missing

**Proposed:**
```bash
# Instead of:
canvas assignments list --course-id 123 --bucket ungraded

# Shortcut:
canvas ungraded 123

# Other shortcuts:
canvas missing 123          # missing submissions
canvas late 123             # late submissions
canvas recent-courses       # courses with recent activity
canvas my-courses           # courses where I'm enrolled
canvas todo                 # assignments due soon
```

**Why it matters:**
- Common workflows should be fast
- Discoverable via `canvas --help`
- Reduces cognitive load

**Implementation:** Aliases that ship by default

---

### 💡 LOW IMPACT (Future Considerations)

#### 13. **Template/Snippet System**
**Status:** ❌ Missing

**Proposed:**
```bash
# Save command templates
canvas template save weekly-hw \
  'assignments create --course-id 123 --points 100 --submission-types online_upload'

# Use with variable substitution
canvas template run weekly-hw --name "Week 5 Homework" --due 2026-03-01

# List templates
canvas template list
```

**Decision:** **LOW PRIORITY** - Aliases cover most use cases

---

#### 14. **Plugin/Extension System**
**Inspired by:** gh extensions
**Status:** ❌ Missing

**Proposed:**
```bash
# Install community extensions
canvas extension install canvas-cli-analytics
canvas extension install canvas-cli-gradebook-export

# Use them
canvas analytics course 123
canvas gradebook-export --format excel
```

**Decision:** **FUTURE** - Only if there's strong community demand

---

#### 15. **Web Browser Integration**
**Inspired by:** gh browse
**Status:** ❌ Missing

**Proposed:**
```bash
# Open course in browser
canvas courses open 123

# Open assignment
canvas assignments open 456 --course-id 123

# Open current context
canvas open  # opens current course from context
```

**Decision:** **LOW PRIORITY** - Nice but not essential

---

#### 16. **Export to Other Formats**
**Status:** ⚠️ CSV exists, could add more

**Proposed:**
```bash
canvas courses list --output excel
canvas grades export --format xlsx
canvas assignments list --format markdown  # for documentation
```

**Decision:** **LOW PRIORITY** - CSV + external tools (csvkit, pandoc) work well

---

## Implementation Priority

### Phase 1: Quick Wins (1-2 weeks)
1. **Shell completion** (2-3 hours) - Massive UX win
2. **Expand dry-run** (4-6 hours) - Safety for destructive operations
3. **Alias system** (1-2 days) - High value, low effort
4. **Better error messages** (1-2 days) - Quality of life

### Phase 2: Core Improvements (2-4 weeks)
5. **Output filtering** (3-5 days) - --filter, --columns, --sort
6. **Context management** (3-5 days) - Persistent course/assignment context
7. **Progress indicators** (2-3 days) - Better feedback for long operations
8. **Interactive prompts** (3-5 days) - User-friendly create commands

### Phase 3: Advanced Features (1-2 months)
9. **Watch mode** (1-2 days) - Real-time monitoring
10. **FZF integration** (2-3 days) - Interactive selection
11. **Diff mode enhancements** (2-3 days) - Better preview
12. **Quick action shortcuts** (1-2 days) - Convenience commands

### Phase 4: Future (Nice to Have)
13. Browser integration
14. Plugin system
15. Additional export formats

---

## Recommendations Summary

### ✅ IMPLEMENT NOW
1. **Shell completion** - Industry standard, Cobra built-in
2. **Alias system** - gh users love this
3. **Output filtering** - `--filter`, `--columns`, `--sort`
4. **Context management** - Like kubectl, reduce repetition
5. **Better error messages** - "Did you mean" suggestions

### ⚠️ CONSIDER
6. **Interactive prompts** - Friendly for new users
7. **Progress indicators** - Better UX for long operations
8. **Watch mode** - Useful for monitoring
9. **FZF integration** - If widely requested

### ❌ SKIP FOR NOW
10. **JQ-like filtering** - Use actual jq, don't reinvent
11. **Plugin system** - Premature, wait for demand
12. **Browser integration** - Low value
13. **Export formats** - CSV + external tools sufficient

---

## Philosophy

**Good CLI UX means:**
- ✅ Fast for daily use (completion, aliases, context)
- ✅ Safe by default (dry-run, confirmations, diff)
- ✅ Helpful errors (suggestions, examples)
- ✅ Scriptable AND interactive (flags + prompts)
- ✅ Composable with Unix tools (JSON output + jq/awk/grep)

**Avoid:**
- ❌ Reinventing existing tools (jq, cron, git)
- ❌ Over-engineering (plugins before there's demand)
- ❌ Breaking Unix philosophy (do one thing well)
- ❌ Surprising behavior (implicit context should be obvious)

---

## Competitive Analysis

| Feature | Canvas CLI | gh | stripe | kubectl | aws-cli |
|---------|-----------|-----|---------|---------|---------|
| Shell completion | ❌ | ✅ | ✅ | ✅ | ✅ |
| Aliases | ❌ | ✅ | ❌ | ✅ | ✅ |
| Output filtering | ⚠️ | ✅ | ⚠️ | ✅ | ✅ |
| Context switching | ⚠️ | ❌ | ❌ | ✅ | ✅ |
| Interactive mode | ✅ | ✅ | ✅ | ❌ | ✅ |
| Dry-run | ⚠️ | ❌ | ❌ | ✅ | ✅ |
| Progress bars | ⚠️ | ✅ | ✅ | ⚠️ | ✅ |
| Watch mode | ❌ | ❌ | ❌ | ✅ | ✅ |
| Multiple profiles | ✅ | ❌ | ❌ | ✅ | ✅ |
| Caching | ✅ | ✅ | ❌ | ❌ | ✅ |

**Key takeaway:** Canvas CLI is competitive but missing industry-standard features (completion, aliases, filtering).

---

## Success Metrics

How to measure UX improvements:

1. **Time to complete common tasks** - Measure before/after
2. **Error rate** - Track failed commands
3. **Command repetition** - Are users typing the same flags repeatedly?
4. **Adoption of new features** - Track alias usage, completion usage
5. **User feedback** - Satisfaction surveys, GitHub issues

---

## References

- [Command Line Interface Guidelines](https://clig.dev/)
- [UX patterns for CLI tools](https://lucasfcosta.com/2022/06/01/ux-patterns-cli-tools.html)
- [12 Factor CLI Apps](https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46)
- [GitHub CLI Manual](https://cli.github.com/manual/)
- [Stripe CLI Documentation](https://docs.stripe.com/stripe-cli)
- [kubectl Best Practices](https://kubernetes.io/docs/reference/kubectl/)
