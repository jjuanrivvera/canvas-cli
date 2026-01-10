# Canvas CLI v1.1.0 - Feature Plan

**Created:** 2026-01-10
**Target:** Version 1.1.0
**Current Version:** 1.0.0

---

## Summary

Version 1.1.0 focuses on expanding Canvas LMS content management capabilities. Based on Canvas API documentation analysis, the following high-value features are planned.

**Current Commands (v1.0.0):**
- `auth` - OAuth 2.0 with PKCE authentication
- `courses` - Course listing and details
- `assignments` - Assignment CRUD operations
- `submissions` - Submission viewing and grading
- `enrollments` - Enrollment management
- `files` - File upload/download management
- `users` - User management
- `accounts` - Account information
- `webhook` - Real-time webhook listener
- `sync` - Batch operations via CSV
- `doctor` - Diagnostics and health checks
- `shell` - Interactive REPL mode

---

## New Features for v1.1.0

### 1. Modules Command (Priority: High)

Modules are the primary way to organize course content in Canvas. This is a highly requested feature.

**API Reference:** `.ai/canvas-lms-docs/modules.md`

**Commands:**
```bash
# List all modules in a course
canvas modules list --course-id 123

# Get module details
canvas modules get --course-id 123 --module-id 456

# Create a new module
canvas modules create --course-id 123 --name "Week 1" --position 1

# Update a module
canvas modules update --course-id 123 --module-id 456 --name "Week 1: Introduction"

# Delete a module
canvas modules delete --course-id 123 --module-id 456

# List module items
canvas modules items list --course-id 123 --module-id 456

# Add item to module
canvas modules items add --course-id 123 --module-id 456 --type Assignment --content-id 789

# Reorder module items
canvas modules items reorder --course-id 123 --module-id 456 --item-ids 1,2,3
```

**API Endpoints:**
- `GET /api/v1/courses/:course_id/modules` - List modules
- `GET /api/v1/courses/:course_id/modules/:id` - Show module
- `POST /api/v1/courses/:course_id/modules` - Create module
- `PUT /api/v1/courses/:course_id/modules/:id` - Update module
- `DELETE /api/v1/courses/:course_id/modules/:id` - Delete module
- `GET /api/v1/courses/:course_id/modules/:module_id/items` - List items
- `POST /api/v1/courses/:course_id/modules/:module_id/items` - Create item

---

### 2. Pages Command (Priority: High)

Wiki pages provide rich content for courses. Essential for content management.

**API Reference:** `.ai/canvas-lms-docs/pages.md`

**Commands:**
```bash
# List all pages in a course
canvas pages list --course-id 123

# Get page content (by URL or ID)
canvas pages get --course-id 123 --page my-page-title

# Create a new page
canvas pages create --course-id 123 --title "Welcome" --body "<p>Hello!</p>"

# Update a page
canvas pages update --course-id 123 --page my-page --title "New Title"

# Delete a page
canvas pages delete --course-id 123 --page my-page

# Get front page
canvas pages front --course-id 123

# Set as front page
canvas pages set-front --course-id 123 --page my-page

# List page revisions
canvas pages revisions --course-id 123 --page my-page

# Revert to revision
canvas pages revert --course-id 123 --page my-page --revision-id 5
```

**API Endpoints:**
- `GET /api/v1/courses/:course_id/pages` - List pages
- `GET /api/v1/courses/:course_id/pages/:url_or_id` - Show page
- `POST /api/v1/courses/:course_id/pages` - Create page
- `PUT /api/v1/courses/:course_id/pages/:url_or_id` - Update page
- `DELETE /api/v1/courses/:course_id/pages/:url_or_id` - Delete page
- `GET /api/v1/courses/:course_id/front_page` - Get front page
- `PUT /api/v1/courses/:course_id/front_page` - Set front page

---

### 3. Announcements Command (Priority: High)

Course announcements for communication with students.

**API Reference:** `.ai/canvas-lms-docs/announcements.md`

**Commands:**
```bash
# List announcements across courses
canvas announcements list --courses course_123,course_456

# List announcements for a single course
canvas announcements list --course-id 123

# Filter by date range
canvas announcements list --course-id 123 --start-date 2026-01-01 --end-date 2026-01-31

# Get latest announcement per course
canvas announcements list --courses course_123,course_456 --latest-only

# Create announcement (uses discussion topics API)
canvas announcements create --course-id 123 --title "Welcome!" --message "Hello students!"

# Include section information
canvas announcements list --course-id 123 --include sections
```

**API Endpoints:**
- `GET /api/v1/announcements` - List announcements (with context_codes)
- Note: Creating announcements uses the Discussion Topics API with `is_announcement=true`

---

### 4. Calendar Events Command (Priority: Medium)

Calendar management for events and scheduling.

**API Reference:** `.ai/canvas-lms-docs/calendar-events.md`

**Commands:**
```bash
# List calendar events
canvas calendar list --start-date 2026-01-01 --end-date 2026-01-31

# Filter by context
canvas calendar list --context course_123

# Get event details
canvas calendar get --event-id 456

# Create a calendar event
canvas calendar create --context course_123 \
  --title "Office Hours" \
  --start "2026-01-15T14:00:00Z" \
  --end "2026-01-15T15:00:00Z" \
  --location "Room 101"

# Update an event
canvas calendar update --event-id 456 --title "Updated Title"

# Delete an event
canvas calendar delete --event-id 456

# Create recurring event
canvas calendar create --context course_123 \
  --title "Weekly Meeting" \
  --start "2026-01-15T10:00:00Z" \
  --end "2026-01-15T11:00:00Z" \
  --rrule "FREQ=WEEKLY;COUNT=10"
```

**API Endpoints:**
- `GET /api/v1/calendar_events` - List events
- `GET /api/v1/calendar_events/:id` - Show event
- `POST /api/v1/calendar_events` - Create event
- `PUT /api/v1/calendar_events/:id` - Update event
- `DELETE /api/v1/calendar_events/:id` - Delete event

---

### 5. Discussion Topics Command (Priority: Medium)

Forums and discussion management.

**API Reference:** `.ai/canvas-lms-docs/discussion-topics.md`

**Commands:**
```bash
# List discussion topics
canvas discussions list --course-id 123

# Filter by type
canvas discussions list --course-id 123 --only-announcements
canvas discussions list --course-id 123 --exclude-announcements

# Get topic details
canvas discussions get --course-id 123 --topic-id 456

# Create a discussion
canvas discussions create --course-id 123 \
  --title "Week 1 Discussion" \
  --message "What did you learn?" \
  --discussion-type threaded

# Update a discussion
canvas discussions update --course-id 123 --topic-id 456 --pinned

# Delete a discussion
canvas discussions delete --course-id 123 --topic-id 456

# List discussion entries
canvas discussions entries --course-id 123 --topic-id 456

# Post a reply
canvas discussions reply --course-id 123 --topic-id 456 \
  --message "Great point!"

# Mark as read
canvas discussions mark-read --course-id 123 --topic-id 456
```

**API Endpoints:**
- `GET /api/v1/courses/:course_id/discussion_topics` - List topics
- `GET /api/v1/courses/:course_id/discussion_topics/:topic_id` - Show topic
- `POST /api/v1/courses/:course_id/discussion_topics` - Create topic
- `PUT /api/v1/courses/:course_id/discussion_topics/:topic_id` - Update topic
- `DELETE /api/v1/courses/:course_id/discussion_topics/:topic_id` - Delete topic
- `GET /api/v1/courses/:course_id/discussion_topics/:topic_id/entries` - List entries
- `POST /api/v1/courses/:course_id/discussion_topics/:topic_id/entries` - Create entry

---

### 6. Planner Command (Priority: Low)

Student planner for task management.

**API Reference:** `.ai/canvas-lms-docs/planner.md`

**Commands:**
```bash
# List planner items
canvas planner list

# Filter by date
canvas planner list --start-date 2026-01-01 --end-date 2026-01-31

# Filter by context
canvas planner list --context course_123

# Filter by type
canvas planner list --filter assignments
canvas planner list --filter quizzes

# List planner notes
canvas planner notes list

# Create a planner note
canvas planner notes create --title "Study for exam" --todo-date 2026-01-20

# Update a planner note
canvas planner notes update --note-id 123 --title "Updated title"

# Delete a planner note
canvas planner notes delete --note-id 123

# Mark item complete
canvas planner complete --item-id 456

# Dismiss item
canvas planner dismiss --item-id 456
```

**API Endpoints:**
- `GET /api/v1/planner/items` - List planner items
- `GET /api/v1/planner_notes` - List notes
- `POST /api/v1/planner_notes` - Create note
- `PUT /api/v1/planner_notes/:id` - Update note
- `DELETE /api/v1/planner_notes/:id` - Delete note
- `PUT /api/v1/planner/overrides/:id` - Update override

---

## Implementation Order

### Phase 1 (Core Content Management)
1. **Modules** - Most requested, fundamental to course organization
2. **Pages** - Essential for rich content

### Phase 2 (Communication)
3. **Announcements** - Course-wide communication
4. **Discussion Topics** - Student engagement

### Phase 3 (Scheduling & Planning)
5. **Calendar Events** - Time-based content
6. **Planner** - Student task management

---

## Technical Considerations

### API Client Extensions

New files needed in `internal/api/`:
- `modules.go` - Module and ModuleItem types and methods
- `pages.go` - Page and PageRevision types and methods
- `announcements.go` - Announcement types (uses Discussion types)
- `calendar.go` - CalendarEvent types and methods
- `discussions.go` - DiscussionTopic and Entry types and methods
- `planner.go` - PlannerItem and PlannerNote types and methods

### Command Files

New files needed in `commands/`:
- `modules.go` - Module commands
- `pages.go` - Page commands
- `announcements.go` - Announcement commands
- `calendar.go` - Calendar commands
- `discussions.go` - Discussion commands
- `planner.go` - Planner commands

### Tests

Corresponding test files:
- `internal/api/*_test.go` - API client tests
- `commands/*_test.go` - Command integration tests

---

## Breaking Changes

None planned. All changes are additive.

---

## Migration Notes

No migration required from v1.0.0 to v1.1.0.

---

## Release Checklist

- [ ] Implement Modules command
- [ ] Implement Pages command
- [ ] Implement Announcements command
- [ ] Implement Calendar command
- [ ] Implement Discussions command
- [ ] Implement Planner command
- [ ] Update README with new commands
- [ ] Update shell completion for new commands
- [ ] Add new commands to REPL
- [ ] Write comprehensive tests
- [ ] Update documentation
- [ ] Create v1.1.0 tag
