# Canvas CLI Feature Analysis

## Legend
- ✅ **KEEP** - Valuable feature that belongs in Canvas CLI
- ⚠️ **MAYBE** - Could be useful but needs careful consideration
- ❌ **REJECT** - Belongs in external tools or not valuable enough
- 🟢 **EXISTS** - Already implemented in develop branch

---

## 1. Advanced Reporting & Analytics

### Grade Distribution Reports
**Status:** ❌ **REJECT** - Better handled by external tools
**Reasoning:**
- Canvas UI already provides basic grade distributions
- Visualization (histograms, charts) doesn't fit CLI paradigm
- Users can pipe JSON output to external tools: `canvas grades list | jq | gnuplot`
- External tools (Python/R/Excel) are better suited for statistical analysis

### Engagement Dashboards
**Status:** 🟢 **EXISTS** + ❌ **REJECT**
**What exists:** `canvas analytics` commands already provide activity, participation, communication metrics
**Why reject more:** Dashboards imply visual/interactive UIs - not CLI's strength

### Comparative Analytics
**Status:** ❌ **REJECT** - Complex analysis belongs elsewhere
**Reasoning:**
- Comparing courses over time requires data warehousing
- This is BI/analytics software territory (Tableau, PowerBI, custom Python)
- CLI should focus on data extraction, not complex analytics

### Custom Report Builder
**Status:** ❌ **REJECT** - Scope creep
**Reasoning:**
- Canvas already has report builder in UI
- Complex query building needs a proper query language (SQL-like)
- Better: `canvas api` command + external tools (jq, awk, pandas)

---

## 2. Bulk Content Operations

### Course Cloning Wizard
**Status:** 🟢 **EXISTS** - Already implemented
**What exists:** `canvas content-migrations` supports course copying with selective import

### Batch Assignment Creator
**Status:** ✅ **KEEP** - HIGH VALUE
**Reasoning:**
- Common workflow: creating 10-20 similar assignments (weekly homeworks, etc.)
- CSV/YAML input is natural for CLI: `canvas assignments batch-create --from assignments.csv`
- Fits existing pattern (SIS imports, bulk grading use CSV)
- **UNIQUE VALUE:** Canvas UI requires clicking through forms 20 times

**Suggested Implementation:**
```bash
# CSV with columns: name, points, due_at, description, submission_types
canvas assignments batch-create --course-id 123 --from assignments.csv --dry-run
```

### Content Migration Tools
**Status:** 🟢 **EXISTS** - Already implemented
**What exists:** `canvas content-migrations` handles this

### Bulk File Operations
**Status:** ⚠️ **MAYBE** - Limited value
**Reasoning:**
- Uploading directories: `canvas files upload --course-id 123 --recursive ./materials/`
- Could be useful for instructors migrating materials
- **BUT:** Canvas UI file manager works well for this
- **DECISION:** Only worth if there's significant demand

---

## 3. Automation & Scheduling

### Scheduled Commands
**Status:** ❌ **REJECT** - Use crontab/systemd timers
**Reasoning:**
- Unix philosophy: do one thing well
- `crontab -e` is the standard for scheduling
- Example: `0 9 * * 1 canvas announcements create --course-id 123 --message "Weekly reminder"`
- Adding scheduling duplicates system functionality
- Users who need this already know cron

### Workflow Automation
**Status:** ❌ **REJECT** - Use shell scripts/make
**Reasoning:**
- Chaining commands: just use bash scripts or Makefiles
- Example: `create-assignment.sh` that calls multiple canvas commands
- CLI should be composable, not an orchestration engine
- **COUNTER-ARGUMENT:** Maybe add a simple `--then` flag?
  - Still reject: bash `&&` is clearer

### Event-Driven Actions
**Status:** 🟢 **EXISTS** - Webhook listener already does this
**What exists:** `canvas webhook listen` with event filtering
**Users can:** Write scripts that respond to webhook events

### Template System
**Status:** 🟢 **EXISTS** - Blueprint courses
**What exists:** `canvas blueprint` commands for course templates

---

## 4. Advanced Grading Features

### Grade Import/Export
**Status:** 🟢 **EXISTS** - Bulk grading implemented
**What exists:** `canvas submissions bulk-grade --csv grades.csv --dry-run`

### Grading Workflows
**Status:** ❌ **REJECT** - Business logic doesn't belong in CLI
**Reasoning:**
- Multi-step approval processes are institutional policy
- This requires workflow engine (Airflow, Temporal, etc.)
- Out of scope for a Canvas API client

### Auto-Grading Rules
**Status:** ❌ **REJECT** - Belongs in Canvas or external graders
**Reasoning:**
- Canvas supports auto-grading via LTI tools
- Rule engines (Drools, Python scripts) handle this better
- CLI is for API operations, not business logic

### Grade Analytics
**Status:** ⚠️ **MAYBE** - Simple checks only
**What could work:** `canvas grades audit --course-id 123` to find:
  - Assignments with no grades
  - Students with missing submissions
  - Grade outliers (>2 std dev)
**Keep it simple:** Detection only, not analysis

---

## 5. Content Quality & Accessibility

### Accessibility Scanner
**Status:** ✅ **KEEP** - HIGH VALUE
**Reasoning:**
- WCAG compliance is legally required for many institutions
- Canvas API provides accessibility checking: `/api/v1/courses/:id/content_migrations?type=course_audit`
- CLI can aggregate results across courses
- **UNIQUE VALUE:** Batch scan all courses, export violations to CSV

**Suggested Implementation:**
```bash
canvas courses audit --accessibility --all
canvas pages audit --course-id 123 --wcag-level AA
```

### Broken Link Checker
**Status:** ✅ **KEEP** - MEDIUM VALUE
**Reasoning:**
- Common problem: instructors link to removed files or expired URLs
- Can be implemented by crawling course content via API
- **CAVEAT:** External tools (linkchecker, wget --spider) might be better
- **DECISION:** Implement if it can check Canvas-specific links (files, pages, assignments)

**Suggested Implementation:**
```bash
canvas courses check-links --course-id 123
# Output: pages with broken internal links, missing file references
```

### Content Health Check
**Status:** ✅ **KEEP** - HIGH VALUE
**Reasoning:**
- Common instructor mistake: unpublished modules, missing due dates
- Quick sanity check before semester starts
- Fits CLI pattern: `canvas courses health-check --course-id 123`
- **OUTPUT:** List of issues with actionable fixes

**Suggested Implementation:**
```bash
canvas courses health-check --course-id 123
# Reports:
# - Unpublished modules: 3
# - Assignments without due dates: 5
# - Empty modules: 2
# - Files over quota: 0
```

### Style Consistency
**Status:** ❌ **REJECT** - Too subjective
**Reasoning:**
- "Standards" vary widely by institution
- Linting rules need configuration, exceptions, etc.
- Better: institutions write custom scripts using `canvas api`

---

## 6. Enhanced Scripting Support

### Pipeline Mode
**Status:** ⚠️ **MAYBE** - Interesting but complex
**Reasoning:**
- Unix pipes already work: `canvas courses list -o json | jq '.[] | select(.published == false)'`
- Custom piping between canvas commands adds complexity
- **DECISION:** Reject - existing JSON output + jq is sufficient

### JQ-style Filtering
**Status:** ❌ **REJECT** - Use actual jq
**Reasoning:**
- Don't reinvent jq - it's ubiquitous and powerful
- Better: ensure all commands support `--output json`
- CLI philosophy: compose with other tools

### Batch Files
**Status:** ⚠️ **MAYBE** - Shell scripts already do this
**Reasoning:**
- YAML/JSON command files: `canvas batch run --file commands.yaml`
- **PRO:** Could include error handling, retries, dry-run
- **CON:** Bash scripts with error handling do this
- **DECISION:** Reject - doesn't add enough value over scripts

### Variables & Templating
**Status:** ❌ **REJECT** - Use shell variables
**Reasoning:**
- `export COURSE_ID=123; canvas assignments list --course-id $COURSE_ID` works fine
- envsubst, Jinja2, etc. exist for templating
- Keep CLI simple

---

## 7. Interactive Tools

### Course Builder Wizard
**Status:** ❌ **REJECT** - Canvas UI is better
**Reasoning:**
- Interactive wizards in CLI are awkward (vim/nano are exceptions)
- Canvas UI course setup is actually good
- CLI is for automation/scripting, not guided setup

### Assignment Designer
**Status:** ❌ **REJECT** - Same as above

### Configuration Wizard
**Status:** 🟢 **EXISTS** - Already has interactive config
**What exists:** `canvas config add` prompts for instance URL

### Query Builder
**Status:** ❌ **REJECT** - Use canvas api + curl/HTTPie
**Reasoning:**
- `canvas api GET /api/v1/courses` already exists
- HTTPie, Postman, curl are better for exploring APIs

---

## 8. Third-Party Integrations

### Git Integration
**Status:** ❌ **REJECT** - Git is for code, not course content
**Reasoning:**
- Version control for course pages is interesting but niche
- Canvas already has page revision history
- Complex to implement (converting Canvas content to files)
- **ALTERNATIVE:** Users can export content and version themselves

### Slack/Teams Notifications
**Status:** ⚠️ **MAYBE** - Webhook listener + scripts
**What exists:** `canvas webhook listen` can output events
**Better approach:** Let users pipe webhook events to slack CLI or curl
**Example:**
```bash
# User writes:
canvas webhook listen | while read event; do
  curl -X POST $SLACK_WEBHOOK -d "$event"
done
```
**DECISION:** Reject - composability is better than built-in integrations

### Cloud Storage Sync
**Status:** ❌ **REJECT** - Use rclone/s3cmd
**Reasoning:**
- rclone, s3cmd, gsutil are mature tools for cloud sync
- Canvas files can be downloaded: `canvas files download --all`
- Don't duplicate cloud storage tools

### Calendar Export
**Status:** ✅ **KEEP** - MEDIUM VALUE
**Reasoning:**
- Canvas calendar API exists: `/api/v1/calendar_events`
- Export to iCal format is straightforward
- **USE CASE:** Students want all assignments in Google Calendar
- **IMPLEMENTATION:** `canvas calendar export --output assignments.ics`

### LMS Migration
**Status:** ❌ **REJECT** - Use Canvas's import tools
**Reasoning:**
- Canvas UI has Moodle/Blackboard import
- Migration is complex (format conversions, mapping)
- One-time operation, not worth CLI support

---

## 9. CI/CD & DevOps

### GitHub Actions Integration
**Status:** ✅ **KEEP** - HIGH VALUE
**Reasoning:**
- Common workflow: auto-publish course changes from git repo
- Canvas CLI in CI can sync course content from markdown files
- **EXAMPLE USE CASE:**
  - Store course pages in GitHub as markdown
  - On push to main, CI updates Canvas pages
  - Version control + peer review for course materials
- **IMPLEMENTATION:** Just documentation + examples (canvas already works in CI)

**Suggested Documentation:**
```yaml
# .github/workflows/sync-canvas.yml
name: Sync to Canvas
on: push
jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: |
          echo "$CANVAS_TOKEN" | canvas auth token set
          canvas pages update --course-id 123 --title "Week 1" --body @week1.md
```

### Docker Support
**Status:** ⚠️ **MAYBE** - Just provide Dockerfile
**Reasoning:**
- Canvas CLI doesn't need to run as daemon (unlike webhook listener)
- For CI: users can install from binary or use `go install`
- **DECISION:** Provide example Dockerfile in docs, but low priority

### API Mocking
**Status:** ❌ **REJECT** - Use dedicated mock servers
**Reasoning:**
- Testing tools: wiremock, mockserver, httptest (Go)
- Not CLI's responsibility
- **FOR DEVELOPERS:** Use Go's httptest (already in tests)

### Performance Benchmarking
**Status:** ❌ **REJECT** - Use external profiling
**Reasoning:**
- `time canvas courses list` works fine
- Go has built-in profiling: `go tool pprof`
- Adding benchmarking bloats CLI

---

## 10. Advanced Data Export

### Complete Course Backup
**Status:** 🟢 **EXISTS** - Content migrations do this
**What exists:** `canvas content-migrations create --type course_copy_importer`

### Student Data Export
**Status:** ⚠️ **MAYBE** - GDPR compliance
**Reasoning:**
- GDPR "right to data portability" requires exporting student data
- Canvas API supports this
- **IMPLEMENTATION:** `canvas users export --user-id 123 --gdpr`
- Would export: submissions, grades, discussion posts, files
- **DECISION:** Useful, but niche - maybe later

### Gradebook Snapshots
**Status:** ✅ **KEEP** - MEDIUM VALUE
**Reasoning:**
- Common workflow: backup grades before changes
- Canvas UI doesn't make this easy
- **IMPLEMENTATION:**
```bash
canvas grades export --course-id 123 --output gradebook-2024-01-15.csv
# Later restore if needed:
canvas submissions bulk-grade --csv gradebook-2024-01-15.csv
```
**CAVEAT:** This might already work with submissions bulk-grade

### Audit Trail Export
**Status:** ❌ **REJECT** - Canvas handles this
**Reasoning:**
- Canvas has audit logs in UI (for admins)
- API support is limited
- Compliance/auditing is institution-specific

---

## 11. Data Visualization

### ASCII Charts
**Status:** ❌ **REJECT** - Gimmicky
**Reasoning:**
- Sparklines are cute but limited
- Better: export CSV and use real visualization tools
- Doesn't add real value

### HTML Reports
**Status:** ⚠️ **MAYBE** - For grade reports
**Reasoning:**
- Could generate nice HTML grade reports for students
- **USE CASE:** End-of-semester grade summary
- **DECISION:** Low priority - users can template this themselves

### Trend Analysis
**Status:** ❌ **REJECT** - Use proper analytics tools

---

## 12. Student Support Tools

### Progress Tracking
**Status:** 🟢 **EXISTS** - Analytics commands cover this
**What exists:** `canvas analytics students --course-id 123`

### Risk Identification
**Status:** ⚠️ **MAYBE** - Simple heuristics only
**What could work:** `canvas students at-risk --course-id 123`
**Criteria:** Missing >3 submissions, avg grade <60%, no login in 7 days
**DECISION:** Maybe - if demand exists

### Communication Templates
**Status:** ❌ **REJECT** - Email tools do this better

### Office Hours Scheduler
**Status:** 🟢 **EXISTS** - Calendar reservations
**What exists:** `canvas calendar` commands support appointments

---

## 13. Course Management

### Module Templates
**Status:** 🟢 **EXISTS** - Blueprint courses

### Assignment Libraries
**Status:** ⚠️ **MAYBE** - Overlaps with batch create
**Better approach:** Batch assignment creation from templates

### Rubric Library
**Status:** ⚠️ **MAYBE** - Reusable rubrics
**Reasoning:**
- Canvas supports rubric reuse
- Could be useful: `canvas rubrics list --global` to find reusable rubrics
- **DECISION:** Low priority - UI handles this

### Course Standards
**Status:** ✅ **KEEP** - Ties into health check
**Reasoning:**
- Extension of "content health check"
- Institutions have policies: "all assignments must have due dates", "modules must be published in order"
- **IMPLEMENTATION:** Configurable rules in YAML
```yaml
# .canvas-lint.yml
rules:
  - assignments-must-have-due-dates
  - modules-must-be-published
  - no-broken-links
```
**DECISION:** Keep as part of health-check feature

---

## 14. Performance & Scale

### Multi-Threading
**Status:** 🟢 **EXISTS** - Batch operations use worker pools
**What exists:** `internal/batch/` package for concurrency

### Smart Prefetching
**Status:** ❌ **REJECT** - Over-engineering

### Compression
**Status:** ❌ **REJECT** - HTTP handles this (gzip)

### Offline Mode
**Status:** ❌ **REJECT** - API client can't work offline

---

## 15. Enterprise Features

### Multi-Tenant Management
**Status:** 🟢 **EXISTS** - Multi-instance support
**What exists:** `canvas config` supports multiple Canvas instances

### Role-Based Access
**Status:** ❌ **REJECT** - Canvas handles permissions

### Audit Logging
**Status:** 🟢 **EXISTS** - Telemetry system
**What exists:** `canvas telemetry` tracks command usage

### Usage Analytics
**Status:** 🟢 **EXISTS** - Telemetry system

---

## 16. Quality of Life

### Command Aliases
**Status:** ✅ **KEEP** - LOW EFFORT, NICE UX
**Reasoning:**
- Shell aliases work: `alias ca='canvas assignments'`
- **BETTER:** Built-in aliases in config
```yaml
# ~/.canvas/config.yaml
aliases:
  ca: assignments
  cc: courses
```
**DECISION:** Keep if easy to implement

### Smart Defaults
**Status:** ⚠️ **MAYBE** - Could be annoying
**Reasoning:**
- Remembering `--course-id` could be useful
- **RISK:** Surprising behavior, hard to debug
- **DECISION:** Reject - explicit is better

### Autocomplete
**Status:** ✅ **KEEP** - HIGH VALUE
**Reasoning:**
- Cobra supports shell completion generation
- `canvas completion bash > /etc/bash_completion.d/canvas`
- **EASY TO IMPLEMENT:** Cobra built-in
- **HIGH VALUE:** Significantly improves UX

### Undo/Rollback
**Status:** ❌ **REJECT** - Too complex
**Reasoning:**
- Would need to track all operations + state
- Canvas API doesn't support transactions
- Use `--dry-run` instead

### Dry-Run Mode
**Status:** 🟢 **EXISTS** - For bulk grading
**Expand:** Add to more commands (delete, update)

---

## Summary: Features Worth Implementing

### HIGH PRIORITY ✅
1. **Batch Assignment Creator** - Create multiple assignments from CSV/YAML
2. **Accessibility Scanner** - Check courses for WCAG compliance
3. **Content Health Check** - Find unpublished content, missing dates, empty modules
4. **Shell Autocomplete** - Bash/Zsh completion (Cobra built-in)
5. **Expand Dry-Run** - Add `--dry-run` to delete/update commands

### MEDIUM PRIORITY ⚠️
6. **Broken Link Checker** - Detect broken Canvas internal links
7. **Calendar Export** - Export assignments to iCal format
8. **Gradebook Snapshots** - Backup/restore grades (might exist via bulk-grade)
9. **Grade Audit** - Find missing grades, outliers
10. **CI/CD Documentation** - Examples for GitHub Actions integration

### LOW PRIORITY / REJECT ❌
- Everything else belongs in external tools or is already implemented

---

## Recommended Implementation Order

1. **Autocomplete** (1-2 hours) - Cobra built-in, huge UX win
2. **Expand dry-run** (2-4 hours) - Add to delete, update commands
3. **Batch assignment creator** (1-2 days) - High value, fits existing patterns
4. **Content health check** (2-3 days) - Aggregate existing API calls
5. **Accessibility scanner** (1-2 days) - Use Canvas audit API
6. **Calendar export** (1 day) - iCal format generation
7. **Broken link checker** (2-3 days) - Parse HTML content, check links

---

## Philosophy Summary

**Canvas CLI should:**
- ✅ Wrap Canvas API operations cleanly
- ✅ Enable automation and scripting
- ✅ Provide data extraction and bulk operations
- ✅ Integrate well with Unix tools (jq, grep, awk)

**Canvas CLI should NOT:**
- ❌ Duplicate functionality of external tools (cron, git, jq, cloud storage)
- ❌ Implement complex business logic or analytics
- ❌ Try to be a workflow engine or task scheduler
- ❌ Replace Canvas UI for interactive operations

**Guiding principle:** If `canvas COMMAND | external-tool` works, don't build it into canvas-cli.
